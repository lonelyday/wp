package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ServerData struct {
	Url      string `json:"url"`
	Interval int    `json:"interval"`
}

var Data []ServerData

type MyServ struct {
}

func (ms *MyServ) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var local_data ServerData
	rb, err := ioutil.ReadAll(r.Body)
	url_path := r.URL.Path

	params := strings.Split(url_path, "/")

	fmt.Println(params)

	if err != nil {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	switch r.Method {

	case http.MethodGet:
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(Data)

	// case http.MethodPatch:
	// case http.MethodDelete:

	case http.MethodPost:
		if len(rb) > 4096 {
			rw.WriteHeader(413)
			return
		} else if !isJSON(string(rb)) {
			rw.WriteHeader(400)
			return
		} else {
			err := json.Unmarshal(rb, &local_data)
			if err != nil {
				rw.WriteHeader(http.StatusCreated)
			} else {
				Data = append(Data, local_data)
				rw.WriteHeader(201)
				rw.Header().Set("Content-Type", "application/json")
				var ret_json bytes.Buffer
				fmt.Fprintf(&ret_json, `{"id":%d}`, len(Data))
				json.NewEncoder(rw).Encode(ret_json.String())
			}
		}
	}
}

type MyServList struct {
}

func (ms *MyServList) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rb, err := ioutil.ReadAll(r.Body)

	if err != nil {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	url_path := r.URL.Path
	params := strings.Split(url_path, "/")
	fmt.Println(params)

	switch r.Method {

	case http.MethodGet:

	case http.MethodPatch:
		if rec, err := strconv.Atoi(params[3]); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else {
			if rec > len(Data) {
				rw.WriteHeader(http.StatusNotFound)
				return
			} else {
				var local_data ServerData
				json.Unmarshal(rb, &local_data)
				Data[rec-1].Url = local_data.Url
				Data[rec-1].Interval = local_data.Interval
			}
		}
	case http.MethodDelete:
		if rec, err := strconv.Atoi(params[3]); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else {
			if rec > len(Data) {
				rw.WriteHeader(http.StatusNotFound)
				return
			} else {
				Data = append(Data[:rec-1], Data[rec:]...)
			}
		}
	}
}

func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

func main() {

	sm := http.NewServeMux()
	ms := MyServ{}
	ms1 := MyServList{}
	sm.Handle("/api/urls", &ms)
	sm.Handle("/api/urls/", &ms1)
	http.ListenAndServe(":8080", sm)
}
