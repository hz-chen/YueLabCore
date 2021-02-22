package main

import (
	"YusLabCore/src/ObjectModule"
)

func Calculate(inputDataSheet *ObjectModule.InputDataSheet) int {

	globalMinSeMean := 0.0
	targetWithMinSeMean := -1
	for index, _ := range inputDataSheet.Data {
		currentResult := calculateDataSheetForRowNumber(inputDataSheet.Data, index)
		if currentResult < globalMinSeMean {
			globalMinSeMean = currentResult
			targetWithMinSeMean = index
		}
	}

	return targetWithMinSeMean
}

func Adjust(inputDataSheet *ObjectModule.InputDataSheet, baseGeneIndex int) *ObjectModule.OutputDataSheet {
	return nil
}



func calculateDataSheetForRowNumber(sheet [][]float64, index int) float64 {

	return 0.0
}
