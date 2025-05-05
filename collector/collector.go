// SPDX-FileCopyrightText: 2025 Sidings Media <contact@sidingsmedia.com>
// SPDX-License-Identifier: MIT

package collector

import (
	"log/slog"
	"sync"
	"time"

	"github.com/SidingsMedia/icmp_exporter/config"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector collects all Collector metrics, implemented as per the prometheus.Collector interface.
type Collector struct {
	logger       *slog.Logger
	descriptions map[string]*prometheus.Desc
	config       *config.Config
}

// NewCollector returns a new Exporter.
func NewCollector(logger *slog.Logger, config *config.Config) (*Collector, error) {
	labels := []string{"host", "interface"}
	descriptions := map[string]*prometheus.Desc{
		"icmp_avg_rtt":      prometheus.NewDesc("icmp_avg_rtt", "Average round trip time to the target.", labels, nil),
		"icmp_min_rtt":      prometheus.NewDesc("icmp_min_rtt", "Minimum round trip time to the target.", labels, nil),
		"icmp_max_rtt":      prometheus.NewDesc("icmp_max_rtt", "Maximum round trip time to the target.", labels, nil),
		"icmp_std_dev_rtt":  prometheus.NewDesc("icmp_std_dev_rtt", "Standard deviation of round trip time to the target.", labels, nil),
		"icmp_packets_sent": prometheus.NewDesc("icmp_packets_sent", "Number of packets sent in this run.", labels, nil),
		"icmp_packets_recv": prometheus.NewDesc("icmp_packets_recv", "Number of packets received in this run.", labels, nil),
		"icmp_packet_loss": prometheus.NewDesc("icmp_packet_loss", "Percentage of packets lost in this run.", labels, nil),
	}

	return &Collector{
		logger:       logger,
		descriptions: descriptions,
		config:       config,
	}, nil
}

// Collect implemented as per the prometheus.Collector interface.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
    var wg sync.WaitGroup

	wg.Add(len(c.config.Targets))

    for _, target := range c.config.Targets {
        go func() {
            defer wg.Done()
          stats, err := c.ping(&target)
          if err != nil {
            c.logger.Error("Failed to ping target", "target", target.Host, "err", err)
            return
          }

          labels := []string{target.Host, target.Interface}

          ch <- prometheus.MustNewConstMetric(c.descriptions["icmp_avg_rtt"], prometheus.GaugeValue, stats.AvgRtt.Seconds(), labels...)
          ch <- prometheus.MustNewConstMetric(c.descriptions["icmp_min_rtt"], prometheus.GaugeValue, stats.MinRtt.Seconds(), labels...)
          ch <- prometheus.MustNewConstMetric(c.descriptions["icmp_max_rtt"], prometheus.GaugeValue, stats.MaxRtt.Seconds(), labels...)
          ch <- prometheus.MustNewConstMetric(c.descriptions["icmp_std_dev_rtt"], prometheus.GaugeValue, stats.StdDevRtt.Seconds(), labels...)
          ch <- prometheus.MustNewConstMetric(c.descriptions["icmp_packets_sent"], prometheus.GaugeValue, float64(stats.PacketsSent), labels...)
          ch <- prometheus.MustNewConstMetric(c.descriptions["icmp_packets_recv"], prometheus.GaugeValue, float64(stats.PacketsRecv), labels...)
          ch <- prometheus.MustNewConstMetric(c.descriptions["icmp_packet_loss"], prometheus.GaugeValue, float64(stats.PacketLoss), labels...)
        }()
    }

    wg.Wait()
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
    for _, desc := range c.descriptions {
		ch <- desc
	}
}

// Ping the target host
func (c *Collector) ping(target *config.Target) (*probing.Statistics, error) {
	pinger := probing.New(target.Host)

	pinger.Count = target.Count
	pinger.Interval = time.Duration(target.Interval * int(time.Millisecond))
	pinger.Size = target.Size
	pinger.TTL = target.TTL
	pinger.Timeout = time.Duration(c.config.Timeout * int(time.Millisecond))

    if (target.Interface != "") {
        pinger.InterfaceName = target.Interface
    }

	pinger.SetPrivileged(true)

	c.logger.Debug("Starting ping", "pinger", pinger)

	if err := pinger.Run(); err != nil {
		c.logger.Error("Failed to ping host", "err", err)
		return nil, err
	}

	c.logger.Debug("Finished ping", "target", target.Host, "statistics", pinger.Statistics())

	return pinger.Statistics(), nil
}
