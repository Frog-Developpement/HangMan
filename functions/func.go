package functions

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
	"unicode"
)

func Menu() {
	var choice int
	fmt.Println("Que souhaitez vous faire ?")
	fmt.Println("1. Jouer")
	fmt.Println("2. Quitter")
	fmt.Scan(&choice)
	switch choice {
	case 1:
		Hangman()
	case 2:
		return
	}
}

var MotsFromFile []string
var Mistakes int

const MaxMistakes = 6

func Hangman() {
	rand.Seed(time.Now().UnixNano())
	Perdu := false
	Gagné := false
	MotsFromFile = GetMotFromFile()
	RandomMot := MotRandom()
	usedWord = MotToTiré(RandomMot)
	if RandomMot != "" {
		fmt.Println("Mot sélectionné : ", RandomMot)
		fmt.Println("Mot sélectionné : ", usedWord)
	}

	fmt.Println("Quelle lettre souhaitez vous rajouter ?")

	for {

		ChoixLettre(RandomMot)
		fmt.Println("Mot actuel : ", usedWord)
		if usedWord == RandomMot {
			fmt.Println("Félicitations, vous avez trouvé le mot !")
			Gagné = true
		}
		if Mistakes >= MaxMistakes {
			fmt.Println("Vous avez atteint le nombre maximum d'erreurs. Vous avez perdu!")
			Perdu = true
		}
		if Perdu {
			break
		}
		if Gagné {
			break
		}

	}
}

var usedWord string

func AddLetterInWord(letter rune, RandomMot string) {
	usedWordRunes := []rune(usedWord)
	for i, char := range RandomMot {
		if char == letter {
			usedWordRunes[i] = letter
		}
	}
	usedWord = string(usedWordRunes)
}

func ChoixLettre(RandomMot string) {

	for Mistakes < MaxMistakes && usedWord != RandomMot {
		var input string
		fmt.Print("Entrez une lettre : ")
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println("Erreur lors de la saisie de la lettre:", err)
			return
		}

		if len(input) == 0 {
			fmt.Println("Veuillez entrer une lettre valide.")
			return
		}

		letter := rune(input[0])

		letterInWord := false

		letter = unicode.ToLower(letter)

		// Convert the first letter of RandomMot to lowercase for comparison
		if len(RandomMot) > 0 {
			firstLetter := rune(RandomMot[0])
			firstLetter = unicode.ToLower(firstLetter)

			if letter == firstLetter {
				letterInWord = true
				// If the guessed letter matches the lowercase first letter, replace it in usedWord
				usedWordRunes := []rune(usedWord)
				usedWordRunes[0] = rune(RandomMot[0])
				usedWord = string(usedWordRunes)
			}

		}

		for _, char := range RandomMot {
			if char == letter {
				letterInWord = true
				AddLetterInWord(letter, RandomMot)
			}
		}

		fmt.Println("Mot actuel :", usedWord)

		if letterInWord {
			fmt.Println("Cette lettre est bien dans le mot")
		} else {
			fmt.Println("Et non, cette lettre n'est pas dans le mot")
			Mistakes += 1
			fmt.Println("erreurs :", Mistakes)

		}

		AddLetterInWord(letter, RandomMot)

	}
	if Mistakes >= MaxMistakes {
		fmt.Println("Vous avez atteint le nombre maximum d'erreurs. Vous avez perdu!")
		return
	}
}

func MotToTiré(s string) string {
	tirés := make([]rune, len(s))
	for i := range tirés {
		tirés[i] = '_'
	}
	return string(tirés)
}

func MotRandom() string {
	if len(MotsFromFile) == 0 {
		fmt.Println("La liste des mots est vide.")
		return ""
	}

	randomIndex := rand.Intn(len(MotsFromFile))
	RandomMot := MotsFromFile[randomIndex]
	return RandomMot
}

func GetMotFromFile() []string {
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
