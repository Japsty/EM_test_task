package controller

import (
	"EM_test_task/internal/apis"
	"EM_test_task/internal/entities"
	"EM_test_task/pkg/storage"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

type PersonInput struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	//Age         int    `json:"age,omitempty"`
	//Gender      string `json:"gender,omitempty"`
	//Nationality string `json:"nationality,omitempty"`
}

type ApiInput struct {
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type PersonHandler struct {
	PersonRepo         storage.Repository
	AgifyService       apis.AgifyGateway
	GenderizeService   apis.GenderizeGateway
	NationalizeService apis.NationalizeGateway
}

func (ph *PersonHandler) AddPerson(c *gin.Context) {
	var personInput PersonInput
	if err := c.BindJSON(&personInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("AddPerson BindJSON Error: ", err)
		return
	}

	age, err := ph.AgifyService.GetAge(personInput.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get age"})
		slog.Error("AddPerson GetAge Error: ", err)
		return
	}
	gender, err := ph.GenderizeService.GetGender(personInput.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gender"})
		slog.Error("AddPerson GetGender Error: ", err)
		return
	}
	nation, err := ph.NationalizeService.GetNationality(personInput.Surname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nationality"})
		slog.Error("AddPerson GetNationality Error: ", err)
		return
	}

	person, err := ph.PersonRepo.CreatePerson(c.Request.Context(), storage.PersonParams{
		Name:        personInput.Name,
		Surname:     personInput.Surname,
		Patronymic:  personInput.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nation,
	})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		slog.Error("AddPerson CreatePerson Error: ", err)
		return
	}

	c.JSON(http.StatusCreated, person)
}

func (ph *PersonHandler) GetPersons(c *gin.Context) {
	people, err := ph.PersonRepo.GetPeople()
	if err != nil {
		slog.Error("GetPersons Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(people) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No people found"})
		return
	}

	c.JSON(http.StatusOK, people)
}

func (ph *PersonHandler) UpdatePerson(c *gin.Context) {
	var personInput PersonInput
	if err := c.BindJSON(&personInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("UpdatePerson BindJSON Error: ", err)
		return
	}

	log.Println(c.Param("id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	err = ph.PersonRepo.UpdatePerson(id, entities.Person{
		Name:       personInput.Name,
		Surname:    personInput.Surname,
		Patronymic: personInput.Patronymic,
		//Age:         personInput.Age,
		//Gender:      personInput.Gender,
		//Nationality: personInput.Nationality,
	})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
		slog.Error("UpdatePerson UpdatePerson Error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person updated successfully"})
}

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

	c.JSON(http.StatusOK, gin.H{"message": "Person was deleted"})
}
