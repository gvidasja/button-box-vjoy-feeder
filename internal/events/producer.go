package events

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Producer interface {
	Produce(event string, data any)
}

type AppEventProducer struct {
	last     map[string]time.Time
	interval time.Duration
	app      *application.App
}

func NewAppEventProducer(interval time.Duration) *AppEventProducer {
	return &AppEventProducer{
		last:     map[string]time.Time{},
		interval: interval,
	}
}

func (p *AppEventProducer) Produce(event string, data any) {
	logrus.Debug("Producing event: ", event, " with data: ", data, p.interval, p.app != nil)

	if p.app == nil {
		return
	}

	if p.app.Window.Current() == nil || !p.app.Window.Current().IsVisible() {
		return
	}

	if time.Since(p.last[event]) < p.interval {
		return
	}

	p.app.Event.Emit(event, data)
	p.last[event] = time.Now()
}

func (p *AppEventProducer) SetApp(app *application.App) {
	p.app = app
}
