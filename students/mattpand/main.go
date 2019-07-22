package main

import (
	"flag"	
	"os"
	"fmt"
	"encoding/csv"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second) 

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		
		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <- answerCh:
			if answer == p.a {
				fmt.Println("Correct!")
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

// TODO:
// - Break up some of above operations in functions in order to make a unit test
// - Write some code that is easier to test


func parseLines (lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}



	return ret
}


type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}


