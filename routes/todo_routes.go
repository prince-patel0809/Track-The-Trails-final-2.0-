package routes

import (
	"net/http"
	"track-the-trails/controllers"
	"track-the-trails/middlewares"
)

func TodoRoutes() {
	http.Handle(
		"/todos/create",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreateTodo)),
	)

	http.Handle(
		"/todos/get",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.GetAllTodos)),
	)

	http.Handle(
		"/todos/Update/{id}",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.UpdateTodo)),
	)

	http.Handle(
		"/todos/Delete/{id}",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.DeleteTodo)),
	)

	http.Handle(
		"/todos/Complate/{id}",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.CompleteTodo)),
	)

	http.Handle(
		"/todos/GetOne/{id}",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.GetTodoByID)),
	)
}
