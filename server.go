package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Grille [][]string
}

func initGrille() [][]string {
	grille := make([][]string, 6)
	for i := range grille {
		grille[i] = make([]string, 7)
	}
	return grille
}

func handlerAccueil(w http.ResponseWriter, r *http.Request) {
	data := PageData{Grille: initGrille()}
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, data)
}

func main() {
	initGrille()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlerAccueil)
	http.ListenAndServe(":8080", nil)
}
