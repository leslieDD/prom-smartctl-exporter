package exporter

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var _ prometheus.Collector = &Collector{}

type Collector struct {
	device       string
	PowerOnHours *prometheus.Desc
	Temperature  *prometheus.Desc
}

func NewCollector(device string) *Collector {
	var (
		labels = []string{
			"device",
			"model",
		}
	)
	return &Collector{
		device: device,
		PowerOnHours: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "power_on_hours"),
			"Power on hours",
			labels,
			nil,
		),
		Temperature: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "temperature"),
			"Temperature",
			labels,
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.PowerOnHours,
		c.Temperature,
	}
	for _, d := range ds {
		ch <- d
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	// if desc, err := c.collect(ch); err != nil {
	// 	log.Printf("[ERROR] failed collecting metric %v: %v", desc, err)
	// 	ch <- prometheus.NewInvalidMetric(desc, err)
	// 	return
	// }
	_, _ = c.collect(ch)
}

func (c *Collector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	if c.device == "" {
		return nil, nil
	}

	// out, err := exec.Command("smartctl", "-iA", c.device).CombinedOutput()
	code, out, err := ExecCmd("smartctl", "-iA", c.device)
	if err != nil {
		log.Printf("[ERROR] smart log: \n%s\n device: %s", out, c.device)
		return nil, err
	}
	if code != 0 {
		log.Printf("[ERROR] code: %d, device: %s, smart log: \n%s\n", code, c.device, out)
		return nil, err
	}
	smart := ParseSmart(string(out))

	if len(smart.attrs) == 0 || len(smart.info) == 0 {
		log.Printf("[ERROR] attrs or info is empty\n%s\n", out)
		return nil, err
	}

	labels := []string{
		c.device,
		smart.GetInfo("Device Model", "Model Family"),
	}

	ch <- prometheus.MustNewConstMetric(
		c.PowerOnHours,
		prometheus.GaugeValue,
		float64(smart.GetAttr(9).rawValue),
		labels...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.Temperature,
		prometheus.GaugeValue,
		float64(smart.GetAttr(190, 194).rawValue),
		labels...,
	)

	return nil, nil
}

// 返回结果状态码，终端【标准输出，错误输出】输出，错误信息
func ExecCmd(command string, args ...string) (int, string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	var res int
	if err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			res = ex.Sys().(syscall.WaitStatus).ExitStatus()
		}
		return res, "", err
	}
	if utf8Output, err := GbkToUtf8(output); err != nil {
		return res, string(output), nil
	} else {
		return res, string(utf8Output), nil
	}
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
