package functions

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unicode"
)

func Menu() {
	ClearConsole()
	var choice int
	fmt.Print("Bienvenue dans le jeu du pendu !\n")
	fmt.Printf("\nLe but est de trouver un mot en entrant les lettres que vous pensez être présentes dans le mot. \nVous pouvez aussi entrer le mot directement si vous pensez que c'est le bon.\nAttention aux erreurs ce serait une erreur !\n")
	fmt.Println("\nQue souhaitez vous faire ?")
	fmt.Println("\n1. Lancer une partie")
	fmt.Println("2. Quitter")
	fmt.Scan(&choice)
	switch choice {
	case 1:
		ClearConsole()
		Hangman()
	case 2:
		ClearConsole()
		println("Au revoir ...")
		os.Exit(1)
	}
}

var MotsFromFile []string
var Mistakes int

const MaxMistakes = 10

func Hangman() {
	rand.Seed(time.Now().UnixNano())
	Perdu := false
	Gagné := false
	MotsFromFile = GetMotFromFile()
	RandomMot := MotRandom()
	usedWord = MotToTire(RandomMot)
	if RandomMot != "" {

		fmt.Println("\nMot à trouver : ", usedWord)
	}

	for {

		ChoixLettre(RandomMot)
		if usedWord == RandomMot {
			fmt.Println("Félicitations, vous avez trouvé le mot !")
			Gagné = true
		}
		if Mistakes >= MaxMistakes {
			fmt.Println("Vous avez atteint le nombre maximum d'erreurs. Vous êtes mort !")
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

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func ChoixLettre(RandomMot string) {
	var SUsedLetter []string
	LetterInWord := false
	for Mistakes < MaxMistakes && usedWord != RandomMot {
		var choice int
		fmt.Println("\nQue voulez-vous faire ?")
		fmt.Println("\n1. Choisir une lettre")
		fmt.Println("2. Saisir le mot entier")
		fmt.Scan(&choice)
		switch choice {
		case 1:

			var input string
			fmt.Print("\nEntrez une lettre : ")
			_, err := fmt.Scan(&input)
			if err != nil {
				fmt.Println("Erreur lors de la saisie de la lettre:", err)
				return
			}
			if len(input) == 0 || len(input) > 1 || !unicode.IsLetter([]rune(input)[0]) {
				fmt.Println("Veuillez entrer une lettre valide.")
				return

			} else {
				input = strings.ToLower(input)
				letter := rune(input[0])
				LetterInWord = false
				for _, char := range RandomMot {
					if unicode.ToLower(char) == letter {
						LetterInWord = true
						AddLetterInWord(char, RandomMot)
					}
				}

				if LetterInWord {
					ClearConsole()
					fmt.Println("Cette lettre est bien dans le mot")
					fmt.Println("Mot actuel :", usedWord)
					fmt.Println("Lettres saisies et non présentes dans le mot :", SUsedLetter)
				} else {
					ClearConsole()
					fmt.Print("Et non, cette lettre n'est pas dans le mot\n\n")
					Mistakes++
					fmt.Println("Mot actuel :", usedWord)
					if !contains(SUsedLetter, input) {
						SUsedLetter = append(SUsedLetter, input)
						fmt.Println("Lettres utilisées et non présentes dans le mot :", SUsedLetter)

					} else {
						fmt.Println("Vous avez déjà essayé cette lettre mais n'y étais pas...")
					}

					hangman := GetHangman(Mistakes)
					fmt.Println(hangman)
				}
			}

		case 2:
			ClearConsole()
			fmt.Println("Mot actuel :", usedWord)
			fmt.Println("Lettres saisies et non présentes dans le mot :", SUsedLetter)
			var input string
			fmt.Print("\nEntrez un mot : ")

			_, err := fmt.Scan(&input)
			if err != nil {
				fmt.Println("Erreur lors de la saisie de la lettre:", err)
				return
			}

			if len(input) < 4 {
				ClearConsole()
				fmt.Println("Veuillez entrer un mot valide.")
				return
			} else if len(input) == len(RandomMot) {
				if strings.EqualFold(input, RandomMot) {
					usedWord = RandomMot
				} else {
					ClearConsole()
					Mistakes += 2
					println("Perdu ! Ce n'était pas le bon mot. Vous avez perdu 2 points.")
					fmt.Println("Mot actuel :", usedWord)
					fmt.Println("Lettres saisies et non présentes dans le mot :", SUsedLetter)
					hangman := GetHangman(Mistakes)
					fmt.Println(hangman)
				}
			} else {
				ClearConsole()
				Mistakes += 2
				println("Perdu ! Ce n'était pas le bon mot. Vous avez perdu 2 points.")
				fmt.Println("Mot actuel :", usedWord)
				hangman := GetHangman(Mistakes)
				fmt.Println(hangman)
			}

		}
	}
}
func MotToTire(s string) string {
	rand.Seed(time.Now().Unix())
	lettreR := rand.Intn(len(s))
	tires := make([]rune, len(s))
	for i := range tires {
		tires[i] = '_'
		if lettreR == i {
			tires[i] = rune(s[lettreR])
		}
	}
	return string(tires)
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

func ClearConsole() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
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
func GetHangman(mistakes int) string {
	hangmanFigures := []string{
		`
=========
`, `	
      |
      |
      |
      |
      |
=========
`, `
  +---+
      |
      |
      |
      |
      |
=========
`, `
  +---+
  |   |
      |
      |
      |
      |
=========
`, `
  +---+
  |   |
  O   |
      |
      |
      |
=========
`, `
  +---+
  |   |
  O   |
  |   |
      |
      |
=========
`, `
  +---+
  |   |
  O   |
 /|   |
      |
      |
=========
`, `
  +---+
  |   |
  O   |
 /|\  |
      |
      |
=========
`, `
  +---+
  |   |
  O   |
 /|\  |
 /    |
      |
=========
`, `
  +---+
  |   |
  O   |
 /|\  |
 / \  |
      |
=========
`}
	if mistakes >= 0 && mistakes <= len(hangmanFigures) {
		return hangmanFigures[(mistakes - 1)]
	}
	return "erreur"
}
