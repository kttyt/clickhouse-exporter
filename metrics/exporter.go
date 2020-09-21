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

package metrics

import (
	"sync"

	log "github.com/golang/glog"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter implements prometheus.Collector interface
type Exporter struct {
	chAccessInfo *CHAccessInfo

	mutex               sync.RWMutex
	toRemoveFromWatched sync.Map
}

var exporter *Exporter

// NewExporter returns a new instance of Exporter type
func NewExporter(chAccess *CHAccessInfo) *Exporter {
	return &Exporter{
		// chInstallations: make(map[string]*WatchedCHI),
		chAccessInfo: chAccess,
	}
}

// Collect implements prometheus.Collector Collect method
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if ch == nil {
		log.V(2).Info("Prometheus channel is closed. Skipping")
		return
	}

	log.V(2).Info("Starting Collect")
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func(c chan<- prometheus.Metric) {
		defer wg.Done()
		e.collectFromHost(c)
	}(ch)
	wg.Wait()
	log.V(2).Info("Finished Collect")
}

// Describe implements prometheus.Collector Describe method
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(e, ch)
}

// newFetcher returns new Metrics Fetcher for specified host
func (e *Exporter) newFetcher() *ClickHouseFetcher {
	return NewClickHouseFetcher(e.chAccessInfo.Hostname, e.chAccessInfo.Username, e.chAccessInfo.Password, e.chAccessInfo.Port)
}

// collectFromHost collects metrics from one host and writes them into chan
func (e *Exporter) collectFromHost(c chan<- prometheus.Metric) {
	fetcher := e.newFetcher()
	writer := NewPrometheusWriter(c, e.chAccessInfo.Hostname)

	log.V(2).Infof("Querying metrics for %s\n", e.chAccessInfo.Hostname)
	if metrics, err := fetcher.getClickHouseQueryMetrics(); err == nil {
		log.V(2).Infof("Extracted %d metrics for %s\n", len(metrics), e.chAccessInfo.Hostname)
		writer.WriteMetrics(metrics)
		writer.WriteOKFetch("system.metrics")
	} else {
		// In case of an error fetching data from clickhouse store CHI name in e.cleanup
		log.V(2).Infof("Error querying metrics for %s: %s\n", e.chAccessInfo.Hostname, err)
		writer.WriteErrorFetch("system.metrics")
		return
	}

	log.V(2).Infof("Querying table sizes for %s\n", e.chAccessInfo.Hostname)
	if tableSizes, err := fetcher.getClickHouseQueryTableSizes(); err == nil {
		log.V(2).Infof("Extracted %d table sizes for %s\n", len(tableSizes), e.chAccessInfo.Hostname)
		writer.WriteTableSizes(tableSizes)
		writer.WriteOKFetch("table sizes")
	} else {
		// In case of an error fetching data from clickhouse store CHI name in e.cleanup
		log.V(2).Infof("Error querying table sizes for %s: %s\n", e.chAccessInfo.Hostname, err)
		writer.WriteErrorFetch("table sizes")
		return
	}

	log.V(2).Infof("Querying system replicas for %s\n", e.chAccessInfo.Hostname)
	if systemReplicas, err := fetcher.getClickHouseQuerySystemReplicas(); err == nil {
		log.V(2).Infof("Extracted %d system replicas for %s\n", len(systemReplicas), e.chAccessInfo.Hostname)
		writer.WriteSystemReplicas(systemReplicas)
		writer.WriteOKFetch("system.replicas")
	} else {
		// In case of an error fetching data from clickhouse store CHI name in e.cleanup
		log.V(2).Infof("Error querying system replicas for %s: %s\n", e.chAccessInfo.Hostname, err)
		writer.WriteErrorFetch("system.replicas")
		return
	}

	log.V(2).Infof("Querying mutations for %s\n", e.chAccessInfo.Hostname)
	if mutations, err := fetcher.getClickHouseQueryMutations(); err == nil {
		log.V(2).Infof("Extracted %d mutations for %s\n", len(mutations), e.chAccessInfo.Hostname)
		writer.WriteMutations(mutations)
		writer.WriteOKFetch("system.mutations")
	} else {
		// In case of an error fetching data from clickhouse store CHI name in e.cleanup
		log.V(2).Infof("Error querying mutations for %s: %s\n", e.chAccessInfo.Hostname, err)
		writer.WriteErrorFetch("system.mutations")
		return
	}
}
