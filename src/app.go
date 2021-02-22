package main

import (
	"YusLabCore/src/DataAccess"
	"fmt"
)

func main() {
	fmt.Println("hello world")

	inputDataSheet := DataAccess.ReadFromCsv("data/testcase1/input.csv")
	baseGeneIndex := Calculate(inputDataSheet)
	outputDataSheet := Adjust(inputDataSheet, baseGeneIndex)
	DataAccess.WriteToCsv(outputDataSheet)
}
