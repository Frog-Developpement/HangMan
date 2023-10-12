package jeu

import (
	"bufio"
	"fmt"
	"os"
)

type Hangman struct {
	Mots     []Mot
	Position int
}

type Mot struct {
	NbLettres         int
	NbLettresProposes int
}

func MotRandom() {

}

func getMotFromFile() []int {
	readFile, err := os.ReadFile("noms_monstres.txt")
	if err != nil {
		fmt.Println("Erreur de lecture")
	} else {
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)
		var lines []string
		for fileScanner.Scan() {
			lines = append(lines, fileScanner.Text())

		}
		text := string
		readFile.Close()
		for _, line := range lines {
			fmt.Println("The name is :", line)
		}
	}
}
