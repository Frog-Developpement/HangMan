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

var game Game

type Game struct {
	Word             string
	CurrentWord      string
	IncorrectGuesses int
	MaxIncorrect     int
	Difference       int
	Difficulte       string
	IncorrectLetters []string
	Outcome          string
	Pseudo           string
}

var currentGame *Game
var tpl *template.Template

func getRandomWord() (string, error) {
	file, err := os.Open("noms_monstres.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	return words[rand.Intn(len(words))], nil
}

func getMaxIncorrect(difficulty string) int {
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

func startNewGame(pseudo, difficulty string) *Game {
	wordToGuess, err := getRandomWord()
	if err != nil {
		panic(err)
	}

	maxIncorrect := getMaxIncorrect(difficulty)

	return &Game{
		Word:             wordToGuess,
		CurrentWord:      strings.Repeat("_", len(wordToGuess)),
		IncorrectGuesses: 0,
		MaxIncorrect:     maxIncorrect,
		Difficulte:       difficulty,
		IncorrectLetters: []string{},
		Outcome:          "",
		Pseudo:           pseudo,
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		currentGame = startNewGame(currentGame.Pseudo, currentGame.Difficulte)

	}

	game := currentGame

	if r.Method == http.MethodPost {
		guess := strings.ToLower(r.FormValue("guess"))

		if len(guess) == 1 {
			if strings.Contains(strings.ToLower(game.Word), strings.ToLower(guess)) {
				for i, letter := range game.Word {
					if strings.EqualFold(string(letter), guess) {
						game.CurrentWord = game.CurrentWord[:i] + string(letter) + game.CurrentWord[i+1:]
						game.Outcome = "Oui, cette lettre est dans le mot"

					}
				}
			} else {
				game.IncorrectGuesses++
				game.IncorrectLetters = append(game.IncorrectLetters, guess)
				game.Difference = game.MaxIncorrect - game.IncorrectGuesses
				game.Outcome = "Non, cette lettre n'est pas dans le mot"

			}
		} else if len(guess) > 1 {
			// Devinez un mot
			if strings.EqualFold(guess, game.Word) {
				game.CurrentWord = game.Word
				game.Outcome = "Vous avez gagné ! Le mot était bien " + game.Word
			} else {
				game.IncorrectGuesses += 2
				game.Outcome = "Non ce n'est pas le mot à trouver + 2 erreurs"

			}
		}

		if game.IncorrectGuesses >= game.MaxIncorrect {
			game.Outcome = "Vous avez perdu ! Le mot correct était : " + game.Word
		}

		if !strings.Contains(game.CurrentWord, "_") {
			game.Outcome = "Vous avez gagné ! Le mot était bien" + game.Word

		}
	}

	renderTemplate(w, "template", game)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
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

func treatmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pseudo := r.FormValue("pseudo")
		difficulty := r.FormValue("difficulty")
		currentGame = startNewGame(pseudo, difficulty)
		http.Redirect(w, r, "/templates/template", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "init", nil)
}

func main() {
	tpl = template.Must(template.ParseFiles("templates/template.html"))

	http.HandleFunc("/templates/template", indexHandler)
	http.HandleFunc("/templates/init", initHandler)
	http.HandleFunc("/templates/treatment", treatmentHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
