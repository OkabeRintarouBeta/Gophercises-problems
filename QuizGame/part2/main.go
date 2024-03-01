package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	equation string
	answer   string
}

func printError(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}

func stripEmptyFields(record *[]string) {
	idx := len(*record) - 1
	for idx >= 0 {
		if len((*record)[idx]) == 0 {
			idx--
		} else {
			break
		}
	}
	*record = (*record)[:idx+1]
}

func parseFile(filepath string) []problem {

	problemList := []problem{}
	csvFile, err := os.Open(filepath)
	if err != nil {
		printError(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Disable the check for a consistent number of fields
	_, err1 := reader.Read()
	if err1 != nil {
		printError(err1)
	}

	for {
		record, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break // Break the loop when there are no more records to read
			}
			fmt.Println("Error reading CSV:", err)
			continue
		}
		stripEmptyFields(&record)
		if len(record) != 2 {
			fmt.Println(record)
			fmt.Println("Does not have the correct number of parameters!")
			continue
		}

		equation := record[0]
		correctAns := record[1]

		problemList = append(problemList, problem{equation, correctAns})
	}
	return problemList
}

func main() {
	var points int = 0

	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question answer")
	shuffle := flag.Bool("shuffle", false, "if true, shuffle the questions each time")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	problemList := parseFile(*csvFileName)

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problemList), func(i, j int) {
			problemList[i], problemList[j] = problemList[j], problemList[i]
		})
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemloop:
	for _, problem := range problemList {
		fmt.Printf("%s: ", problem.equation)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answer = strings.TrimSpace(answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			break problemloop
		case answer := <-answerCh:
			if answer == problem.answer {
				points += 1
			} else {
				fmt.Println("Answer is incorrect")
			}
		}
	}
	fmt.Printf("\nYou do %d out of %d questions correctly\n", points, len(problemList))
}
