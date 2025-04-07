package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	const usdToEur = 0.91
	const usdToRub = 84.44
	const eurToRub = usdToRub / usdToEur

	currency := askCurrency()
	amount := askAmmount()
	targetCurrency := askTargetCurrency(currency)
	converted, err := convert(amount, currency, targetCurrency)
	if err != nil {
		fmt.Println("Error during conversion", err)
		return
	}

	fmt.Printf("You want to convert %.2f %s в %s\n", amount, currency, targetCurrency)
	fmt.Printf("Result: %.2f %s = %.2f %s\n", amount, currency, converted, targetCurrency)

}

func usersInput(prompt string) (string, error) {
	var input string
	fmt.Println(prompt)
	_, err := fmt.Scan(&input)
	return input, err
}

func askCurrency() string {
	for {
		currency, err := usersInput("Enter a valid currency value (USD, EUR, RUB): ")

		if err != nil {
			fmt.Println("Input error", err)
			continue
		}

		if currency != "USD" && currency != "EUR" && currency != "RUB" {
			err := errors.New("invalid currency")
			fmt.Println(err)
			continue
		}

		return currency
	}
}

func askAmmount() float64 {

	for {
		amountStr, err := usersInput("Enter a valid ammount: ")

		if err != nil {
			fmt.Println("Input error", err)
			continue
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			err := errors.New("invalid ammount")
			fmt.Println(err)
			continue
		}

		return amount
	}
}

func askTargetCurrency(from string) string {
	for {
		fmt.Println("Available currencies: ")

		options := []string{}
		for _, currency := range []string{"USD", "EUR", "RUB"} {
			if currency != from {
				options = append(options, currency)
			}
		}

		fmt.Println(strings.Join(options, ", "))

		targetCurrency, err := usersInput("Enter a target currency: ")

		if err != nil {
			fmt.Println("Input error", err)
			continue
		}

		valid := false
		for _, option := range options {
			if targetCurrency == option {
				valid = true
				break
			}
		}

		if !valid {
			err := errors.New("invalid target currency")
			fmt.Println(err)
			continue
		}

		return targetCurrency
	}
}

func convert(amount float64, from, to string) (float64, error) {
	const usdToEur = 0.91
	const usdToRub = 84.44
	const eurToRub = usdToRub / usdToEur

	switch from + "->" + to {
	case "USD->EUR":
		return amount * usdToEur, nil
	case "USD->RUB":
		return amount * usdToRub, nil
	case "EUR->USD":
		return amount / usdToEur, nil
	case "EUR->RUB":
		return amount * eurToRub, nil
	case "RUB->USD":
		return amount / usdToRub, nil
	case "RUB->EUR":
		return amount / eurToRub, nil
	default:
		return 0, errors.New("invalid conversion")
	}

}
