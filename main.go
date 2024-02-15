package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayedulkrm/go-mongo-curd/controllers"
)

func main() {

	router := gin.Default() // Address : Localhost:8000

	router.GET("/", func(ct *gin.Context) {
		// Set the content type to HTML
		ct.Header("Content-Type", "text/html")

		// Send the HTML response with an h1 tag
		html := `<h1>Server is working. To See Frontend <a href="http://localhost:3000"> Click Here </a></h1>`
		ct.String(http.StatusOK, html)
	})

	// create person
	router.POST("/create-user", controllers.CreatePerson)

	// get person by id
	router.GET("/user/:id", controllers.GetPerson)

	// Get All Person
	router.GET("/all-users", controllers.GetAllPerson)

	// Delete User

	router.DELETE("/user/:id", controllers.DeletePerson)

	// Update User

	router.PUT("/user/:id", controllers.UpdatePerson)

	router.Run(":8000")
}
