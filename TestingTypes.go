package main

import "gohttp/HttpContext"

type BigYeetsRequest struct {
	HttpContext.RequestBody

	Username string                 `json:"username"`
	Data     map[string]interface{} `json:"data"`
}
