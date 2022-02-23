package main

import (
	"gohttp/Routing"
	"gohttp/Routing/Route"
	"gohttp/lib/storage"
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
