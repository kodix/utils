// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmptyString(t *testing.T) {
	if "" != EmptyString() {
		t.Errorf("EmptyString() is not empty")
	}
}

func TestSetCapacity(t *testing.T) {
	SetCapacity(50)
	if capacity != 50 {
		t.Error("SetCapacity() - has not been set")
	}
	SetCapacity(100)
	if capacity != 100 {
		t.Error("SetCapacity() - has not been set")
	}
}

func TestLen(t *testing.T) {
	throttling = 50
	if Len() != 50 {
		t.Error("Len() error")
	}
	throttling = 0
}

func TestCap(t *testing.T) {
	capacity = 50
	if Cap() != 50 {
		t.Error("Cap() error")
	}
	capacity = 100
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`London is the capital of GB`))
}

func TestBackPressure(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	h := BackPressure(http.HandlerFunc(testHandler))
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("BackPressure() - want status 200, got %d", w.Code)
	}

	throttling = capacity
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("BackPressure() - want status 429, got %d", w.Code)
	}
	throttling = 0
}

func TestHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	h := http.HandlerFunc(Health(EmptyString))
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Health() - want status 200, got status %d", w.Code)
	}
	if w.Body.String() != `{"health":{"len":0,"cap":100}, "additional":{}}` {
		t.Errorf("Health() - got response body %s", w.Body.String())
	}
}
