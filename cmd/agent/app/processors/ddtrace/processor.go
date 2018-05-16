package ddtrace

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/DataDog/datadog-trace-agent/cmd/ddtrace"
	ddconfig "github.com/DataDog/datadog-trace-agent/config"
	"github.com/fsnotify/fsnotify"
	"github.com/jaegertracing/jaeger/cmd/agent/app/auth"
	"go.uber.org/zap"
)

const logFilePath = "/var/log/archsaber/archsaber-trace-agent.log"

// ProcessorConfig is the config for DDTraceProcessor
type ProcessorConfig struct {
	Enabled         bool    `yaml:"enabled"`
	HostPort        string  `yaml:"hostPort" validate:"nonzero"`
	ConnectionLimit int     `yaml:"connectionLimit"`
	ReceiverTimeout int     `yaml:"receiverTimeout"`
	ExtraSampleRate float64 `yaml:"extraSampleRate"`
	MaxTPS          float64 `yaml:"maxtps"`
}

// Processor is a collector that uses HTTP protocol and just holds
// a chan where the spans received are sent one by one
type Processor struct {
	ddAgent     *ddtrace.Agent
	ddAgentStop context.CancelFunc

	conf ProcessorConfig

	isRunning  bool
	stopCalled bool
	*sync.Mutex
}

// NewProcessor returns a pointer to a new DDServer
func (c ProcessorConfig) NewProcessor() (*Processor, error) {
	ddAgentConfig, err := newDDAgentConfig(c)
	if err != nil {
		return nil, err
	}
	ddAgentContext, ddAgentCancel := context.WithCancel(context.Background())
	ddtraceProcessor := Processor{
		ddAgent:     ddtrace.NewAgent(ddAgentContext, ddAgentConfig),
		ddAgentStop: ddAgentCancel,
		conf:        c,
		isRunning:   false,
		stopCalled:  false,
		Mutex:       &sync.Mutex{},
	}
	auth.AddTokenUpdateAction(ddtraceProcessor.restartOnTokenUpdate)

	return &ddtraceProcessor, nil
}

// Serve starts doing the HTTP server and is ready to receive traces
func (r *Processor) Serve() {
	r.Lock()
	defer r.Unlock()
	// Start the agent if it is not already running
	if !r.isRunning {
		go r.ddAgent.Run()
		r.isRunning = true
		r.stopCalled = false
	}
}

// Stop stops the processor
func (r *Processor) Stop() {
	r.Lock()
	r.isRunning = false
	r.stopCalled = true
	r.ddAgentStop()
	r.Unlock()
}

func (r *Processor) restartOnTokenUpdate(event fsnotify.Event, logger *zap.Logger) {
	r.Lock()
	defer r.Unlock()
	if r.stopCalled == true {
		return
	}
	if event.Op&fsnotify.Write != fsnotify.Write &&
		event.Op&fsnotify.Create != fsnotify.Create {
		return
	}
	logger.Info("Restarting ddagent due to token update")
	r.ddAgentStop()
	ddAgentConfig, err := newDDAgentConfig(r.conf)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	ddAgentContext, ddAgentCancel := context.WithCancel(context.Background())
	r.ddAgent = ddtrace.NewAgent(ddAgentContext, ddAgentConfig)
	r.ddAgentStop = ddAgentCancel
	go r.ddAgent.Run()
	r.isRunning = true
}

func newDDAgentConfig(c ProcessorConfig) (*ddconfig.AgentConfig, error) {
	conf := ddconfig.NewDefaultAgentConfig()
	conf.ConnectionLimit = c.ConnectionLimit
	conf.ReceiverTimeout = c.ReceiverTimeout
	hostPortSplit := strings.Split(c.HostPort, ":")
	if len(hostPortSplit) != 2 {
		return nil, errors.New("Invalid host-port string")
	}
	conf.ReceiverHost = hostPortSplit[0]
	receiverPort, err := strconv.ParseInt(hostPortSplit[1], 10, 32)
	if err != nil {
		return nil, err
	}
	conf.ReceiverPort = int(receiverPort)
	conf.APIKey = auth.GetToken()
	conf.LogFilePath = logFilePath
	conf.ExtraSampleRate = c.ExtraSampleRate
	conf.MaxTPS = c.MaxTPS
	return conf, nil
}
