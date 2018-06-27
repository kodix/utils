// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mw

import (
	"net/http"
	"github.com/abramd/log"
	"fmt"
	"os"
	"flag"
)

// Matcher checks access of request
type Matcher interface {
	Match(*http.Request) bool
}

// MatcherFunc is func implementation of Matcher
type MatcherFunc func(*http.Request) bool

// Match is MatcherFunc method
func (m MatcherFunc) Match(r *http.Request) bool {
	return m(r)
}

// KeyFunc returns header keys for match
type KeyFunc func() string

// HeaderMatcher matches existence of the given header value
func HeaderMatcher(kf KeyFunc, input string) MatcherFunc {
	return func(r *http.Request) bool {
		if h, ok := r.Header[kf()]; ok {
			for _, v := range h {
				if input == v {
					return true
				}
			}
		}
		return false
	}
}

// RoleKeyFunc returns the func, which returns the X-Auth-Access-Roles header key with service ID
func RoleKeyFunc(cid string) KeyFunc {
	return func() string {
		return fmt.Sprintf("X-Auth-Access-Roles-%s", http.CanonicalHeaderKey(cid))
	}
}

// AuthMw check the matcher value. If true - refers to next http.Handler, if false - returns 403 http status
func AuthMw(next http.Handler, m Matcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.Match(r) {
			next.ServeHTTP(w, r)
		} else {
			log.Infoln("forbidden")
			http.Error(w, "forbidden", http.StatusForbidden)
		}
	}
}

// RoleAuthMw checks the X-Auth-Access-Roles header with given role
func RoleAuthMw(next http.Handler, allowedRole string) http.HandlerFunc {
	return AuthMw(next, HeaderMatcher(RoleKeyFunc(clientID), allowedRole))
}

// DefAuthMw checks X-Auth-Access-Roles header with default role
func DefAuthMw(next http.Handler) http.HandlerFunc {
	return RoleAuthMw(next, allowedRole)
}

// Default role and service ID variables
var (
	allowedRole = ""
	clientID    = ""
)

// SetAllowedRole
func SetAllowedRole(role string) {
	allowedRole = role
}

// SetClientID
func SetClientID(id string) {
	clientID = id
}

// init - initialization of default variables
func init() {
	flag.StringVar(&allowedRole, "role", os.Getenv("ALLOWED_ROLE"), "x-auth-access role for auth")
	flag.StringVar(&clientID, "cid", os.Getenv("CLIENT_ID"), "service id")
}
