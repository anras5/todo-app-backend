package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
	{"all-todos-incompleted", "/todos?completed=false", "GET", http.StatusOK},
	{"all-todos-completed-wrong", "/todos?completed=one", "GET", http.StatusBadRequest},
	{"one-todo", "/todos/1", "GET", http.StatusOK},
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
