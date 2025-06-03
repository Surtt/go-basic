package main

import (
	"errors"
	"fmt"
	"strconv"
)

const usdToEur = 0.91
const usdToRub = 84.44
const eurToRub = usdToRub / usdToEur

var conversions = map[string]func(float64) float64{
	"USD->EUR": func(a float64) float64 { return a * usdToEur },
	"USD->RUB": func(a float64) float64 { return a * usdToRub },
	"EUR->USD": func(a float64) float64 { return a / usdToEur },
	"EUR->RUB": func(a float64) float64 { return a * eurToRub },
	"RUB->USD": func(a float64) float64 { return a / usdToRub },
	"RUB->EUR": func(a float64) float64 { return a / eurToRub },
}

func main() {
	currency := askCurrency()
	amount := askAmount()
	targetCurrency := askTargetCurrency(currency)
	converted, err := convert(amount, currency, targetCurrency, &conversions)
	if err != nil {
		fmt.Println("Error during conversion", err)
		return
	}

	fmt.Printf("You want to convert %.2f %s в %s\n", amount, currency, targetCurrency)
	fmt.Printf("Result conversion: %.2f %s = %.2f %s\n", amount, currency, converted, targetCurrency)

}

func usersInput(prompt string) (string, error) {
	var input string
	fmt.Println(prompt)
	_, err := fmt.Scan(&input)
	return input, err
}

func askCurrency() string {
	validCurrencies := map[string]bool{"USD": true, "EUR": true, "RUB": true}
	for {
		currency, err := usersInput("Enter a valid currency value (USD, EUR, RUB): ")

		if err != nil {
			fmt.Println("Input error", err)
			continue
		}

		if !validCurrencies[currency] {
			err := errors.New("invalid currency")
			fmt.Println(err)
			continue
		}

		return currency
	}
}

func askAmount() float64 {

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

		options := map[string]bool{}
		for _, currency := range []string{"USD", "EUR", "RUB"} {
			if currency != from {
				options[currency] = true
				fmt.Print(currency + " ")
			}
		}

		fmt.Println()

		targetCurrency, err := usersInput("Enter a target currency: ")

		if err != nil {
			fmt.Println("Input error", err)
			continue
		}

		if !options[targetCurrency] {
			err := errors.New("invalid target currency")
			fmt.Println(err)
			continue
		}

		return targetCurrency
	}
}

func convert(amount float64, from, to string, conversions *map[string]func(float64) float64) (float64, error) {
	key := from + "->" + to

	if fn, ok := (*conversions)[key]; ok {
		return fn(amount), nil
	}

	return 0, errors.New("invalid conversion")
}
