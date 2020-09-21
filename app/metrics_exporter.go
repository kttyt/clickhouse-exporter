// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/golang/glog"
	"github.com/jamiealquiza/envy"

	"clickhouse-exporter/metrics"
	"clickhouse-exporter/version"
)

// Prometheus exporter defaults
const (
	defaultMetricsEndpoint = ":8888"
	metricsPath            = "/metrics"
	defaultCHHostname      = "localhost"
	defaultCHPort          = 8123
)

// CLI parameter variables
var (
	// versionRequest defines request for clickhouse-exporter version report. Operator should exit after version printed
	versionRequest bool

	// metricsEP defines metrics end-point IP address
	metricsEP string
	// Username for getting metrics from ClickHouse
	CHUsername string
	// Password of user for getting metrics from ClickHouse
	CHPassword string
	// ClickHouse hostname
	CHHostname string
	// ClickHouse port
	CHPort int
)

func init() {
	flag.BoolVar(&versionRequest, "version", false, "Display clickhouse-exporter version and exit")
	flag.StringVar(&metricsEP, "metrics-endpoint", defaultMetricsEndpoint, "The Prometheus exporter endpoint.")
	flag.StringVar(&CHHostname, "hostname", defaultCHHostname, "Hostname of CK installation.")
	flag.StringVar(&CHUsername, "username", "", "Username for getting metrics.")
	flag.StringVar(&CHPassword, "password", "", "Password of user for getting metrics.")
	flag.IntVar(&CHPort, "port", defaultCHPort, "The Clickhouse http port.")
	envy.Parse("CH_EXPORTER")
	flag.Parse()
}

// Run is an entry point of the application
func Run() {
	if versionRequest {
		fmt.Printf("%s\n", version.Version)
		os.Exit(0)
	}

	// Set OS signals and termination context
	ctx, cancelFunc := context.WithCancel(context.Background())
	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		cancelFunc()
		<-stopChan
		os.Exit(1)
	}()

	log.Infof("Starting metrics exporter. Version:%s GitSHA:%s BuiltAt:%s\n", version.Version, version.GitSHA, version.BuiltAt)

	metrics.StartMetricsREST(
		metrics.NewCHAccessInfo(
			CHUsername,
			CHPassword,
			CHHostname,
			CHPort,
		),

		metricsEP,
		metricsPath,
	)

	<-ctx.Done()
}
