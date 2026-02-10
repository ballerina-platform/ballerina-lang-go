//go:build debug

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
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Auto-registers pprof handlers
	"os"
	"time"

	"github.com/spf13/cobra"
)

type enabledProfiler struct {
	server  *http.Server
	enabled bool
	addr    string
}

func init() {
	profiler = &enabledProfiler{}
}

func (p *enabledProfiler) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&p.enabled, "prof", false, "Enable profiling")
	cmd.Flags().StringVar(&p.addr, "prof-addr", ":6060", "Profiling server address")
}

func (p *enabledProfiler) Start() error {
	if !p.enabled {
		return nil
	}

	// Create HTTP server (pprof handlers auto-registered by import)
	p.server = &http.Server{
		Addr:         p.addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		fmt.Fprintf(os.Stderr, "Profiling enabled at http://%s/debug/pprof/\n", p.addr)
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Profiling server error: %v\n", err)
		}
	}()

	// Brief delay to catch immediate startup errors (port conflicts)
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (p *enabledProfiler) Stop() error {
	if p.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.server.Shutdown(ctx)
}
