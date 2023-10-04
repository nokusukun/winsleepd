package service

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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
	service.installed = service.IsInstalled()
	return service
}

func Get() *Service {
	if service != nil {
		return service
	}
	return NewService()
}

func runElevated() {
	cmd := exec.Command("powershell", "Start-Process", "powershell", "-Verb", "runas", "-ArgumentList", os.Args[0]+" elevated")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error elevating", err)
	}
}

func (s *Service) Install() {
	if s.installed {
		return
	}
	err := daemon.InstallService(s.ServiceName, s.Description)
	if err != nil {
		log.Fatalf("failed to install service: %v", err)
		return
	}
	s.installed = true
}

func (s *Service) Uninstall() {
	if s.IsInstalled() {
		err := daemon.RemoveService(s.ServiceName)
		if err != nil {
			log.Fatalf("failed to remove service: %v", err)
			return
		}
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
	err := daemon.StartService(s.ServiceName)
	if err != nil {
		log.Fatalf("failed to start service: %v", err)
	}
}

func (s *Service) Stop() {
	if s.running {
		err := daemon.ControlService(s.ServiceName, svc.Stop, svc.Stopped)
		if err != nil {
			log.Fatalf("failed to stop service: %v", err)
		}
	}
	s.running = false
}

func (s *Service) Pause() {
	if s.running {
		err := daemon.ControlService(s.ServiceName, svc.Pause, svc.Paused)
		if err != nil {
			log.Fatalf("failed to pause service: %v", err)
		}
	}
	s.running = false
}

func (s *Service) Continue() {
	if s.running {
		err := daemon.ControlService(s.ServiceName, svc.Continue, svc.Running)
		if err != nil {
			log.Fatalf("failed to continue service: %v", err)
		}
	}
	s.running = false
}

func (s *Service) Sleep() {
	winsleepd.Sleep()
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
