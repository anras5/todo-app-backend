package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anras5/todo-app-backend/internal/config"
	"github.com/anras5/todo-app-backend/internal/models"
)

const (
	requestCount = 100000
	apiURL       = "http://localhost:8080"
)

func main() {

	// REST POST
	startTime := time.Now()
	ids := []int{}
	for i := 0; i < requestCount; i++ {
		todo := models.Todo{
			Name:        fmt.Sprintf("Todo number %d", i),
			Description: "This is a todo",
			Deadline:    time.Now().Add(time.Hour * 24),
			Completed:   false,
		}
		jsonResponse, err := sendCreateRequest(todo)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", i, err)
		}
		id := int(jsonResponse.Data.(float64))
		ids = append(ids, id)
	}
	duration := time.Since(startTime)
	fmt.Printf("POST %v.\n", duration)

	// REST GET ONE
	startTime = time.Now()
	for _, id := range ids {
		_, err := sendGetOneRequest(id)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("GET ONE %v.\n", duration)

	// REST UPDATE
	startTime = time.Now()
	for i, id := range ids {
		todo := models.Todo{
			ID:          id,
			Name:        fmt.Sprintf("Todo number %d", i),
			Description: "This is an updated todo",
			Deadline:    time.Now().Add(time.Hour * 24),
			Completed:   false,
		}
		_, err := sendUpdateRequest(todo)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("UPDATE %v.\n", duration)

	// REST GET ALL
	startTime = time.Now()
	_, err := sendGetAllRequest()
	if err != nil {
		fmt.Printf("Error sending request for all todos: %v\n", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("GET ALL %v.\n", duration)

	// REST DELETE
	startTime = time.Now()
	for _, id := range ids {
		_, err := sendDeleteRequest(id)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("DELETE %v.\n", duration)
}

func sendGetOneRequest(id int) (*models.Todo, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/todos/%d", apiURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var todo models.Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &todo, nil
}

func sendCreateRequest(todo models.Todo) (*config.JSONResponse, error) {
	requestBody, err := json.Marshal(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal todo: %w\n", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/todos", apiURL), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("received non-200 response: %d\n", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var jsonResponse config.JSONResponse
	err = json.Unmarshal(responseBody, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &jsonResponse, nil
}

func sendUpdateRequest(todo models.Todo) (*config.JSONResponse, error) {
	requestBody, err := json.Marshal(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal todo: %w\n", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/todos/%d", apiURL, todo.ID), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("received non-200 response: %d\n", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var jsonResponse config.JSONResponse
	err = json.Unmarshal(responseBody, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &jsonResponse, nil
}

func sendGetAllRequest() ([]models.Todo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/todos", apiURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var todos []models.Todo
	err = json.Unmarshal(body, &todos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return todos, nil
}

func sendDeleteRequest(id int) (*config.JSONResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/todos/%d", apiURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("received non-200 response: %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var jsonResponse config.JSONResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &jsonResponse, nil
}
