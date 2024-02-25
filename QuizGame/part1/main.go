package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type problem struct {
	equation string
	answer   int
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
		correctAns, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Printf("The correct answer is not in correct format\n")
			continue
		}

		problemList = append(problemList, problem{equation, correctAns})
	}
	return problemList
}

func main() {
	var points int = 0

	flag.String("csv", "problems.csv", "a csv file in the format of 'question answer")
	shuffle := flag.Bool("shuffle", false, "if true, shuffle the questions each time")
	flag.Parse()

	problemList := parseFile("problems.csv")

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problemList), func(i, j int) {
			problemList[i], problemList[j] = problemList[j], problemList[i]
		})
	}

	for _, problem := range problemList {
		fmt.Printf("%s: ", problem.equation)
		var num int
		_, err2 := fmt.Scanf("%d", &num)
		if err2 != nil {
			fmt.Printf("Error: %v\n", err2)
		} else if num != int(problem.answer) {
			fmt.Println("Answer is incorrect")
		} else {
			points += 1
		}
	}

	fmt.Printf("You do %d out of %d questions correctly\n", points, len(problemList))
}
