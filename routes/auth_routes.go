package routes

import (
	"net/http"
	"track-the-trails/controllers"
)

func AuthRoutes() {
	http.HandleFunc(
		"/register",
		controllers.RegisterUser,
	)
	http.HandleFunc(
		"/login",
		controllers.LoginUser,
	)
	http.HandleFunc(
		"/profile/upload",
		controllers.UploadProfileImage,
	)

	http.HandleFunc(
		"/getAll/users",
		controllers.GetAllUsers,
	)

	http.HandleFunc("/auth/google", controllers.GoogleLogin)

	http.HandleFunc("/auth/google/callback", controllers.GoogleCallback)
}
