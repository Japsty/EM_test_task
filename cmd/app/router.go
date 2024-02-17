package main

import (
	"EM_test_task/internal/api/implement"
	"EM_test_task/internal/handlers"
	"EM_test_task/internal/middleware"
	"EM_test_task/pkg/logger"
	"EM_test_task/pkg/storage"
	"github.com/gin-gonic/gin"
)

func SetupRouter(repo storage.Repository) *gin.Engine {
	router := gin.Default()
	router.Use(logger.LogMiddleware())

	agifyService := implement.AgifyService{}
	genderizeService := implement.GenderizeService{}
	nationalizeService := implement.NationalizeService{}
	ph := handlers.PersonHandler{
		PersonRepo:         repo,
		AgifyService:       &agifyService,
		GenderizeService:   &genderizeService,
		NationalizeService: &nationalizeService,
	}

	uh := handlers.UserHandler{UserRepo: repo}
	router.POST("/registration", uh.Registration)
	router.POST("/login", uh.Login)

	router.Use(middleware.AuthMiddleware())

	router.GET("/people", ph.GetPersons)
	router.GET("/people/:id", ph.GetPerson)
	router.GET("/people-filtered", ph.GetPersonsFiltered)
	router.DELETE("/people/:id", ph.DeletePersonByID)
	router.PUT("/people/:id", ph.UpdatePerson)
	router.POST("/people", ph.AddPerson)

	return router
}
