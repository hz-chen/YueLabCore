package main

import (
	"YusLabCore/src/DataAccess"
	"fmt"
)

func main() {
	fmt.Println("hello world")

	inputDataSheet := DataAccess.ReadFromCsv("data/testcase1/input.csv")
	baseGeneIndex := calculate(inputDataSheet)
	outputDataSheet := adjust(inputDataSheet, baseGeneIndex)
	DataAccess.WriteToCsv(outputDataSheet)
}
