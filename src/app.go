package main

import (
	"YusLabCore/src/BusinessLogic"
	"YusLabCore/src/DataAccess"
)

func main() {
	performOperation("/Users/hongzhoc/Documents/input_2.csv", "/Users/hongzhoc/Documents/result_2_2.csv")
	performOperation("/Users/hongzhoc/Documents/input_2_deseq2.csv", "/Users/hongzhoc/Documents/result_2_deseq2_2.csv")
}

func performOperation(inputFileName string, outputFileName string) {
	inputDataSheet := DataAccess.ReadFromCsv(inputFileName)
	baseGene := BusinessLogic.Calculate(*inputDataSheet)
	outputDataSheet := BusinessLogic.Adjust(inputDataSheet, baseGene.Index)
	DataAccess.WriteToCsv(outputDataSheet, outputFileName, true)

}
