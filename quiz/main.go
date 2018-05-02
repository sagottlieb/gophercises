package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	defaultFilename         = "problems.csv"
	defaultQuizLengthInSec  = 30
	defaultRandomizeEnabled = false
)

func main() {

	var filename = flag.String("filename", defaultFilename, "CSV file containing quiz in the format 'question,answer'")
	var quizLengthInSec = flag.Int("length", defaultQuizLengthInSec, "Number of seconds to complete quiz")
	var randomizeEnabled = flag.Bool("randomize", defaultRandomizeEnabled, "Randomize the order of the questions in the quiz")
	flag.Parse()

	lines := readCSVRecords(*filename)

	problems := parseRecords(lines)

	if *randomizeEnabled {
		problems = randomizeProblemOrder(problems)
	}

	fmt.Printf("You will have %d seconds to complete the quiz. Press enter to begin.\n", *quizLengthInSec)
	var anything string
	fmt.Scanf("%s", &anything)

	timer := time.NewTimer(time.Duration(*quizLengthInSec) * time.Second)

	numCorrect := 0

	for i, p := range problems {

		fmt.Printf("Problem #%d: %s = ?\n", i+1, p.question)

		answerChan := make(chan string)
		go waitForAnswer(answerChan)

		select {

		case <-timer.C:
			fmt.Println("\nTime's up!")
			printScore(numCorrect, len(lines))
			os.Exit(0)

		case answer := <-answerChan:
			if answer == p.answer {
				numCorrect++
			}

		}
	}

	printScore(numCorrect, len(lines))
	os.Exit(0)

}

type problem struct {
	question string
	answer   string
}

func readCSVRecords(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		exit(fmt.Sprintf("Error opening file '%s': %s\n", filename, err.Error()))
	}

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Error reading CSV data from file '%s': %s\n", filename, err.Error()))
	}

	file.Close()

	return lines
}

func parseRecords(lines [][]string) []problem {
	problems := []problem{}

	for _, x := range lines {

		// I'm not sure that ths check is even necessary
		if len(x) != 2 {
			fmt.Printf("Skipping CSV record with unexpected format: '%s'\n, x")
			continue
		}

		p := problem{
			question: sanitizeInput(x[0]),
			answer:   sanitizeInput(x[1]),
		}

		problems = append(problems, p)

	}

	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func printScore(correct, total int) {
	fmt.Printf("Your score: %d/%d\n", correct, total)
}

func randomizeProblemOrder(in []problem) []problem {
	rand.Seed(time.Now().Unix())

	var out = make([]problem, len(in))

	// gives a new ordering of the valid indexes into in
	indexPerm := rand.Perm(len(in))

	// if indexPerm[i] = j, then store in[i] at index j of out

	for i, j := range indexPerm {
		out[j] = in[i]
	}

	return out
}

func waitForAnswer(answerCh chan string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	answerCh <- answer //sanitizeInput(answer)
}

func sanitizeInput(in string) string {
	out := strings.TrimSpace(in)
	out = strings.ToLower(out)
	return out
}
