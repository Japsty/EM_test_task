package main

import (
	"EM_test_task/internal/apis/implement"
	"EM_test_task/internal/controller"
	"EM_test_task/pkg/storage"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log/slog"
	_ "net/http/pprof"
)

//	@title			People Information Enrichment Service
//	@description	Efficent Mobile test task January 2024

// @contact.name	Danil Vinogradov
// @contact.url		http://t.me/japsty
// @contact.email	danil-vinogradov-92@mail.ru
func main() {
	router := gin.Default()

	db, err := sql.Open("postgres", "dbname=GinUrlShortener sslmode=disable")
	slog.Debug("Db connection opened")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}

	repo := storage.New(db)
	agifyService := implement.AgifyService{}
	genderizeService := implement.GenderizeService{}
	nationalizeService := implement.NationalizeService{}
	ph := controller.PersonHandler{
		PersonRepo:         repo,
		AgifyService:       &agifyService,
		GenderizeService:   &genderizeService,
		NationalizeService: &nationalizeService,
	}

	router.GET("/people", ph.GetPersons)
	router.DELETE("/people/:id", ph.DeletePersonByID)
	router.PUT("/people/:id", ph.UpdatePerson)
	router.POST("/people", ph.AddPerson)

	router.Run("localhost:8080")
	slog.Debug("Server is ready to use")
}
