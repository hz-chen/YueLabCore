package main

import (
	"YusLabCore/src/BusinessLogic"
	"YusLabCore/src/DataAccess"
	"testing"
)

const inputPath = "../data/testcase1/input.csv"
const outputPath = "/tmp/output.csv"

func BenchmarkInput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DataAccess.ReadFromCsv(inputPath)
	}
}

func BenchmarkCalculate(b *testing.B) {

	inputDataSheet := DataAccess.ReadFromCsv(inputPath)

	for i := 0; i < b.N; i++ {
		_ = BusinessLogic.Calculate(*inputDataSheet)
	}
}

func BenchmarkAdjust(b *testing.B) {

	inputDataSheet := DataAccess.ReadFromCsv(inputPath)
	baseGene := BusinessLogic.Calculate(*inputDataSheet)
	for i := 0; i < b.N; i++ {
		_ = BusinessLogic.Adjust(inputDataSheet, baseGene.Index)
	}

}

func BenchmarkOutput(b *testing.B) {

	inputDataSheet := DataAccess.ReadFromCsv(inputPath)
	baseGene := BusinessLogic.Calculate(*inputDataSheet)
	outputDataSheet := BusinessLogic.Adjust(inputDataSheet, baseGene.Index)
	for i := 0; i < b.N; i++ {
		DataAccess.WriteToCsv(outputDataSheet, outputPath, true)
	}
}
