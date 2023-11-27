package main

import (
	"bufio"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Jeu struct {
	Mot        string
	MotActuel  string
	Erreurs    int
	MaxErreurs int
	Difference int
	Difficulte string
	LettresNP  []string
	Sortie     string
	Pseudo     string
}

var JeuActuel *Jeu
var tpl *template.Template

func MotRandom() (string, error) {
	file, err := os.Open("noms_monstres.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var Mots []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Mots = append(Mots, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	return Mots[rand.Intn(len(Mots))], nil
}

func MaxErreurs(difficulty string) int {
	switch difficulty {
	case "Facile":
		return 8
	case "Intermediaire":
		return 6
	case "Difficile":
		return 4
	default:
		return 6
	}
}

func NouvellePartie(pseudo, difficulty string) *Jeu {
	MotToGuess, err := MotRandom()
	if err != nil {
		panic(err)
	}

	MaxErreurs := MaxErreurs(difficulty)

	return &Jeu{
		Mot:        MotToGuess,
		MotActuel:  strings.Repeat("_", len(MotToGuess)),
		Erreurs:    0,
		MaxErreurs: MaxErreurs,
		Difficulte: difficulty,
		LettresNP:  []string{},
		Sortie:     "",
		Pseudo:     pseudo,
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	if JeuActuel == nil {
		JeuActuel = NouvellePartie(JeuActuel.Pseudo, JeuActuel.Difficulte)

	}

	jeu := JeuActuel

	if r.Method == http.MethodPost {
		guess := strings.ToLower(r.FormValue("guess"))

		if len(guess) == 1 {
			if strings.Contains(strings.ToLower(jeu.Mot), strings.ToLower(guess)) {
				for i, letter := range jeu.Mot {
					if strings.EqualFold(string(letter), guess) {
						jeu.MotActuel = jeu.MotActuel[:i] + string(letter) + jeu.MotActuel[i+1:]
						jeu.Sortie = "Oui, cette lettre est dans le mot"

					}
				}
			} else {
				jeu.Erreurs++
				jeu.LettresNP = append(jeu.LettresNP, guess)
				jeu.Difference = jeu.MaxErreurs - jeu.Erreurs
				jeu.Sortie = "Non, cette lettre n'est pas dans le mot"

			}
		} else if len(guess) > 1 {

			if strings.EqualFold(guess, jeu.Mot) {
				jeu.MotActuel = jeu.Mot
				jeu.Sortie = "Vous avez gagné ! Le mot était bien " + jeu.Mot
			} else {
				jeu.Erreurs += 2
				jeu.Sortie = "Non ce n'est pas le mot à trouver + 2 erreurs"

			}
		}

		if jeu.Erreurs >= jeu.MaxErreurs {
			jeu.Erreurs = jeu.MaxErreurs
			jeu.Sortie = "Vous avez perdu ! Le mot correct était : " + jeu.Mot
		}

		if !strings.Contains(jeu.MotActuel, "_") {
			jeu.Sortie = "Vous avez gagné ! Le mot était bien" + jeu.Mot

		}
	}

	AfficherTemplate(w, "template", jeu)
}

func AfficherTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplFile := fmt.Sprintf("templates/%s.html", tmpl)
	t, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Traitement(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pseudo := r.FormValue("pseudo")
		difficulty := r.FormValue("difficulty")
		JeuActuel = NouvellePartie(pseudo, difficulty)
		http.Redirect(w, r, "/templates/template", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Init(w http.ResponseWriter, r *http.Request) {
	AfficherTemplate(w, "init", nil)
}

func Menu(w http.ResponseWriter, r *http.Request) {
	AfficherTemplate(w, "menu", nil)
}

func Handle() {
	tpl = template.Must(template.ParseFiles("templates/template.html"))

	http.HandleFunc("/templates/menu", Menu)
	http.HandleFunc("/templates/template", Index)
	http.HandleFunc("/templates/init", Init)
	http.HandleFunc("/templates/treatment", Traitement)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
