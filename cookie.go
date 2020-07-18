package support

// MIT Licensed - see LICENSE
// Copyright (C) 2016-2018 Philip Schlump

import (
	"fmt"
	"net/http"
	"time"
)

func HaveCookie(www http.ResponseWriter, req *http.Request, cookieName string) bool {
	for _, cookie := range req.Cookies() {
		fmt.Println("Found a cookie named:", cookie.Name)
		if cookie.Name == cookieName {
			fmt.Println("!!! Match !!! Found a cookie named:", cookie.Name)
			return true
		}
	}
	return false
}

// addCookie will apply a new cookie to the response of a http
// request, with the key/value this method is passed.
func AddCookie(w http.ResponseWriter, name, value string, inDays int) {
	expire := time.Now().AddDate(0, 0, inDays)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func AddSecureCookie(w http.ResponseWriter, name, value string, inDays int) {
	// See: https://www.calhoun.io/securing-cookies-in-go/
	// May need to add Domain, Path
	expire := time.Now().AddDate(0, 0, inDays)
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expire,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func SetCookie(www http.ResponseWriter, req *http.Request, cookieName, value string, expireInDays int) {
	AddCookie(www, cookieName, value, expireInDays)
}

func SetSecureCookie(www http.ResponseWriter, req *http.Request, cookieName, value string, expireInDays int) {
	// if https: then use a secure cookie! IsTLS(req)
	if IsTLS(req) {
		AddSecureCookie(www, cookieName, value, expireInDays)
	} else {
		AddCookie(www, cookieName, value, expireInDays)
	}
}

func GetCookie(www http.ResponseWriter, req *http.Request, cookieName string) (val string) {
	Ck := req.Cookies()
	for _, v := range Ck {
		if v.Name == cookieName {
			val = v.Value
			return
		}
	}
	return ""
}

/* vim: set noai ts=4 sw=4: */
