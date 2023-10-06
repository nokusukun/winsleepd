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
	ServiceName   string
	Description   string
	Configuration *daemon.Configuration
	debug         bool
}

var service *Service
var config *daemon.Configuration

func NewService() *Service {
	service = &Service{
		ServiceName: "winsleepd",
		Description: "Stupidly Simple Sleep Daemon",
		debug:       false,
	}
	return service
}

func Get() *Service {
	if service != nil {
		return service
	}
	return NewService()
}

func GetConfig() *daemon.Configuration {
	if config != nil {
		return config
	}
	c, err := daemon.GetConfiguration() // also creates a config file if it doesn't exist
	if err != nil {
		log.Fatalf("failed to get configuration: %v", err)
		return nil
	}
	return c
}

func NewConfig() *daemon.Configuration {
	if config != nil {
		return config
	}
	config, err := daemon.NewConfiguration()
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
		return nil
	}
	return config
}

func runElevated() {
	cmd := exec.Command("powershell", "Start-Process", "powershell", "-Verb", "runas", "-ArgumentList", os.Args[0]+" elevated")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error elevating", err)
	}
}

func (s *Service) Install(asUser bool) {
	s.Configuration = GetConfig() // also creates a config file if it doesn't exist
	if s.IsInstalled() {
		return
	}
	err := daemon.InstallService(s.ServiceName, s.Description, asUser)
	if err != nil {
		log.Fatalf("failed to install service: %v", err)
		return
	}
	s.IsInstalled()
}

func (s *Service) Uninstall() {
	if s.IsInstalled() {
		if s.IsPaused() {
			s.Continue()
		}
		if s.IsRunning() {
			s.Stop()
		}
		err := daemon.RemoveService(s.ServiceName)
		if err != nil {
			log.Fatalf("failed to remove service: %v", err)
			return
		}
	}
}

func (s *Service) OpenConfig() {
	//dir, err := os.UserHomeDir()
	//if err != nil {
	//	log.Fatalf("failed to get user home directory: %v", err)
	//	return
	//}
	//dir, err := os.UserHomeDir()
	//if err != nil {
	//	return nil, err
	//}
	dir := "C:\\"
	err := exec.Command("cmd", "/C", "start", "", filepath.Join(dir, ".winsleepd.json")).Run()
	if err != nil {
		log.Fatalf("failed to open configuration: %v", err)
		return
	}
}

func (s *Service) Start() {
	if s.IsRunning() || s.IsPaused() {
		return
	}
	err := daemon.StartService(s.ServiceName)
	if err != nil {
		log.Fatalf("failed to start service: %v", err)
	}
	s.IsRunning()
}

func (s *Service) Stop() {
	if s.IsRunning() || s.IsPaused() {
		err := daemon.ControlService(s.ServiceName, svc.Stop, svc.Stopped)
		if err != nil {
			log.Fatalf("failed to stop service: %v", err)
		}
	}
	s.IsRunning()
}

func (s *Service) Pause() {
	if s.IsRunning() && !s.IsPaused() {
		err := daemon.ControlService(s.ServiceName, svc.Pause, svc.Paused)
		if err != nil {
			log.Fatalf("failed to pause service: %v", err)
		}
	}
}

func (s *Service) Continue() {
	if s.IsPaused() {
		err := daemon.ControlService(s.ServiceName, svc.Continue, svc.Running)
		if err != nil {
			log.Fatalf("failed to continue service: %v", err)
		}
	}
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
		return true
	}
	return false
}

func (s *Service) QueryState() svc.State {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatalf("could not connect to service manager: %v", err)
		return svc.Stopped
	}
	defer m.Disconnect()

	service, err := m.OpenService(s.ServiceName)
	if err != nil {
		log.Fatalf("could not access service: %v", err)
		return svc.Stopped
	}
	defer service.Close()

	status, err := service.Query()
	if err != nil {
		log.Fatalf("could not query service status: %v", err)
		return svc.Stopped
	}
	return status.State
}

func (s *Service) IsRunning() bool {
	return s.QueryState() == svc.Running
}

func (s *Service) IsPaused() bool {
	return s.QueryState() == svc.Paused
}
