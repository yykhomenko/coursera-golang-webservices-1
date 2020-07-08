package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

const DatasetFilename = "dataset.xml"
const Token = "test token"

type Dataset struct {
	Version string  `xml:"version,attr"`
	Entries []Entry `xml:"row"`
}

type Entry struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

func SearchServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req, err := parseRequest(r)
		if err != nil {
			fmt.Printf("read url params error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dataset, err := dataset(DatasetFilename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := SearchUsers(dataset, req)
		json.NewEncoder(w).Encode(resp.Users)
	}
}

func parseRequest(r *http.Request) (*SearchRequest, error) {
	q := r.URL.Query()

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		return nil, err
	}

	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil {
		return nil, err
	}

	query := q.Get("query")
	orderField := q.Get("order_field")

	orderBy, err := strconv.Atoi(q.Get("order_by"))
	if err != nil {
		return nil, err
	}

	return &SearchRequest{
		Limit:      limit,
		Offset:     offset,
		Query:      query,
		OrderField: orderField,
		OrderBy:    orderBy,
	}, nil
}

func dataset(filename string) (*Dataset, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("read dataset error: %v\n", err)
		return nil, err
	}

	dataset := new(Dataset)
	if err := xml.Unmarshal(content, &dataset); err != nil {
		fmt.Printf("dataset parse error %v\n", err)
		return nil, err
	}

	return dataset, nil
}

func SearchUsers(d *Dataset, r *SearchRequest) *SearchResponse {

	var users []User

	for _, e := range d.Entries {
		if strings.Contains(e.FirstName, r.Query) ||
			strings.Contains(e.LastName, r.Query) ||
			strings.Contains(e.About, r.Query) {

			users = append(users, User{
				Id:     e.ID,
				Name:   fmt.Sprintf("%s %s", e.FirstName, e.LastName),
				Age:    e.Age,
				About:  e.About,
				Gender: e.Gender,
			})
		}
	}

	return &SearchResponse{
		Users:    users,
		NextPage: false,
	}
}

func TestFindUsers(t *testing.T) {
	s := httptest.NewServer(SearchServer())

	type fields struct {
		AccessToken string
		URL         string
	}
	type args struct {
		req SearchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{name: "valid",
			fields: fields{AccessToken: Token, URL: s.URL},
			args: args{req: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "Owen",
				OrderField: "",
				OrderBy:    0,
			}},
			want: &SearchResponse{
				Users:    []User{},
				NextPage: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &SearchClient{
				AccessToken: tt.fields.AccessToken,
				URL:         tt.fields.URL,
			}
			got, err := srv.FindUsers(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
