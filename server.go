package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	Grille [][]string
	Joueur string
}

var grille [][]string
var joueur string = "🔴"

func initGrille() [][]string {
	grille := make([][]string, 6)
	for i := range grille {
		grille[i] = make([]string, 7)
	}
	return grille
}

func handlerAccueil(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Grille: grille,
		Joueur: joueur,
	}
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, data)
}

func handlerPlay(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	colStr := r.FormValue("col")
	c, err := strconv.Atoi(colStr)
	if err != nil || c < 0 || c > 6 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerWin(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Grille: grille,
		Joueur: joueur,
	}
	tmpl := template.Must(template.ParseFiles("win.html"))
	tmpl.Execute(w, data)
}

func handlerDraw(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Grille: grille,
		Joueur: joueur,
	}
	tmpl := template.Must(template.ParseFiles("full.html"))
	tmpl.Execute(w, data)
}

func handlerReset(w http.ResponseWriter, r *http.Request) {
	grille = initGrille()
	joueur = "🔴"
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	grille = initGrille()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlerAccueil)
	http.HandleFunc("/play", handlerPlay)
	http.HandleFunc("/win", handlerWin)
	http.HandleFunc("/full", handlerDraw)
	http.HandleFunc("/reset", handlerReset)
	http.ListenAndServe(":8080", nil)
}

func checkVictory(player string) bool {
	// Horizontal
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
	// Vertical
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
	// Diagonale ↘
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
	// Diagonale ↙
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
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			if grille[row][col] == "" {
				return false
			}
		}
	}
	return true
}
