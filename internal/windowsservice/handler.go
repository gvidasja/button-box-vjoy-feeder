package windowsservice

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

type handler struct {
	service *cli
}

func (s *handler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	err := s.service.start()

	if err != nil {
		changes <- svc.Status{State: svc.StopPending}
		return
	}

	for c := range r {
		log.Infof("service status request: %v", c.Cmd)
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			s.service.stop()
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}
