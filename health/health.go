// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package health

import (
	"context"
	"fmt"
	"github.com/kodix/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"sync/atomic"
)

// count of current requests
var throttling uint64

// limit of simultaneous requests
var capacity uint64 = 100

// AdditionalFunc string generator for add to health "additional" section
// (returning value must be json object without outer brackets, e.g. "cache":{"len":8})
type AdditionalFunc func() string

// EmptyString is default AdditionalFunc for Health param
func EmptyString() string {
	return ""
}

// SetCapacity store the capacity value (should be used in init section)
func SetCapacity(cap uint64) {
	atomic.StoreUint64(&capacity, cap)
}

// Health returns http handler for health endpoint
func Health(additionalFunc AdditionalFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := []byte(fmt.Sprintf(`{"health":{"len":%d,"cap":%d}, "additional":{%s}}`, atomic.LoadUint64(&throttling), atomic.LoadUint64(&capacity), additionalFunc()))
		_, err := w.Write(b)
		if err != nil {
			log.Debugln(err)
		}
	}
}

// Len is current length of BP
func Len() uint64 {
	return atomic.LoadUint64(&throttling)
}

// Cap is capacity of BP
func Cap() uint64 {
	return atomic.LoadUint64(&capacity)
}

// BackPressure http middleware
func BackPressure(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddUint64(&throttling, 1)
		defer atomic.AddUint64(&throttling, ^uint64(0))
		if current < atomic.LoadUint64(&capacity) {
			h.ServeHTTP(w, r)
		} else {
			log.Debugln("requests limit exceeded")
			http.Error(w, "requests limit exceeded", http.StatusTooManyRequests)
		}
	}
}

//BackPressure grpc middleware
func BackPressureInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		current := atomic.AddUint64(&throttling, 1)
		defer atomic.AddUint64(&throttling, ^uint64(0))
		if current < atomic.LoadUint64(&capacity) {
			return handler(ctx, req)
		}
		log.Debugln("requests limit exceeded")
		return nil, status.Errorf(codes.ResourceExhausted, "requests limit exceeded")
	}
}
