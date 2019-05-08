// Copyright 2019 The Prometheus Authors
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

// +build !loginsession

package collector

import (
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type loginSessionCollector struct {
	loginSession *prometheus.Desc
}

func init() {
	registerCollector("loginsession", defaultEnabled, newLoginSessionCollector)
}

func newLoginSessionCollector() (Collector, error) {
	ls := &loginSessionCollector{
		loginSession: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "login_session"),
			"login session from w(who) command",
			[]string{"source"}, nil,
		),
	}
	return ls, nil
}

func (ls *loginSessionCollector) Update(ch chan<- prometheus.Metric) error {
	lines, err := ls.getSessions()
	if err != nil {
		return err
	}
	for _, line := range lines {
		ch <- prometheus.MustNewConstMetric(ls.loginSession, prometheus.GaugeValue, 1, line)
	}
	return nil
}

func (ls *loginSessionCollector) getSessions() ([]string, error) {
	cmd := "w|awk '{print $3}'"
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	retLines := make([]string, 0)
	for i, line := range lines {
		if i <= 1 {
			continue
		}
		if len(line) < 1 {
			continue
		}
		if line == "-" {
			continue
		}
		retLines = append(retLines, line)
	}
	return retLines, nil
}
