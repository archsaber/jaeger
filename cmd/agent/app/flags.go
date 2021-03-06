// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	suffixWorkers             = "workers"
	suffixServerQueueSize     = "server-queue-size"
	suffixServerMaxPacketSize = "server-max-packet-size"
	suffixServerHostPort      = "server-host-port"
	collectorHostPort         = "collector.host-port"
	httpServerHostPort        = "http-server.host-port"
	ddServerEnabled           = "dd-server.enabled"
	ddServerWorkers           = "dd-server.workers"
	ddServerHostPort          = "dd-server.host-port"
	ddServerConnLimit         = "dd-server.conn-limit"
	ddServerReceiverTimeout   = "dd-server.receiver-timeout"
	discoveryMinPeers         = "discovery.min-peers"
)

var defaultProcessors = []struct {
	model    Model
	protocol Protocol
	hostPort string
}{
	{model: "zipkin", protocol: "compact", hostPort: ":5775"},
	{model: "jaeger", protocol: "compact", hostPort: ":6831"},
	{model: "jaeger", protocol: "binary", hostPort: ":6832"},
}

// AddFlags adds flags for Builder.
func AddFlags(flags *flag.FlagSet) {
	for _, processor := range defaultProcessors {
		prefix := fmt.Sprintf("processor.%s-%s.", processor.model, processor.protocol)
		flags.Int(prefix+suffixWorkers, defaultServerWorkers, "how many workers the processor should run")
		flags.Int(prefix+suffixServerQueueSize, defaultQueueSize, "length of the queue for the UDP server")
		flags.Int(prefix+suffixServerMaxPacketSize, defaultMaxPacketSize, "max packet size for the UDP server")
		flags.String(prefix+suffixServerHostPort, processor.hostPort, "host:port for the UDP server")
	}
	flags.String(
		collectorHostPort,
		"",
		"comma-separated string representing host:ports of a static list of collectors to connect to directly (e.g. when not using service discovery)")
	flags.String(
		httpServerHostPort,
		defaultHTTPServerHostPort,
		"host:port of the http server (e.g. for /sampling point and /baggage endpoint)")
	flags.Int(
		discoveryMinPeers,
		defaultMinPeers,
		"if using service discovery, the min number of connections to maintain to the backend")
	flags.String(
		ddServerHostPort,
		defaultDDServerHostPort,
		"host:port of the server to receive traces from datadog agents")
	flags.Int(
		ddServerConnLimit,
		defaultDDConnLimit,
		"connection limit of the dd-trace processor receiving traces")
	flags.Int(
		ddServerReceiverTimeout,
		defaultDDReceiverTimeout,
		"receiver timeout of the dd-trace processor receiving traces")
	flags.Bool(
		ddServerEnabled,
		defaultDDServerEnabled,
		"whether dd-trace processor is enabled")
	flags.Int(
		ddServerWorkers,
		defaultDDServerWorkers,
		"how many workers the dd-trace processor should run")
}

// InitFromViper initializes Builder with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.Metrics.InitFromViper(v)

	for _, processor := range defaultProcessors {
		prefix := fmt.Sprintf("processor.%s-%s.", processor.model, processor.protocol)
		p := &ProcessorConfiguration{Model: processor.model, Protocol: processor.protocol}
		p.Workers = v.GetInt(prefix + suffixWorkers)
		p.Server.QueueSize = v.GetInt(prefix + suffixServerQueueSize)
		p.Server.MaxPacketSize = v.GetInt(prefix + suffixServerMaxPacketSize)
		p.Server.HostPort = v.GetString(prefix + suffixServerHostPort)
		b.Processors = append(b.Processors, *p)
	}

	if len(v.GetString(collectorHostPort)) > 0 {
		b.CollectorHostPorts = strings.Split(v.GetString(collectorHostPort), ",")
	}
	b.HTTPServer.HostPort = v.GetString(httpServerHostPort)

	b.DDTraceProcessorConfig.Enabled = v.GetBool(ddServerEnabled)
	b.DDTraceProcessorConfig.NumProcessors = v.GetInt(ddServerWorkers)
	b.DDTraceProcessorConfig.HostPort = v.GetString(ddServerHostPort)
	b.DDTraceProcessorConfig.ConnectionLimit = v.GetInt(ddServerConnLimit)
	b.DDTraceProcessorConfig.ReceiverTimeout = v.GetInt(ddServerReceiverTimeout)

	b.DiscoveryMinPeers = v.GetInt(discoveryMinPeers)
	return b
}
