package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/anras5/todo-app-backend/internal/models"
	"github.com/go-chi/chi/v5"
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

var theUpdateTests = []struct {
	name               string
	todo               interface{}
	id                 int
	method             string
	expectedStatusCode int
}{
	{
		name: "valid update",
		todo: models.Todo{
			ID:          1,
			Name:        "",
			Description: "",
			Deadline:    time.Time{},
			Completed:   false,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		id:                 1,
		method:             "PUT",
		expectedStatusCode: http.StatusAccepted,
	},
	{
		name: "invalid update",
		todo: struct {
			TestField string
		}{
			TestField: "onetwothree",
		},
		id:                 1,
		method:             "PUT",
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name: "wrong id update",
		todo: models.Todo{
			ID:          2,
			Name:        "",
			Description: "",
			Deadline:    time.Time{},
			Completed:   false,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		id:                 2,
		method:             "PUT",
		expectedStatusCode: http.StatusBadRequest,
	},
}

func TestRepository_UpdateTodo(t *testing.T) {
	for _, e := range theUpdateTests {
		var req *http.Request
		jsonTestTodo, _ := json.Marshal(e.todo)
		req, _ = http.NewRequest(e.method, fmt.Sprintf("/todos/%d", e.id), bytes.NewBuffer(jsonTestTodo))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.UpdateTodo)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

var theUpdateCompletedTests = []struct {
	name               string
	id                 string
	completed          string
	expectedStatusCode int
}{
	{
		name:               "valid-update-completed",
		id:                 "1",
		completed:          "complete",
		expectedStatusCode: http.StatusAccepted,
	},
	{
		name:               "invalid-id-update-completed",
		id:                 "one",
		completed:          "complete",
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:               "invalid-complete-variable-update-completed",
		id:                 "1",
		completed:          "some string",
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:               "db-error-complete-update-completed",
		id:                 "2",
		completed:          "complete",
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		name:               "db-error-incomplete-update-completed",
		id:                 "2",
		completed:          "incomplete",
		expectedStatusCode: http.StatusBadRequest,
	},
}

func TestRepository_UpdateTodoCompleted(t *testing.T) {
	for _, e := range theUpdateCompletedTests {

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/todos/%s/%s", e.id, e.completed), nil)

		// handle context and path variables in chi router
		ctx := req.Context()
		pathVariables := make(map[string]string)
		pathVariables["id"] = e.id
		pathVariables["complete"] = e.completed
		ctx = addChiContext(ctx, pathVariables)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.UpdateTodoCompleted)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

func addChiContext(parentCtx context.Context, pathVariables map[string]string) context.Context {
	chiCtx := chi.NewRouteContext()
	for k, v := range pathVariables {
		chiCtx.URLParams.Add(k, v)
	}
	return context.WithValue(parentCtx, chi.RouteCtxKey, chiCtx)
}
