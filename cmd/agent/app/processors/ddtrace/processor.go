package ddtrace

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/DataDog/datadog-trace-agent/cmd/ddtrace"
	ddconfig "github.com/DataDog/datadog-trace-agent/config"
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
	LogLevel        string  `yaml:"logLevel"`
}

// Processor is a collector that uses HTTP protocol and just holds
// a chan where the spans received are sent one by one
type Processor struct {
	agent *ddAgent

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

	err = ddtrace.SetupLoggerFromAgentConfig(ddAgentConfig)
	if err != nil {
		return nil, err
	}

	ddtraceProcessor := Processor{
		agent:      newDDAgent(ddAgentConfig),
		conf:       c,
		isRunning:  false,
		stopCalled: false,
		Mutex:      &sync.Mutex{},
	}
	auth.AddTokenUpdateAction(ddtraceProcessor.restartOnTokenChange)

	return &ddtraceProcessor, nil
}

// Serve starts doing the HTTP server and is ready to receive traces
func (r *Processor) Serve() {
	r.Lock()
	defer r.Unlock()
	// Start the agent if it is not already running
	if !r.isRunning {
		go r.agent.start()
		r.isRunning = true
		r.stopCalled = false
	}
}

// Stop stops the processor
func (r *Processor) Stop() {
	r.Lock()
	r.isRunning = false
	r.stopCalled = true
	r.agent.stop()
	r.Unlock()
}

func (r *Processor) restartOnTokenChange(logger *zap.Logger) {
	r.Lock()
	defer r.Unlock()
	if r.stopCalled == true {
		return
	}

	ddAgentConfig, err := newDDAgentConfig(r.conf)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Make sure that the agent is stopped before starting a new one
	r.agent.stop()
	logger.Info("Stopped ddagent due to token update")

	r.agent = newDDAgent(ddAgentConfig)
	go r.agent.start()
	r.isRunning = true
	logger.Info("Started ddagent due to token update")
}

func newDDAgentConfig(c ProcessorConfig) (*ddconfig.AgentConfig, error) {
	conf := ddconfig.New()
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
	conf.LogLevel = c.LogLevel
	conf.LogFilePath = logFilePath
	conf.ExtraSampleRate = c.ExtraSampleRate
	conf.MaxTPS = c.MaxTPS
	return conf, nil
}
