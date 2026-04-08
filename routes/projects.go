package routes

import (
	"net/http"
	"track-the-trails/controllers"
	"track-the-trails/middlewares"
)

func ProjectsRoutes() {

	http.Handle(
		"/project/create",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreateProject)),
	)

	http.Handle(
		"/project/delete/{id}",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.DeleteProject)),
	)

	http.Handle(
		"/project/update",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.UpdateProject)),
	)

	http.Handle(
		"/project/GetmyProject",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.GetMyProjects)),
	)

	http.Handle(
		"/project/GetmyProject/{id}",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.GetProjectDetails)),
	)

	http.Handle(
		"/projects/:projectId/members",
		middlewares.AuthMiddleware(http.HandlerFunc(controllers.AddMember)),
	)
}
