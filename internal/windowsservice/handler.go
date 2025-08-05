package windowsservice

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

type StarStopper interface {
	Start() error
	Stop()
}

type Handler struct {
	service StarStopper
}

func NewHandler(s StarStopper) *Handler {
	return &Handler{service: s}
}

func (s *Handler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPowerEvent | svc.AcceptHardwareProfileChange
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	err := s.service.Start()

	if err != nil {
		changes <- svc.Status{State: svc.StopPending}
		return
	}

main:
	for c := range r {
		log.Infof("service status request: %v", c.Cmd)
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.PowerEvent:
			changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
		case svc.Stop, svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}
			s.service.Stop()
			changes <- svc.Status{State: svc.Stopped}
			break main
		}
	}

	return
}
