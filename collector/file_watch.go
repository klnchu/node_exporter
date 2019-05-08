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

// +build !filewath

package collector

import (
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	files = kingpin.Flag("collector.file.watch.path", "Directory to read files with metrics from.").Default("").String()
)

type fileWatchCollector struct {
	files []string
	fs    *prometheus.Desc
}

func init() {
	registerCollector("filewath", defaultEnabled, NewFileWatchCollector)
}

// NewFileWatchCollector returns a new Collector exposing metrics read from files
func NewFileWatchCollector() (Collector, error) {

	fs := strings.Split(*files, ",")
	fw := &fileWatchCollector{
		files: fs,
		fs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "file_watch"),
			"watch file watch",
			[]string{"filename"}, nil,
		),
	}
	return fw, nil
}

func (fw *fileWatchCollector) Update(ch chan<- prometheus.Metric) error {

	for _, f := range fw.files {
		file, err := os.Open(f)
		if err != nil {
			continue
		}
		stat, err := file.Stat()
		if err != nil {
			continue
		}
		modTime := float64(stat.ModTime().Unix())
		labels := []string{f}
		ch <- prometheus.MustNewConstMetric(fw.fs, prometheus.GaugeValue, modTime, labels...)
	}
	return nil
}
