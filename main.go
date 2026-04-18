/* 

Updates the Windows Update pause period continuously in the background in order
to prevent Windows update from ever triggering.

This is done by modifying the following key:

	HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\WindowsUpdate\UX\Settings

In there the following entries control the pause:

	
	PauseUpdatesStartTime
	PauseUpdatesExpiryTime
	PauseFeatureUpdatesStartTime
	PauseFeatureUpdatesEndTime
	PauseQualityUpdatesStartTime
	PauseQualityUpdatesEndTime

The pause start is set to current time and the end time is set to 7 days in the
future.

The pause is cleared when these entries are deleted.

This program can be run as a Windows service, or from the command line with
"set" or "clear" as its argument.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	serviceName 	= "WUPauseService"
	keyPath 	= `SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`
	isoFormat	= "2006-01-02T15:04:05Z"
)

var pauseValues = []string{
	"PauseUpdatesStartTime",
	"PauseUpdatesExpiryTime",
	"PauseFeatureUpdatesStartTime",
	"PauseFeatureUpdatesEndTime",
	"PauseQualityUpdatesStartTime",
	"PauseQualityUpdatesEndTime",
}

func setPause() error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("Error opening registry: %w", err)
	}
	defer k.Close()

	now := time.Now().UTC()
	start := now.Format(isoFormat)
	end := now.Add(7 * 24 * time.Hour).Format(isoFormat)

	pairs := map[string]string{
		"PauseUpdatesStartTime":        start,
		"PauseUpdatesExpiryTime":       end,
		"PauseFeatureUpdatesStartTime": start,
		"PauseFeatureUpdatesEndTime":   end,
		"PauseQualityUpdatesStartTime": start,
		"PauseQualityUpdatesEndTime":   end,
	}

	for name, val := range pairs {
		if err := k.SetStringValue(name, val); err != nil {
			return fmt.Errorf("Error setting key: %s", name)
		}
	}

	return nil
}

func clearPause() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.SET_VALUE)

	if err != nil {
		return
	}
	defer k.Close()

	for _, v := range pauseValues {
		k.DeleteValue(v)
	}
}


type wupeService struct {
	elog debug.Log
}

func (s *wupeService) applyPause() {
	if err := setPause(); err != nil {
		s.elog.Error(1, fmt.Sprintf("Could not set pause: %v", err))
	} else {
		s.elog.Info(1, "Pause window updated")
	}
}

func (s *wupeService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const accepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}
	s.applyPause()
	changes <- svc.Status{State: svc.Running, Accepts: accepted}
	
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	loop: for {
		select {
		case <-ticker.C:
			s.applyPause()
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			default:
				s.elog.Warning(1, fmt.Sprintf("Unexpected conrol signal: %d", c))
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	clearPause()
	s.elog.Info(1, "Update pause cleared, exiting")
	changes <- svc.Status{State: svc.Stopped}
	return
}

func runService(name string) {
	elog, err := eventlog.Open(name)
	if err != nil {
		return
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("Starting %s", name))
	if err := svc.Run(name, &wupeService{elog: elog}); err != nil {
		elog.Error(1, fmt.Sprintf("Could not run service %s: %v", name, err))
	}
	elog.Info(1, fmt.Sprintf("Stopping %s", name))
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "set":
			if err := setPause(); err != nil {
				fmt.Println("Error setting the pause:", err)
				os.Exit(1)
			}
			fmt.Println("Pause set")
			return
		case "clear":
			clearPause()
			fmt.Println("Pause cleared")
			return
		}
	}

	runService(serviceName)
}
