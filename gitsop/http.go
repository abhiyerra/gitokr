package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func getRepoTemplates(w http.ResponseWriter, r *http.Request) {

	c := NewRepoConfig(r.URL.Path[1:], nil)

	var s string
	s += "<html><head></head><body>"
	for title, task := range c.GetSOPs() {
		log.Println(title)

		_ = task
		s += "<form method='post'>"
		s += "<h1>" + title + "</h1>"
		s += "<input type='hidden' name='title' value='" + title + "'/>"
		s += "Name: <input type='text' name='name' /><br/>"
		if task.Inputs != nil {
			for name, input := range task.Inputs {
				s += name + ": <input name='" + name + "' type='text' value='" + input.Value + "'/><br/>"
			}
		}
		s += "<input type='submit' />"
		s += "</form>"
	}

	s += "</body>"

	fmt.Fprintf(w, s)
}

func postRepoTemplate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	c := NewRepoConfig(r.URL.Path[1:], nil)

	task := c.GetSOP(r.Form.Get("title"))

	for k := range task.Inputs {
		task.Inputs[k] = Input{Value: r.Form.Get(k)}
	}

	task.createSOP(c, r.Form.Get("name"))

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func RunHTTP() {
	r := mux.NewRouter()
	r.HandleFunc("/{owner}/{repo}", getRepoTemplates).Methods("GET")
	r.HandleFunc("/{owner}/{repo}", postRepoTemplate).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
