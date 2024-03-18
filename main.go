package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

func main() {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(home))
	mux.Post("/", http.HandlerFunc(send))
	mux.Get("/confirmation", http.HandlerFunc(confirmation))

	log.Print("listening...")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/home.html", nil)
}

func send(w http.ResponseWriter, r *http.Request) {
	msg := &Message{
		Email:   r.PostFormValue("email"),
		Content: r.PostFormValue("content"),
	}

	if !msg.Validate() {
		render(w, "templates/home.html", msg)
		return
	}

	if err := msg.Deliver(); err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)

}

func confirmation(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/confirmation.html", nil)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
