package service

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"winsleepd"
	daemon "winsleepd/cmd/winsleepd"
)

type Service struct {
	ServiceName string
	Description string
	installed   bool
	running     bool
	debug       bool
}

var service *Service

func NewService() *Service {
	service = &Service{
		ServiceName: "winsleepd",
		Description: "Stupidly Simple Sleep Daemon",
		installed:   false,
		running:     false,
		debug:       false,
	}
	return service
}

func (s *Service) Install() {
	daemon.InstallService(s.ServiceName, s.Description)
}

func (s *Service) Uninstall() {
	if s.IsInstalled() {
		daemon.RemoveService(s.ServiceName)
	}
	s.installed = false
}

func (s *Service) Config() {
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
}

func (s *Service) Start() {
	daemon.StartService(s.ServiceName)
}

func (s *Service) Stop() {
	if s.running {
		daemon.ControlService(s.ServiceName, svc.Stop, svc.Stopped)
	}
	s.running = false
}

func (s *Service) Pause() {
	if s.running {
		daemon.ControlService(s.ServiceName, svc.Pause, svc.Paused)
	}
	s.running = false
}

func (s *Service) Continue() {
	if s.running {
		daemon.ControlService(s.ServiceName, svc.Continue, svc.Running)
	}
	s.running = false
}

func (s *Service) Sleep() {
	winsleepd.Sleep()
}

func Get() *Service {
	if service != nil {
		return service
	}
	service = NewService()
	return service
}

func (s *Service) IsInstalled() bool {
	servicesmsc, err := mgr.Connect()
	if err != nil {
		return false
	}
	defer servicesmsc.Disconnect()
	o, err := servicesmsc.OpenService(s.ServiceName)
	if err == nil {
		o.Close()
		s.installed = true
		return true
	}
	s.installed = false
	return false
}
