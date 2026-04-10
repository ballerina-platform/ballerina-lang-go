// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"os"
	"time"
)

type (
	benchExport struct {
		Results []benchResult `json:"results"`
	}
	benchResult struct {
		Command string  `json:"command"`
		Mean    float64 `json:"mean"`
		Stddev  float64 `json:"stddev"`
		Median  float64 `json:"median"`
	}
)

func parseHyperfineExport(path string) (*benchExport, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read hyperfine export file: %w", err)
	}

	var export benchExport
	if err := json.Unmarshal(b, &export); err != nil {
		return nil, fmt.Errorf("failed to parse hyperfine export JSON: %w", err)
	}

	return &export, nil
}

type (
	report struct {
		BaseRef   string
		HeadRef   string
		Generated time.Time
		results   []runResult
	}
	row struct {
		Label          string
		Base           *benchResult
		Head           *benchResult
		BaseMean       string
		BaseStddev     string
		HeadMean       string
		HeadStddev     string
		DeltaAvailable bool
		DeltaRatio     string
		DeltaErr       string
		DeltaWinnerRef string
	}
)

func (r *report) export(outPath string) error {
	rows := make([]row, 0, len(r.results))
	for _, run := range r.results {
		var base, head *benchResult
		if len(run.export.Results) >= 2 {
			base = &run.export.Results[0]
			head = &run.export.Results[1]
		}
		row := row{
			Label: run.label,
			Base:  base,
			Head:  head,
		}
		if base != nil {
			row.BaseMean = fmt.Sprintf("%.3f", base.Mean*1000.0)
			row.BaseStddev = fmt.Sprintf("%.3f", base.Stddev*1000.0)
		}
		if head != nil {
			row.HeadMean = fmt.Sprintf("%.3f", head.Mean*1000.0)
			row.HeadStddev = fmt.Sprintf("%.3f", head.Stddev*1000.0)
		}
		row.DeltaAvailable, row.DeltaRatio, row.DeltaErr, row.DeltaWinnerRef = computeDelta(base, head, r.BaseRef, r.HeadRef)
		rows = append(rows, row)
	}

	tpl := template.Must(template.New("report").Parse(htmlTemplate))
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to generate html report %q: %w", outPath, err)
	}
	defer func() { _ = f.Close() }()

	data := struct {
		Report report
		Rows   []row
	}{
		Report: *r,
		Rows:   rows,
	}

	if err := tpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to render html report %q: %w", outPath, err)
	}
	return nil
}

func computeDelta(base, head *benchResult, baseRef, headRef string) (bool, string, string, string) {
	if base == nil || head == nil {
		return false, "", "", ""
	}
	if base.Mean <= 0 || head.Mean <= 0 {
		return false, "", "", ""
	}

	rawRatio := base.Mean / head.Mean
	rawRatioErr := rawRatio * math.Sqrt(
		math.Pow(base.Stddev/base.Mean, 2)+math.Pow(head.Stddev/head.Mean, 2),
	)

	winnerRef := headRef
	ratio := rawRatio
	ratioErr := rawRatioErr
	if rawRatio < 1.0 {
		winnerRef = baseRef
		ratio = 1.0 / rawRatio
		// If q = 1/r then dq = dr / r^2.
		ratioErr = rawRatioErr / (rawRatio * rawRatio)
	}
	return true, fmt.Sprintf("%.2f", ratio), fmt.Sprintf("%.2f", math.Abs(ratioErr)), winnerRef
}
