package Middleware

type UsesMiddlewares struct {
	Middlewares *MiddlewareList
}

func (u *UsesMiddlewares) Created() {
	u.Middlewares = NewMiddlewareList()
}

func (u *UsesMiddlewares) Middleware(middlewares ...any) {
	mws := MapFromVariadic(middlewares...)

	for name, middleware := range mws {
		u.Middlewares.Set(name, middleware)
	}
}

func (u *UsesMiddlewares) HasMiddleware(middleware any) bool {
	mw, name := FromParamWithName(middleware)
	if mw == nil {
		return false
	}

	return u.Middlewares.Has(name)
}
