package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validator "github.com/gobeam/custom-validator"

	"restapi/controllers"
)

func main() {
	r := gin.Default()

	// make extra validator
	validate := []validator.ExtraValidation{
		{Tag: "name", Message: "%s must be a string format."},
	}

	validator.MakeExtraValidation(validate)

	// handle url no route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Not Found"})
	})

	r.Use(validator.Errors())
	{

		// logistic
		logistic := r.Group("/logistics")
		{
			logistic.POST("/create", controllers.CreateLogistics)
			logistic.PUT("/update/:id", controllers.UpdateLogistics)
			logistic.DELETE("/delete/:id", controllers.DeleteLogistics)
			logistic.GET("/all-logistics", controllers.GetAllLogistics)
		}

		log_service := r.Group("/logistic-services")
		{
			log_service.POST("/create", controllers.CreateLogisticService)
			log_service.PUT("/update/:id", controllers.UpdateLogisticService)
			log_service.DELETE("/delete/:id", controllers.DeleteLogisticService)
			log_service.GET("/all-logistic-services", controllers.GetAllLogisticServices)
		}

		r.POST("/add-new-stock", controllers.AddStock)
	}

	r.Run()
}
