package main

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/btc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello BITCH!"))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text;charset-utf-8")

		fmt.Println(staticFiles.ReadDir("static"))

		templ, err := template.ParseFS(staticFiles, "static/*.html", "static/*.css") //.ParseFiles("./static/index.html")

		if err != nil {
			fmt.Println(err)
			return
		}
		// t := templ.Lookup("index.html")
		templ.ExecuteTemplate(w, "index.html", nil)
	})

	//http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":3000", r)
}
