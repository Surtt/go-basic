package main

import "fmt"

func main() {
	const usdToEur = 0.91
	const usdToRub = 84.44
	const eurToRub = usdToRub / usdToEur

	input := usersInput()
	fmt.Println(input)
}

func usersInput() string {
	var input string
	fmt.Scan(&input)
	return input
}

func getValue(amount, from, to string) {

}
