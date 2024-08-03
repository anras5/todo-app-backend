package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anras5/todo-app-backend/internal/models"
)

const (
	requestCount = 100000
	apiURL       = "http://localhost:8080/graphql"
)

type graphQLResponse struct {
	Data map[string]interface{} `json:"data"`
}

func main() {
	// GraphQL requests

	ids := []int{}

	// create todos
	startTime := time.Now()
	for i := 0; i < requestCount; i++ {
		todo := models.Todo{
			Name:        fmt.Sprintf("Todo number %d", i),
			Description: "This is a todo from graphql",
			Deadline:    time.Now().Add(time.Hour * 24),
			Completed:   false,
		}
		graphqlResponse, err := sendCreateRequest(todo)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", i, err)
		} else {
			id := int(graphqlResponse.Data["createTodo"].(map[string]interface{})["id"].(float64))
			ids = append(ids, id)
		}
	}
	duration := time.Since(startTime)
	fmt.Printf("createTodo %v.\n", duration)

	// get todos by one
	startTime = time.Now()
	for _, id := range ids {
		_, err := sendGetOneRequest(id)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}

	}
	duration = time.Since(startTime)
	fmt.Printf("getTodo %v.\n", duration)

	// update todos
	startTime = time.Now()
	for i, id := range ids {
		todo := models.Todo{
			ID:          id,
			Name:        fmt.Sprintf("Todo number %d", i),
			Description: "This is an updated todo from graphql",
			Deadline:    time.Now().Add(time.Hour * 24),
			Completed:   true,
		}
		_, err := sendUpdateRequest(todo)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("updateTodo %v.\n", duration)

	// get all todos
	startTime = time.Now()
	_, err := sendGetAllRequest()
	if err != nil {
		fmt.Printf("Error sending request for all todos: %v\n", err)
	}
	duration = time.Since(startTime)
	fmt.Printf("getTodos %v.\n", duration)

	// delete todos
	startTime = time.Now()
	for _, id := range ids {
		_, err := sendDeleteRequest(id)
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("deleteTodo %v.\n", duration)
}

func sendCreateRequest(todo models.Todo) (*graphQLResponse, error) {
	deadlineStr := todo.Deadline.Format(time.RFC3339)
	createTodoMutation := fmt.Sprintf(`mutation {
			createTodo(name: "%s", description: "%s", deadline: "%s", completed: %t) {
				id
			}
		}`, todo.Name, todo.Description, deadlineStr, todo.Completed)
	requestBody := fmt.Sprintf(`{"query": %q}`, createTodoMutation)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResp graphQLResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, err
	}
	return &jsonResp, nil
}

func sendGetOneRequest(id int) (*models.Todo, error) {
	getTodoQuery := fmt.Sprintf(`query {
			getTodo(id: %d) {
				id
				name
				description
				deadline
				completed
			}
		}`, id)
	requestBody := fmt.Sprintf(`{"query": %q}`, getTodoQuery)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResp graphQLResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, err
	}

	parsedTime, _ := time.Parse(time.RFC3339, jsonResp.Data["getTodo"].(map[string]interface{})["deadline"].(string))
	todo := &models.Todo{
		ID:          int(jsonResp.Data["getTodo"].(map[string]interface{})["id"].(float64)),
		Name:        jsonResp.Data["getTodo"].(map[string]interface{})["name"].(string),
		Description: jsonResp.Data["getTodo"].(map[string]interface{})["description"].(string),
		Deadline:    parsedTime,
		Completed:   jsonResp.Data["getTodo"].(map[string]interface{})["completed"].(bool),
	}
	return todo, nil
}

func sendUpdateRequest(todo models.Todo) (*graphQLResponse, error) {
	deadlineStr := todo.Deadline.Format(time.RFC3339)
	updateTodoMutation := fmt.Sprintf(`mutation {
			updateTodo(id: %d, name: "%s", description: "%s", deadline: "%s", completed: %t) {
				id
			}
		}`, todo.ID, todo.Name, todo.Description, deadlineStr, todo.Completed)
	requestBody := fmt.Sprintf(`{"query": %q}`, updateTodoMutation)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResp graphQLResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, err
	}
	return &jsonResp, nil
}

func sendGetAllRequest() ([]models.Todo, error) {
	getAllTodosQuery := `query {
			getTodos {
				id,
				name,
				description,
				deadline,
				completed
			}
		}`
	requestBody := fmt.Sprintf(`{"query": %q}`, getAllTodosQuery)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResp graphQLResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, err
	}
	var todos []models.Todo
	for _, todo := range jsonResp.Data["getTodos"].([]interface{}) {
		parsedTime, _ := time.Parse(time.RFC3339, todo.(map[string]interface{})["deadline"].(string))
		todos = append(todos, models.Todo{
			ID:          int(todo.(map[string]interface{})["id"].(float64)),
			Name:        todo.(map[string]interface{})["name"].(string),
			Description: todo.(map[string]interface{})["description"].(string),
			Deadline:    parsedTime,
			Completed:   todo.(map[string]interface{})["completed"].(bool),
		})
	}
	return todos, nil
}

func sendDeleteRequest(id int) (*graphQLResponse, error) {
	deleteTodoMutation := fmt.Sprintf(`mutation {
			deleteTodo(id: %d) {
				id
			}
		}`, id)
	requestBody := fmt.Sprintf(`{"query": %q}`, deleteTodoMutation)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonResp graphQLResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, err
	}
	return &jsonResp, nil
}
