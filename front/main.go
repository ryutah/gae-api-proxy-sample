package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	backendHost   = os.Getenv("BACKEND_HOST")
	backendSchema = os.Getenv("BACKEND_SCHEMA")
)

func init() {
	if backendHost == "" || backendSchema == "" {
		panic("env var BACKEND_HOST and BACKEND_SCHEMA is required.")
	}
}

func requestURL(path string, query map[string]string) *url.URL {
	var rawQuery string
	if query != nil {
		queryPair := make([]string, 0, len(query))
		for k, v := range query {
			queryPair = append(queryPair, fmt.Sprintf("%s=%s", k, v))
		}
		rawQuery = strings.Join(queryPair, "&")
	}
	return &url.URL{
		Host:     backendHost,
		Scheme:   backendSchema,
		Path:     path,
		RawQuery: rawQuery,
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(requestURL("ping", nil).String())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		w.Write(body)
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		echo := r.FormValue("val")
		if echo == "" {
			echo = "HelloWorld!!"
		}
		resp, err := http.Get(requestURL("/echo", map[string]string{
			"val": echo,
		}).String())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		w.Write(body)
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		reqPayload, _ := json.Marshal(map[string]string{
			"request": "send_request",
		})
		body := bytes.NewBuffer(reqPayload)
		resp, err := http.Post(requestURL("/json", nil).String(), "application/json", body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
