package handlers

import (
	"EM_test_task/internal/api"
	"EM_test_task/internal/entities"
	"EM_test_task/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type PersonInput struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type ApiInput struct {
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type PersonHandler struct {
	PersonRepo         storage.Repository
	AgifyService       api.AgifyGateway
	GenderizeService   api.GenderizeGateway
	NationalizeService api.NationalizeGateway
}

// AddPerson - метод, вызывающий CreatePerson и предоставляющий всю необходимую
// информацию для создания записи в бд
func (ph *PersonHandler) AddPerson(c *gin.Context) {
	var personInput PersonInput
	if err := c.BindJSON(&personInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("AddPerson BindJSON Error: ", err)
		return
	}

	ageCh := make(chan int)
	genderCh := make(chan string)
	nationCh := make(chan string)

	go func() {
		age, err := ph.AgifyService.GetAge(personInput.Name)
		if err != nil {
			slog.Error("AddPerson GetAge Error: ", err)
		}
		ageCh <- age
	}()

	go func() {
		gender, err := ph.GenderizeService.GetGender(personInput.Name)
		if err != nil {
			slog.Error("AddPerson GetGender Error: ", err)
		}
		genderCh <- gender
	}()

	go func() {
		nation, err := ph.NationalizeService.GetNationality(personInput.Surname)
		if err != nil {
			slog.Error("AddPerson GetNationality Error: ", err)
		}
		nationCh <- nation
	}()

	age := <-ageCh
	gender := <-genderCh
	nation := <-nationCh

	person, err := ph.PersonRepo.CreatePerson(c.Request.Context(), storage.PersonParams{
		Name:        personInput.Name,
		Surname:     personInput.Surname,
		Patronymic:  personInput.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nation,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
		slog.Error("AddPerson CreatePerson Error: ", err)
		return
	}

	slog.Info("Person added successfully")
	c.JSON(http.StatusCreated, person)
}

// GetPersons - метод, вызывающий GetPeople, возвращает список записей в бд,
// в соответствии с переданными настройками пагинации
func (ph *PersonHandler) GetPersons(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		slog.Error("Invalid 'page' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'page' parameter"})
		return
	}

	perPage, err := strconv.Atoi(c.Query("per_page"))
	if err != nil || perPage < 1 {
		slog.Error("Invalid 'perPage' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'per_page' parameter"})
		return
	}

	people, err := ph.PersonRepo.GetPeople(page, perPage)
	if err != nil {
		slog.Error("GetPersons GetPeople Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Get People": err.Error()})
		return
	}

	if len(people) == 0 {
		log.Info("No people found")
		c.JSON(http.StatusOK, gin.H{"message": "No people found"})
		return
	}

	slog.Info("People found successfully")
	c.JSON(http.StatusOK, people)
}

// GetPerson - метод, вызывающий метод GetPerson у бд, возвращает
// найденную по id запись
func (ph *PersonHandler) GetPerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("Invalid 'id' ")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	people, err := ph.PersonRepo.GetPerson(id)
	if err != nil {
		slog.Error("GetPerson GetPerson Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Get Person": err.Error()})
		return
	}

	slog.Info("Person found successfully")
	c.JSON(http.StatusOK, people)
}

func (ph *PersonHandler) GetPersonsFiltered(c *gin.Context) {
	var filter storage.Filter
	if err := c.BindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("GetPersonsFiltered BindQuery Error: ", err)
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		slog.Error("Invalid 'page' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'page' parameter"})
		return
	}

	perPage, err := strconv.Atoi(c.Query("per_page"))
	if err != nil || perPage < 1 {
		slog.Error("Invalid 'perPage' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'per_page' parameter"})
		return
	}

	people, err := ph.PersonRepo.GetPeopleFiltered(filter.SortBy, filter.From, filter.To, page, perPage)
	if err != nil {
		slog.Error("GetPersons GetPeople Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Get People": err.Error()})
		return
	}

	if len(people) == 0 {
		log.Info("No people found")
		c.JSON(http.StatusOK, gin.H{"message": "No people found"})
		return
	}

	slog.Info("People found successfully")
	c.JSON(http.StatusOK, people)
}

// UpdatePerson - метод обновляющий запись в бд в соответствии с новой информацией
func (ph *PersonHandler) UpdatePerson(c *gin.Context) {
	var personInput PersonInput
	if err := c.BindJSON(&personInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("UpdatePerson BindJSON Error: ", err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("Invalid 'id' ")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	ageCh := make(chan int)
	genderCh := make(chan string)
	nationCh := make(chan string)

	go func() {
		age, err := ph.AgifyService.GetAge(personInput.Name)
		if err != nil {
			slog.Error("AddPerson GetAge Error: ", err)
		}
		ageCh <- age
	}()

	go func() {
		gender, err := ph.GenderizeService.GetGender(personInput.Name)
		if err != nil {
			slog.Error("AddPerson GetGender Error: ", err)
		}
		genderCh <- gender
	}()

	go func() {
		nation, err := ph.NationalizeService.GetNationality(personInput.Surname)
		if err != nil {
			slog.Error("AddPerson GetNationality Error: ", err)
		}
		nationCh <- nation
	}()

	age := <-ageCh
	gender := <-genderCh
	nation := <-nationCh

	if age == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get age"})
		return
	}

	if gender == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gender"})
		return
	}

	if nation == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nationality"})
		return
	}

	currentPerson, err := ph.PersonRepo.GetPerson(id)
	if err != nil {
		slog.Error("UpdatePerson GetPerson Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"UpdatePerson": err.Error()})
		return
	}

	if personInput.Name == "" {
		personInput.Name = currentPerson.Name
	}
	if personInput.Surname == "" {
		personInput.Surname = currentPerson.Surname
	}
	if personInput.Patronymic == "" {
		personInput.Patronymic = currentPerson.Patronymic
	}

	err = ph.PersonRepo.UpdatePerson(id, entities.Person{
		Name:        personInput.Name,
		Surname:     personInput.Surname,
		Patronymic:  personInput.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nation,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
		slog.Error("UpdatePerson UpdatePerson Error: ", err)
		return
	}

	slog.Info("Person updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Person updated successfully"})
}

// DeletePersonByID - метод вызывающий DeletePerson у бд, удаляет запись в соответствии с id
func (ph *PersonHandler) DeletePersonByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	err = ph.PersonRepo.DeletePerson(id)
	if err != nil {
		slog.Error("DeletePersonByID Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Person deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Person was deleted"})
}
