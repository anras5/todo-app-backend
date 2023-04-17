package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var theTests = []struct {
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
}

func TestRepository_Home(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.Home)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected %d but got %d", http.StatusOK, rr.Code)
	}

}
