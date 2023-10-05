// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

// Example service program that beeps.
//
// The program demonstrates how to create Windows service and
// install / remove it on a computer. It also shows how to
// stop / start / pause / continue any service, and how to
// write to event log. It also shows how to use debug
// facilities available in debug package.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"winsleepd"
	tui "winsleepd/cmd/tui/cmd"
	"winsleepd/cmd/winsleepd"

	"golang.org/x/sys/windows/svc"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

var ServiceName = "winsleepd"
var Description = "Stupidly Simple Sleep Daemon"
var InstallAsUser = false

func main() {
	flag.StringVar(&ServiceName, "name", ServiceName, "name of the service")
	flag.BoolVar(&InstallAsUser, "as-user", InstallAsUser, "install as user")
	flag.Parse()

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}
	if inService {
		daemon.RunService(ServiceName, false)
		return
	}

	if len(os.Args) < 2 {
		//usage("no command specified")
		err := tui.Run()
		if err != nil {
			log.Fatalf("failed to run program: %v", err)
		}
		return
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		daemon.RunService(ServiceName, true)
		return
	case "install":
		_, err = daemon.NewConfiguration()
		if err != nil {
			log.Fatalf("failed to create configuration: %v", err)
		}

		err = daemon.InstallService(ServiceName, Description, InstallAsUser)
		if err != nil {
			log.Fatalf("failed to install service: %v", err)
		}
	case "config":
		dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get user home directory: %v", err)
			return
		}
		err = exec.Command("cmd", "/C", "start", "", filepath.Join(dir, ".winsleepd.json")).Run()
		if err != nil {
			log.Fatalf("failed to open configuration: %v", err)
			return
		}
	case "remove":
		err = daemon.RemoveService(ServiceName)
	case "start":
		err = daemon.StartService(ServiceName)
	case "stop":
		err = daemon.ControlService(ServiceName, svc.Stop, svc.Stopped)
	case "pause":
		err = daemon.ControlService(ServiceName, svc.Pause, svc.Paused)
	case "continue":
		err = daemon.ControlService(ServiceName, svc.Continue, svc.Running)
	case "tui":
		err = tui.Run()
	case "debug:sleep":
		log.Println("sleeping")
		winsleepd.Sleep()
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, ServiceName, err)
	}
	return
}
