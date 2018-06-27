// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mw

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kodix.ru/utils/health"
)

// MetricsMw add http duration metric and BP metric
func MetricsMw(dur prometheus.ObserverVec, bp prometheus.Observer, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.InstrumentHandlerDuration(dur.MustCurryWith(prometheus.Labels{"endpoint": r.URL.Path}), next).ServeHTTP(w,r)
		bp.Observe(float64(health.Len()) / float64(health.Cap()))
	}
}
