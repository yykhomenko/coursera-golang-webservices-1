package main

import (
	"encoding/json"
	"encoding/xml"
	_ "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type row struct {
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
}

type dataset struct {
	Version string `xml:"version"`
	Row     []row  `xml:"row"`
}

const pageSize = 25

func SearchServerSuccess(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("dataset.xml")
	checkError(err)

	dataset := &dataset{}
	xml.Unmarshal(f, &dataset)

	var users []User
	for _, user := range dataset.Row {
		users = append(users, User{
			Id:     user.Id,
			Name:   user.FirstName,
			Age:    user.Age,
			About:  user.About,
			Gender: user.Gender,
		})
	}

	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	var startRow int
	if offset > 0 {
		startRow = offset * pageSize
	}

	endRow := startRow + limit
	users = users[startRow:endRow]

	response, err := json.Marshal(users)
	checkError(err)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func SearchServerLimitFail(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("dataset.xml")
	checkError(err)

	dataset := &dataset{}
	xml.Unmarshal(f, &dataset)

	var users []User
	for _, user := range dataset.Row {
		users = append(users, User{
			Id:     user.Id,
			Name:   user.FirstName,
			Age:    user.Age,
			About:  user.About,
			Gender: user.Gender,
		})
	}

	response, err := json.Marshal(users)
	checkError(err)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func TestErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServerSuccess))
	defer server.Close()
	client := &SearchClient{URL: server.URL}
	request := SearchRequest{Limit: 5, Offset: 0}

	_, err := client.FindUsers(request)
	if err != nil {
		t.Error("Doesn't work success request")
	}

	request.Limit = -1

	_, err = client.FindUsers(request)
	if err.Error() != "limit must be > 0" {
		t.Error("limit must be > 0")
	}

	request.Limit = 1
	request.Offset = -1
	_, err = client.FindUsers(request)
	if err.Error() != "offset must be > 0" {
		t.Error("offset must be > 0")
	}
}

func TestOverLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServerSuccess))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	response, _ := client.FindUsers(SearchRequest{Limit: 26})

	if pageSize != len(response.Users) {
		t.Error("Over limit")
	}
}

func TestLimitFailed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServerLimitFail))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	limit := 7
	response, _ := client.FindUsers(SearchRequest{Limit: limit})
	if limit == len(response.Users) {
		t.Error("Limit not true")
	}
}

func TestBadJson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `"err": "bad json"}`)
	}))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	_, err := client.FindUsers(SearchRequest{})
	if err.Error() != `cant unpack result json: invalid character ':' after top-level value` {
		t.Error("Bad json test")
	}
}

func TestTimeoutError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Error("Timeout check error")
	}
}

func TestStatusUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	_, err := client.FindUsers(SearchRequest{})
	if err.Error() != "Bad AccessToken" {
		t.Error("Bad AccessToken")
	}
}

func TestBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	searchClient := &SearchClient{URL: server.URL}

	_, err := searchClient.FindUsers(SearchRequest{})
	if err.Error() != "cant unpack error json: unexpected end of JSON input" {
		t.Error("TestBadRequest is not done")
	}
}

func TestBadField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
		w.Write(jsonResponse)
	}))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	_, err := client.FindUsers(SearchRequest{})
	if err.Error() != "OrderFeld  invalid" {
		t.Error("ErrorBadOrderField is not done")
	}
}

func TestBadRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(SearchErrorResponse{Error: "Unknown error"})
		w.Write(jsonResponse)
	}))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Error("TestBadRequestError is not done")
	}
}

func TestStatusInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	client := &SearchClient{URL: server.URL}

	_, err := client.FindUsers(SearchRequest{})
	if err.Error() != "SearchServer fatal error" {
		t.Error("SearchServer error")
	}
}

func TestUnknownError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { /* NOP */ }))
	defer server.Close()
	client := &SearchClient{URL: "bad_link"}

	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Error("Test unknown error")
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
