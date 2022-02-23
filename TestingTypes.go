package main

import "github.com/envuso/go-http/HttpContext"

type BigYeetsRequest struct {
	HttpContext.RequestBody

	Username string                 `json:"username"`
	Data     map[string]interface{} `json:"data"`
}
