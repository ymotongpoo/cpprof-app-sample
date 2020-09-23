// Copyright 2020 Yoshi Yamaguchi
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

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/rs/zerolog"
)

const (
	Version                = "0.3.0"
	Service                = "cpprofhistory"
	BlockSimulatedDuration = 1 * time.Second
)

var logger zerolog.Logger

func levelFieldMarshalFunc(l zerolog.Level) string {
	switch l {
	case zerolog.TraceLevel:
		return "DEFAULT"
	case zerolog.DebugLevel:
		return "DEBUG"
	case zerolog.InfoLevel:
		return "INFO"
	case zerolog.WarnLevel:
		return "WARNING"
	case zerolog.ErrorLevel:
		return "ERROR"
	case zerolog.FatalLevel:
		return "CRITICAL"
	case zerolog.PanicLevel:
		return "ALERT"
	default:
		return "DEFAULT"
	}
}

func init() {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelFieldName = "severity"
	zerolog.LevelFieldMarshalFunc = levelFieldMarshalFunc
	logger = zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger()
}

func main() {
	logger.Info().Msgf("%s: starting server...", Service)
	go startProfiler()

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info().Msgf("%s: listening on port %s", Service, port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		logger.Error().Msgf("error on running HTTP server: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger.Debug().Msgf("%s: received a request", Service)
	time.Sleep(BlockSimulatedDuration)

	m, _ := url.ParseQuery(r.URL.RawQuery)
	target := "world"
	if v, ok := m["q"]; ok {
		target = v[0]
	}
	fmt.Fprintf(w, "Hello %s\n", target)
}

func startProfiler() {
	cfg := profiler.Config{
		Service:        Service,
		ServiceVersion: Version,
	}
	if err := profiler.Start(cfg); err != nil {
		logger.Error().Msgf("error on starting profiler: %v", err)
	}
}
