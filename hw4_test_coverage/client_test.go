package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Row struct {
	ID    int    `xml:"id"`
	Login string `xml:"first_name"`
	Name  string `xml:"last_name"`
}

type Dataset struct {
	Version string `xml:"version,attr"`
	Rows    []Row  `xml:"row"`
}

func SearchServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("dataset.xml")
		if err != nil {
			log.Fatal(err)
		}

		v := new(Dataset)
		if err := xml.Unmarshal(content, &v); err != nil {
			fmt.Printf("dataset parse error %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(v.Rows)
	}
}

func TestFindUsers(t *testing.T) {

	s := httptest.NewServer(SearchServer())
	http.Get(s.URL)

	// type fields struct {
	// 	AccessToken string
	// 	URL         string
	// }
	// type args struct {
	// 	req SearchRequest
	// }
	// tests := []struct {
	// 	name    string
	// 	fields  fields
	// 	args    args
	// 	want    *SearchResponse
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		srv := &SearchClient{
	// 			AccessToken: tt.fields.AccessToken,
	// 			URL:         tt.fields.URL,
	// 		}
	// 		got, err := srv.FindUsers(tt.args.req)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("FindUsers() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("FindUsers() got = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
