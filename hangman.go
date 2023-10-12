package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	menu()

}

func menu() {
	var choice int
	fmt.Println("Que souhaitez vous faire ?")
	fmt.Println("1. Jouer")
	fmt.Println("2. Quitter")
	fmt.Scan(&choice)
	switch choice {
	case 1:
		Hangman(lettresRandom)
	case 2:
		return
	}
}

type HangmanS struct {
	Mots     []Mot
	Position int
}

type Mot struct {
	NbLettres         int
	NbLettresProposes int
}

var MotFromFile []string
var lettresRandom string

func Hangman(s string) {
	rand.Seed(time.Now().UnixNano())
	MotFromFile = getMotFromFile()

	randomMot := MotRandom()
	if randomMot != "" {
		lettresRandom = PrintLettersWord(randomMot)
		indices := PrintLettersWord(randomMot)

		PrintRemainingLetters(randomMot, indices)

	}

	fmt.Println("Quelle lettre souhaitez vous rajouter ?")
	var choice string
	fmt.Scan(&choice)
	if choice == s {
		fmt.Println()
	}
}

func PrintRemainingLetters(s string, indices []int) {
	RemainingLetters := ""
	for i, char := range s {
		if !contains(indices, i) {
			RemainingLetters += string(char)
		}
	}
	fmt.Println("Lettres restantes :", RemainingLetters)
}

func contains(arr []int, val int) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func PrintLettersWord(s string) []int {
	lettresRandom := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	indices := r.Perm(len(s))[:2]

	for i, char := range s {
		if contains(indices, i) {
			fmt.Printf("%c ", char)
			lettresRandom += string(char)
		} else {
			fmt.Print("_")
		}
	}
	fmt.Println() // Nouvelle ligne après avoir affiché les lettres
	return indices
}

func MotRandom() string {
	if len(MotFromFile) == 0 {
		fmt.Println("La liste des mots est vide.")
		return ""
	}

	randomIndex := rand.Intn(len(MotFromFile))
	randomMot := MotFromFile[randomIndex]
	return randomMot
}

func getMotFromFile() []string {
	readFile, err := os.ReadFile("noms_monstres.txt")
	if err != nil {
		fmt.Println("Erreur de lecture")
		return nil
	}
	fileScanner := bufio.NewScanner(bytes.NewReader(readFile))
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	return lines
}
