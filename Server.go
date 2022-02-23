package main

import (
	"github.com/envuso/go-http/Routing"
	"github.com/envuso/go-http/Routing/Route"
	"github.com/envuso/go-http/lib/storage"
)

func main() {
	Application.Boot(":8080", boot)
}

func boot(Router Routing.RouterContract) {
	// Router.Get("/yeet", new(CustomController), "DifferentResTypesPog")
	// Router.Get("/user/{user}", new(CustomController), "UsernameRoute")
	Router.Post("/user/{user}", new(CustomController), "UsernameRoute")

	Router.ControllerGroup(new(CustomController), func(stack *Route.RouteGroupStack) {
		// stack.Get("/yeet", "DifferentResTypesPog")
		stack.Get("/yeet", func() map[string]interface{} {
			return map[string]interface{}{
				"message": "hi",
			}
		})

		stack.Get("/storage", func(s storage.Storage) map[string]interface{} {
			f := s.AllDirectories("")
			print(f)

			return map[string]interface{}{
				"message": "hi",
			}
		})
	})

}
