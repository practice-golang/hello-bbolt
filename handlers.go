package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("It works!"))
}

func dataInitHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{"status": "fail", "message": ""}

	InitData()

	result["status"] = "success"
	ResponseJSON(w, http.StatusOK, result)
}

func addPersonHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{"status": "fail", "message": ""}
	var person Person

	if db == nil {
		result["message"] = "Enter password first"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		result["message"] = "Invalid request body"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	err = AddPerson(person)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		if strings.Contains(err.Error(), "required") {
			httpStatus = http.StatusBadRequest
		}

		result["message"] = "Failed to add person"
		ResponseJSON(w, httpStatus, result)
		return
	}

	result["status"] = "success"
	ResponseJSON(w, http.StatusOK, result)
}

func deletePersonHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{"status": "fail", "message": ""}

	if db == nil {
		result["message"] = "Enter password first"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	personID := r.URL.Query().Get("id")
	if personID == "" {
		result["message"] = "'id' is required"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	id, _ := strconv.Atoi(personID)
	err := DeletePerson(id)
	if err != nil {
		result["message"] = "Failed to delete person"
		ResponseJSON(w, http.StatusInternalServerError, result)
		return
	}

	result["status"] = "success"
	ResponseJSON(w, http.StatusOK, result)
}

func updatePersonHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{"status": "fail", "message": ""}

	if db == nil {
		result["message"] = "Enter password first"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	personID := r.URL.Query().Get("id")
	if personID == "" {
		result["message"] = "'id' is required"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	var updatedPerson Person
	err := json.NewDecoder(r.Body).Decode(&updatedPerson)
	if err != nil {
		result["message"] = "Invalid request body"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	// Here: Validate - Todo

	id, _ := strconv.Atoi(personID)
	err = UpdatePerson(id, updatedPerson)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		if strings.Contains(err.Error(), "Key not found") {
			httpStatus = http.StatusBadRequest
		}

		result["message"] = "Failed to update person"
		ResponseJSON(w, httpStatus, result)
		return
	}

	result["status"] = "success"
	ResponseJSON(w, http.StatusOK, result)
}

func getPersonListHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{"status": "fail", "message": ""}

	if db == nil {
		result["message"] = "Enter password first"
		ResponseJSON(w, http.StatusBadRequest, result)
		return
	}

	query := r.URL.Query()
	search := PersonSearch{
		Name:   query.Get("name"),
		Gender: query.Get("gender"),
		From:   query.Get("from"),
		To:     query.Get("to"),
	}

	limit := 300 // 목록 limit 기본값 = 300
	limitStr := query.Get("limit")
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	persons, err := GetPersonList(search, limit)
	if err != nil {
		result["message"] = "Failed to get persons: " + err.Error()
		ResponseJSON(w, http.StatusInternalServerError, result)
		return
	}

	result["status"] = "success"
	result["persons"] = persons
	ResponseJSON(w, http.StatusOK, result)
}

func setTextFileHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{"status": "fail", "message": ""}

	content := "Hello, Go!"
	file, err := os.Create("example.txt")
	if err != nil {
		result["message"] = err.Error()
		ResponseJSON(w, http.StatusInternalServerError, result)
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		result["message"] = err.Error()
		ResponseJSON(w, http.StatusInternalServerError, result)
		return
	}

	result["status"] = "success"
	ResponseJSON(w, http.StatusOK, result)
}

func getTextFileHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{"status": "fail", "message": ""}

	filePath := "example.txt"
	file, err := os.Open(filePath)
	if err != nil {
		result["message"] = err.Error()
		ResponseJSON(w, http.StatusInternalServerError, result)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		result["message"] = err.Error()
		ResponseJSON(w, http.StatusInternalServerError, result)
		return
	}

	result["status"] = "success"
	result["content"] = string(content)
	ResponseJSON(w, http.StatusOK, result)
}
