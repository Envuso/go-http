package Route

import "github.com/envuso/go-http/HttpContext"

type Controller struct {
	Context *HttpContext.RequestContext
}

func (c *Controller) SetContext(ctx interface{}) {
	c.Context = ctx.(*HttpContext.RequestContext)
}
