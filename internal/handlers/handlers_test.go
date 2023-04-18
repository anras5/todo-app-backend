package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/anras5/todo-app-backend/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var theGetTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"all-todos", "/todos", "GET", http.StatusOK},
	{"all-todos-completed", "/todos?completed=true", "GET", http.StatusOK},
	{"all-todos-incompleted", "/todos?completed=false", "GET", http.StatusBadRequest},
	{"all-todos-completed-wrong", "/todos?completed=one", "GET", http.StatusBadRequest},
	{"one-todo-1", "/todos/1", "GET", http.StatusOK},
	{"one-todo-2", "/todos/2", "GET", http.StatusBadRequest},
	{"one-todo-invalid-parameter", "/todos/one", "GET", http.StatusBadRequest},
}

// TestHandlers tests all routes that don't require extra tests (GET handlers)
func TestHandlers(t *testing.T) {

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theGetTests {
		response, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
		}
	}

}

var theInsertTests = []struct {
	name               string
	todo               interface{}
	method             string
	expectedStatusCode int
}{
	{
		name: "valid todo",
		todo: models.Todo{
			Name:        "valid todo",
			Description: "todo made for the purposes of testing",
			Deadline:    time.Time{},
			Completed:   false,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		method:             "POST",
		expectedStatusCode: http.StatusAccepted,
	},
	{
		name: "invalid todo",
		todo: struct {
			TestField string
		}{
			TestField: "onetwothree",
		},
		method:             "POST",
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name: "empty name todo",
		todo: models.Todo{
			Name:        "",
			Description: "todo made for the purposes of testing",
			Deadline:    time.Time{},
			Completed:   false,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		method:             "POST",
		expectedStatusCode: http.StatusBadRequest,
	},
}

func TestRepository_InsertTodo(t *testing.T) {
	for _, e := range theInsertTests {
		var req *http.Request
		jsonTestTodo, _ := json.Marshal(e.todo)
		req, _ = http.NewRequest(e.method, "/todos", bytes.NewBuffer(jsonTestTodo))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.InsertTodo)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

	}

}
