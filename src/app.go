package main

import (
	"YusLabCore/src/BusinessLogic"
	"YusLabCore/src/DataAccess"
	"fmt"
)

func main() {
	fmt.Println("hello world")

	inputDataSheet := DataAccess.ReadFromCsv("data/testcase1/input.csv")
	baseGeneIndex := BusinessLogic.Calculate(*inputDataSheet)
	outputDataSheet := BusinessLogic.Adjust(inputDataSheet, baseGeneIndex)
	DataAccess.WriteToCsv(outputDataSheet, "result.csv")
}
