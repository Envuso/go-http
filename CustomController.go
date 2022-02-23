package main

import (
	"github.com/envuso/go-http/HttpContext"
	"github.com/envuso/go-http/Routing/Route"
)

type CustomController struct {
	Route.Controller

	SomeBsService SomeBullShitServiceContract

	SomeOtherBs SomeOtherBullshitServiceContract
}

func (c *CustomController) UsernameRoute(ctx *HttpContext.RequestContext, data *BigYeetsRequest) map[string]interface{} {
	username := data.Get("data.user.username").String()

	return map[string]interface{}{
		"name":   ctx.Params().Get("name"),
		"second": username, // data.Get("username").String(),
		"data":   data.Json().Value(),
	}
}

func (c *CustomController) DifferentResTypesPog(ctx *HttpContext.RequestContext, req *HttpContext.Request, response *HttpContext.Response) *HttpContext.Response {
	return response.Json(map[string]interface{}{"message": "pogu"})
}
func (c *CustomController) PureMapAids(ctx *HttpContext.RequestContext, req *HttpContext.Request) map[string]interface{} {
	return map[string]interface{}{
		"ctx.idString()": ctx.IdString(),
		"req.IdString()": req.IdString(),
	}
}

type FFF struct {
	Message string `json:"message,omitempty"`
}

func (c *CustomController) StructRes(ctx *HttpContext.RequestContext, req *HttpContext.Request, response *HttpContext.Response) FFF {
	return FFF{Message: "big yeet?"}
}
