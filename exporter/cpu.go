package exporter

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

var _ prometheus.Collector = &Collector{}

type CollectCpuInfo struct {
	device    string
	ModelName *prometheus.Desc
}

func NewCollectCpuInfo() *CollectCpuInfo {
	var (
		labels = []string{
			"modename",
		}
	)
	return &CollectCpuInfo{
		ModelName: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "modename"),
			"cpu mode name",
			labels,
			nil,
		),
	}
}

func (c *CollectCpuInfo) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.ModelName,
	}
	for _, d := range ds {
		ch <- d
	}
}

func (c *CollectCpuInfo) Collect(ch chan<- prometheus.Metric) {
	_, _ = c.collect(ch)
}

func (c *CollectCpuInfo) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	if c.device == "" {
		return nil, nil
	}

	info := GetCpuInfo()
	if info == nil || len(info) == 0 {
		return nil, nil
	}

	for _, each := range info {
		labels := []string{
			fmt.Sprintf("core_%s", each.CoreID),
			each.ModelName,
		}

		ch <- prometheus.MustNewConstMetric(
			c.ModelName,
			prometheus.GaugeValue,
			float64(1),
			labels...,
		)
	}

	return nil, nil
}

func GetCpuInfo() []cpu.InfoStat {
	info, err := cpu.Info()
	if err != nil {
		log.Printf("[ERROR] failed collecting cpu metric: %v", err)
		return nil
	}
	return info
}
