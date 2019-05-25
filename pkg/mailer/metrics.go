/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type meters struct {
	ReqCount    *kitprometheus.Counter
	ReqLatency  *kitprometheus.Summary
	CountResult *kitprometheus.Summary
}

// Instrumentation
func instrumentationMeters() meters {
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "poslan",
		Subsystem: "auth",
		Name:      "request_count",
		Help:      "Nº of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "poslan",
		Subsystem: "auth",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in μSeconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "poslan",
		Subsystem: "auth",
		Name:      "count_result",
		Help:      "Sumary.",
	}, []string{}) // no fields here

	return meters{
		ReqCount:    requestCount,
		ReqLatency:  requestLatency,
		CountResult: countResult,
	}
}
