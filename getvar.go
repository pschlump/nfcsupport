package support

// MIT Licensed - see LICENSE
// Copyright (C) 2017 Philip Schlump

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pschlump/godebug"
)

// GetVar returns a variable by name from GET or POST data.
func GetVar(name string, www http.ResponseWriter, req *http.Request) (found bool, value string) {

	method := req.Method

	godebug.DbPf(dbFlag["GetVar"], "GetVar name=%s req.Method %s AT:%s\n", name, method, godebug.LF())

	if method == "POST" || method == "PUT" {
		if str := req.PostFormValue(name); str != "" { // xyzzy - actually have to check if exists
			value = str
			found = true
		} else {
			if dbFlag["GetVar"] {
				fmt.Printf("AT:%s\n", godebug.LF())
			}
			qq := req.URL.Query()
			strArr, ok := qq[name]
			if dbFlag["GetVar"] {
				fmt.Printf("AT:%s strArr = %s ok = %v\n", godebug.LF(), godebug.SVar(strArr), ok)
			}
			if ok {
				if dbFlag["GetVar"] {
					fmt.Printf("AT:%s\n", godebug.LF())
				}
				if len(strArr) > 0 {
					if dbFlag["GetVar"] {
						fmt.Printf("AT:%s\n", godebug.LF())
					}
					value = strArr[0]
					found = true
				} else {
					if dbFlag["GetVar"] {
						fmt.Printf("AT:%s\n", godebug.LF())
					}
					fmt.Fprintf(os.Stderr, "Multiple values for [%s]\n", name)
					found = false
				}
			}
		}
	} else if method == "GET" || method == "DELETE" {
		godebug.DbPf(dbFlag["GetVar"], "AT:%s\n", godebug.LF())
		qq := req.URL.Query()
		strArr, ok := qq[name]
		godebug.DbPf(dbFlag["GetVar"], "AT:%s strArr = %s ok = %v, parsed Query qq=%s, req.URL=%s\n", godebug.LF(), godebug.SVar(strArr), ok, qq, req.URL)
		if ok {
			if dbFlag["GetVar"] {
				fmt.Printf("AT:%s\n", godebug.LF())
			}
			if len(strArr) > 0 {
				if dbFlag["GetVar"] {
					fmt.Printf("AT:%s\n", godebug.LF())
				}
				value = strArr[0]
				found = true
			} else {
				if dbFlag["GetVar"] {
					fmt.Printf("AT:%s\n", godebug.LF())
				}
				fmt.Fprintf(os.Stderr, "Multiple values for [%s]\n", name)
				found = false
			}
		}
	} else {
		www.WriteHeader(418) // Ha Ha - I Am A Tea Pot
	}
	return
}

func GetMultipleVar(name []string, www http.ResponseWriter, req *http.Request) (found bool, value map[string]string) {
	value = make(map[string]string)
	for _, nn := range name {
		x, v := GetVar(nn, www, req)
		if !x {
			www.WriteHeader(406) // not acceptable
			fmt.Fprintf(www, `{"status":"error", "msg":"missing %s parameter - required"}`, nn)
			return
		}
		value[nn] = v
	}
	found = true
	return
}
