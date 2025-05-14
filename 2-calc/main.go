package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var operations = map[string]func([]int) float64{
	"SUM": sum,
	"AVG": average,
	"MED": median,
}

func main() {
	operation := askOperation()
	numbers := askNumbers()
	result := operations[operation](numbers)
	fmt.Println(result)
}

func askOperation() string {

	for {
		operation, err := usersInput("Enter a valid operation: ")

		if _, ok := operations[operation]; !ok {
			err = errors.New("Invalid operation")
			fmt.Println(err)
			continue
		}

		return operation
	}
}

func askNumbers() []int {
	var numbers []int
	valid := true
	for {
		numbersStr, err := usersInput("Enter numbers separated by commas (e.g. 1,2,3): ")

		if err != nil {
			fmt.Println("Invalid input: ", err)
			continue
		}
		parts := strings.Split(numbersStr, ",")
		fmt.Println(strings.TrimSpace(numbersStr))

		for _, part := range parts {
			numStr := strings.TrimSpace(part)
			num, err := strconv.Atoi(numStr)

			if err != nil {
				err = errors.New("Invalid numbers")
				fmt.Println(err, numStr)
				valid = false
				break
			}
			numbers = append(numbers, num)

		}

		if !valid {
			continue
		}

		return numbers
	}
}

func usersInput(prompt string) (string, error) {
	var input string
	fmt.Println(prompt)
	_, err := fmt.Scan(&input)
	return input, err
}

func sum(numbers []int) float64 {
	var total float64

	for _, n := range numbers {
		total += float64(n)
	}

	return total
}

func average(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}

	return float64(sum(numbers)) / float64(len(numbers))
}

func median(numbers []int) float64 {
	n := len(numbers)
	if n == 0 {
		return 0
	}

	sorted := make([]int, n)
	copy(sorted, numbers)
	sort.Ints(sorted)

	middle := n / 2
	if n%2 == 0 {
		return float64(sorted[middle-1]+sorted[middle]) / 2
	} else {
		return float64(sorted[middle])
	}
}

func selectCalculation(operation string, numbers []int) float64 {
	switch operation {
	case "SUM":
		return float64(sum(numbers))
	case "AVG":
		return average(numbers)
	case "MED":
		return median(numbers)
	default:
		fmt.Println("Unknown operation")
		return 0
	}
}
