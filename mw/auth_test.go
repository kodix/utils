// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mw

import (
	"net/http"
	"testing"
	"net/http/httptest"
)

func next(http.ResponseWriter, *http.Request) {

}

func TestDefAuthMw(t *testing.T) {
	// success
	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	SetAllowedRole("root")
	SetClientID("tester")
	r.Header.Set("X-Auth-Access-Roles-Tester", "root")
	DefAuthMw(http.HandlerFunc(next)).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("want code 200, got %d", w.Code)
	}

	// forbidden
	w = httptest.NewRecorder()
	r, err = http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("X-Auth-Access-Roles-Tester", "anonymous")
	DefAuthMw(http.HandlerFunc(next)).ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("want code 200, got %d", w.Code)
	}
	SetAllowedRole("")
	SetClientID("")
}
