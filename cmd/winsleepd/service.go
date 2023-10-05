// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"winsleepd"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

var (
	ActionSleep      = "sleep"
	ActionHibernate  = "hibernate"
	ActionScreenOff  = "screenoff"
	ActionLockScreen = "lockscreen"
)

type Configuration struct {
	Timeout string `json:"timeout"`
	Action  string `json:"action"`
}

func GetConfiguration() (*Configuration, error) {
	//dir, err := os.UserHomeDir()
	//if err != nil {
	//	return nil, err
	//}
	dir := "C:\\"
	file, err := os.Open(dir + ".winsleepd.json")
	if err != nil {
		if os.IsNotExist(err) {
			return NewConfiguration()
		}
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func NewConfiguration() (*Configuration, error) {
	//dir, err := os.UserHomeDir()
	//if err != nil {
	//	return nil, err
	//}
	dir := "C:\\"
	file, err := os.Create(filepath.Join(dir, ".winsleepd.json"))
	if err != nil {
		return nil, err
	}
	encoder := json.NewEncoder(file)
	configuration := Configuration{
		Timeout: "30m", // 30 minutes
		Action:  ActionScreenOff,
	}
	encoder.SetIndent("", "\t")
	err = encoder.Encode(configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

type daemon struct {
	mouseX     int
	mouseY     int
	lastChange time.Time
	lastLogMsg string
}

func (m *daemon) Tick() {
	x, y := winsleepd.GetMousePos()
	if x != m.mouseX || y != m.mouseY {
		m.mouseX = x
		m.mouseY = y
		m.lastChange = time.Now()
	}
	config, err := GetConfiguration()
	if err != nil {
		elog.Error(1, fmt.Sprintf("error getting configuration: %v", err))
	}
	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		elog.Error(1, fmt.Sprintf("error parsing timeout: %v", err))
	}
	if time.Now().Sub(m.lastChange) > timeout {
		switch config.Action {
		case ActionSleep:
			logmsg := "Sleeping..."
			if logmsg != m.lastLogMsg {
				m.lastLogMsg = logmsg
				elog.Info(1, logmsg)
			}
			winsleepd.Sleep()
		case ActionHibernate:
			logmsg := "Hibernating..."
			if logmsg != m.lastLogMsg {
				m.lastLogMsg = logmsg
				elog.Info(1, logmsg)
			}
			winsleepd.Hibernate()
		case ActionScreenOff:
			logmsg := "Turning screen off..."
			if logmsg != m.lastLogMsg {
				m.lastLogMsg = logmsg
				elog.Info(1, logmsg)
			}
			winsleepd.ScreenOff()
		case ActionLockScreen:
			logmsg := "Locking Screen..."
			if logmsg != m.lastLogMsg {
				m.lastLogMsg = logmsg
				elog.Info(1, logmsg)
			}
			winsleepd.LockScreen()
		}
	}
}

func (m *daemon) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	cfg, err := GetConfiguration()
	if err != nil {
		elog.Error(1, fmt.Sprintf("error getting configuration: %v", err))
	}
	elog.Info(1, fmt.Sprintf("timeout: %v", cfg.Timeout))
	elog.Info(1, fmt.Sprintf("action: %v", cfg.Action))
	changes <- svc.Status{State: svc.StartPending}
	fasttick := time.Tick(10 * time.Second)
	slowtick := time.Tick(time.Minute)
	tick := fasttick
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case <-tick:
			//elog.Info(1, "beep")
			m.Tick()
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				testOutput := strings.Join(args, "-")
				testOutput += fmt.Sprintf("-%d", c.Context)
				elog.Info(1, testOutput)
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				tick = slowtick
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tick = fasttick
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &daemon{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
