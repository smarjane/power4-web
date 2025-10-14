package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageData struct {
	Grille [][]string
	Joueur string
	Winner string
	Looser string
	Nom    string
}

var grille [][]string
var joueur string = "🔴"
var nomUtilisateur string = ""
var templates *template.Template
var rows, cols int = 6, 7 // ✅ dimensions par défaut, modifiables selon difficulté

func initGrille() [][]string {
	grille := make([][]string, rows)
	for i := range grille {
		grille[i] = make([]string, cols)
	}
	return grille
}

func handlerStart(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		player1 := r.FormValue("player1")
		player2 := r.FormValue("player2")
		difficulty := r.FormValue("difficulty")

		if player1 == "" || player2 == "" || difficulty == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		nomUtilisateur = player1 + " vs " + player2

		switch difficulty {
		case "facile":
			rows, cols = 6, 7
		case "normal":
			rows, cols = 6, 9
		case "difficile":
			rows, cols = 7, 8
		default:
			rows, cols = 6, 7
		}

		grille = initGrille()
		joueur = "🔴"
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	} else {
		tmpl, err := template.ParseFiles("html/start.html")
		if err != nil {
			http.Error(w, "Erreur", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func handlerGame(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Grille: grille,
		Joueur: joueur,
		Nom:    nomUtilisateur,
	}
	tmpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
	}
	tmpl.Execute(w, data)
}

func handlerPlay(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	colStr := r.FormValue("col")
	c, err := strconv.Atoi(colStr)
	if err != nil || c < 0 || c > 6 {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	for i := 5; i >= 0; i-- {
		if grille[i][c] == "" {
			grille[i][c] = joueur
			if checkVictory(joueur) {
				http.Redirect(w, r, "/win", http.StatusSeeOther)
				return
			} else if isDraw() {
				http.Redirect(w, r, "/full", http.StatusSeeOther)
				return
			}
			if joueur == "🔴" {
				joueur = "🟡"
			} else {
				joueur = "🔴"
			}
			break
		}
	}
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func handlerWin(w http.ResponseWriter, r *http.Request) {
	looser := "🔴"
	if joueur == "🔴" {
		looser = "🟡"
	}
	data := PageData{
		Grille: grille,
		Joueur: joueur,
		Winner: joueur,
		Looser: looser,
		Nom:    nomUtilisateur,
	}
	tmpl, err := template.ParseFiles("html/win.html")
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
	}
	tmpl.Execute(w, data)
}

func handlerDraw(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Grille: grille,
		Joueur: joueur,
		Nom:    nomUtilisateur,
	}
	tmpl, err := template.ParseFiles("html/full.html")
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
	}
	tmpl.Execute(w, data)
}

func handlerReset(w http.ResponseWriter, r *http.Request) {
	grille = initGrille()
	joueur = "🔴"
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func checkVictory(player string) bool {
	for row := 0; row < 6; row++ {
		for col := 0; col <= 3; col++ {
			if grille[row][col] == player &&
				grille[row][col+1] == player &&
				grille[row][col+2] == player &&
				grille[row][col+3] == player {
				return true
			}
		}
	}
	for col := 0; col < 7; col++ {
		for row := 0; row <= 2; row++ {
			if grille[row][col] == player &&
				grille[row+1][col] == player &&
				grille[row+2][col] == player &&
				grille[row+3][col] == player {
				return true
			}
		}
	}
	for row := 0; row <= 2; row++ {
		for col := 0; col <= 3; col++ {
			if grille[row][col] == player &&
				grille[row+1][col+1] == player &&
				grille[row+2][col+2] == player &&
				grille[row+3][col+3] == player {
				return true
			}
		}
	}
	for row := 3; row < 6; row++ {
		for col := 0; col <= 3; col++ {
			if grille[row][col] == player &&
				grille[row-1][col+1] == player &&
				grille[row-2][col+2] == player &&
				grille[row-3][col+3] == player {
				return true
			}
		}
	}
	return false
}

func isDraw() bool {
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if grille[row][col] == "" {
				return false
			}
		}
	}
	return true
}

func main() {
	grille = initGrille()
	templates = template.Must(template.ParseFiles(
		"html/start.html",
		"html/index.html",
		"html/win.html",
		"html/full.html",
	))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlerStart)
	http.HandleFunc("/start", handlerStart)
	http.HandleFunc("/game", handlerGame)
	http.HandleFunc("/play", handlerPlay)
	http.HandleFunc("/win", handlerWin)
	http.HandleFunc("/full", handlerDraw)
	http.HandleFunc("/reset", handlerReset)
	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
