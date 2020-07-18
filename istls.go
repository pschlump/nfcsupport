package support

import "net/http"

// MIT Licensed - see LICENSE
// Copyright (C) 2013 Philip Schlump

func IsTLS(req *http.Request) bool {
	if req.TLS == nil {
		return false
	}
	return true
}
