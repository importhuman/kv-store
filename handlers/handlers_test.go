package kv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// Test for "Set" API (POST method)
func TestCheckSetHandler(t *testing.T) {
	// Create request to pass to handler (serves as "r *http.Request" in the test)
	data := strings.NewReader(`{"abc-1":"one"}`)
	req, err := http.NewRequest("POST", "/set", data)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response.
	// (This will serve as the "w http.ResponseWriter" in the test.)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Set)

	// Handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"abc-1":"one"}`
	// Length of actual response was more than expected
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// Test for "Get" API
func TestCheckGetHandler(t *testing.T) {
	// table-driven test
	tt := []struct {
		key   string
		value interface{}
		pass  bool
	}{
		{"abc-1", 1, true},
		{"abc-2", 2, true},
		{"xyz-1", "three", true},
		{"xyz-2", 4, true},
		{"xyz-4", 5, false},
	}

	// send set request to initialize kv store
	data := strings.NewReader(`{"abc-1":1,"abc-2":2,"xyz-1":"three","xyz-2":4}`)
	req, err := http.NewRequest("POST", "/set", data)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Set)
	handler.ServeHTTP(rr, req)

	// set up and execute test cases
	for _, tc := range tt {
		path := fmt.Sprintf("/get/%s", tc.key)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// Need to create a router that we can pass the request through so that the vars will be added to the context
		r := mux.NewRouter()
		r.HandleFunc("/get/{key}", Get)
		r.ServeHTTP(rr, req)

		// If test should pass, but it fails
		if rr.Code != http.StatusOK && tc.pass {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}

		// If test should fail, but it passes
		if rr.Code == http.StatusOK && !tc.pass {
			t.Errorf("handler should have failed on routeVariable %s: got %v want %v",
				tc.key, rr.Code, http.StatusBadRequest)
		}

		// Compare expected and received values
		expected := fmt.Sprintf("%v", tc.value)
		received := fmt.Sprintf("%v", strings.Trim(strings.TrimSpace(rr.Body.String()), "\""))
		if received != expected && tc.pass {
			t.Errorf("handler returned unexpected body: got %v want %v", received, expected)
		}
	}
}

// Test for "GetAll" API
func TestCheckGetAllHandler(t *testing.T) {
	// send set request to initialize kv store
	data := strings.NewReader(`{"abc-1":1,"abc-2":2,"xyz-1":"three","xyz-2":4}`)
	req, err := http.NewRequest("POST", "/set", data)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	setHandler := http.HandlerFunc(Set)
	setHandler.ServeHTTP(rr, req)

	// send get request at "/" endpoint for all keys and values
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	getHandler := http.HandlerFunc(GetAll)
	getHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// compare expected and received body
	expected := `{"abc-1":1,"abc-2":2,"xyz-1":"three","xyz-2":4}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// Test for "Search" API
func TestCheckSearchHandler(t *testing.T) {
	// test table for search
	tt := []struct {
		category string
		search   string
		expect   map[string]bool
	}{
		{"prefix", "abc", map[string]bool{"abc-1": true, "abc-2": true}},
		{"prefix", "xyz", map[string]bool{"xyz-1": true, "xyz-2": true}},
		{"suffix", "1", map[string]bool{"abc-1": true, "xyz-1": true}},
		{"suffix", "-2", map[string]bool{"abc-2": true, "xyz-2": true}},
	}

	// send set requests to initialize kv store
	data := strings.NewReader(`{"abc-1":1,"abc-2":2,"xyz-1":"three","xyz-2":4}`)
	req, err := http.NewRequest("POST", "/set", data)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Set)
	handler.ServeHTTP(rr, req)

	// set up and execute search queries
	for _, tc := range tt {
		req, err := http.NewRequest("GET", "/search", nil)
		if err != nil {
			t.Fatal(err)
		}
		q := req.URL.Query()
		q.Add(tc.category, tc.search)
		req.URL.RawQuery = q.Encode()
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Search)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// parse body into an array of the keys without quotes
		arr := strings.Split(strings.ReplaceAll(strings.Trim(strings.TrimSpace(rr.Body.String()), "[]"), "\"", ""), ",")

		if len(arr) != len(tc.expect) {
			t.Errorf("handler returned unexpected number of keys: got %v want %v", len(arr), len(tc.expect))
		}

		// check if key was expected
		for _, v := range arr {
			if tc.expect[v] != true {
				t.Errorf("handler returned unexpected key: got %v want %v", v, tc.expect[v])
			}
		}
	}
}
