package main

import (
	"EM_test_task/internal/apis/implement"
	"EM_test_task/internal/handlers"
	"EM_test_task/pkg/storage"
	"EM_test_task/pkg/storage/db"
	"EM_test_task/pkg/storage/migrations"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"log"
	"log/slog"
	_ "net/http/pprof"
)

//		@title			People Information Enrichment Service
//		@description	Efficent Mobile test task January 2024
//
//	 @contact.name	Danil Vinogradov
//	 @contact.url		http://t.me/japsty
//	 @contact.email	danil-vinogradov-92@mail.ru
func main() {
	//Подгружаем переменные из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Создаем gin router
	gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.DebugMode)
	router := gin.Default()

	slog.Info("Server started")

	//Коннектимся к бд
	db, err := db.ConnectToDb()
	if err != nil {
		log.Fatal("Main ConnectToDb Error")
	}
	slog.Debug("бд работает")
	defer db.Close()

	//делаем миграцию
	provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)
	if err != nil {
		log.Fatal("Main failed to create NewProvider for migration")
	}

	provider.Up(context.Background())
	repo := storage.New(db)

	//встраиваем все сервисы
	agifyService := implement.AgifyService{}
	genderizeService := implement.GenderizeService{}
	nationalizeService := implement.NationalizeService{}
	ph := handlers.PersonHandler{
		PersonRepo:         repo,
		AgifyService:       &agifyService,
		GenderizeService:   &genderizeService,
		NationalizeService: &nationalizeService,
	}
	slog.Debug("сервисы подключены")

	//Сами эндпоинты
	router.GET("/people", ph.GetPersons)
	router.GET("/people/:id", ph.GetPerson)
	//router.GET("/people-filtered", ph.GetPersonsFiltered)
	router.DELETE("/people/:id", ph.DeletePersonByID)
	router.PUT("/people/:id", ph.UpdatePerson)
	router.POST("/people", ph.AddPerson)

	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Server dropped")
	}
}
