// SPDX-FileCopyrightText: 2025 2025 Sidings Media <contact@sidingsmedia.com>
// SPDX-License-Identifier: MIT

package main

import (
	"net/http"
	"os"

	"github.com/SidingsMedia/icmp_exporter/config"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

var (
	telemetryPath = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	webFlagConfig = kingpinflag.AddFlags(kingpin.CommandLine, ":9342")
    configFile = kingpin.Flag("collector.config.file", "Path to exporter configuration.").Required().String()
)

func main() {
	promslogConfig := &promslog.Config{}

	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("icmp_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promslog.New(promslogConfig)

    logger.Info("Starting icmp_exporter", "version", version.Info())
	logger.Info("Build context", "build_context", version.BuildContext())

    _, err := config.ParseConfig(*configFile, logger)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

	prometheus.MustRegister(versioncollector.NewCollector("icmp_exporter"))

    http.Handle(*telemetryPath, promhttp.Handler())
	if *telemetryPath != "/" && *telemetryPath != "" {
		landingConfig := web.LandingConfig{
			Name:        "ICMP Exporter",
			Description: "Prometheus Exporter for monitoring point to point link latency",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{Address: *telemetryPath, Text: "Metrics"},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	server := &http.Server{}
	if err := web.ListenAndServe(server, webFlagConfig, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
