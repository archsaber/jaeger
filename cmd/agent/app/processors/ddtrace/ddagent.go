package ddtrace

import (
	"context"

	"github.com/DataDog/datadog-trace-agent/cmd/ddtrace"
	ddconfig "github.com/DataDog/datadog-trace-agent/config"
)

type ddAgent struct {
	agent      *ddtrace.Agent
	stopFunc   context.CancelFunc
	isFinished chan struct{}
}

func newDDAgent(c *ddconfig.AgentConfig) *ddAgent {
	ctx, cancel := context.WithCancel(context.Background())

	return &ddAgent{
		agent:      ddtrace.NewAgent(ctx, c),
		stopFunc:   cancel,
		isFinished: make(chan struct{}),
	}
}

func (d *ddAgent) start() {
	d.agent.Run()
	close(d.isFinished)
}

func (d *ddAgent) stop() {
	d.stopFunc()
	<-d.isFinished
}
