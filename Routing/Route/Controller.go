package Route

import "gohttp/HttpContext"

type Controller struct {
	Context *HttpContext.RequestContext
}

func (c *Controller) SetContext(ctx interface{}) {
	c.Context = ctx.(*HttpContext.RequestContext)
}
