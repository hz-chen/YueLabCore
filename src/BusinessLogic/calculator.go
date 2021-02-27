package BusinessLogic

import (
	"YusLabCore/src/ObjectModule"
	"fmt"
	"math"
)

func Calculate(inputDataSheet ObjectModule.InputDataSheet) int {

	globalMinSeMean := math.MaxFloat64
	targetWithMinSeMean := -1
	for index, _ := range inputDataSheet.Data {
		currentResult := calculateMinSeMeanAccordingToRow(inputDataSheet.Data, index)
		if currentResult < globalMinSeMean {
			globalMinSeMean = currentResult
			targetWithMinSeMean = index
			fmt.Printf("new target [%v] found with min SeMean [%v]\n", index, currentResult)
		}
	}

	return targetWithMinSeMean
}

func Adjust(inputDataSheet *ObjectModule.InputDataSheet, baseGeneIndex int) *ObjectModule.OutputDataSheet {
	targetRow := inputDataSheet.Data[baseGeneIndex]
	baseNumber := targetRow[0]
	var adjustingFactors []float64
	for _, val := range targetRow {
		adjustingFactors = append(adjustingFactors, val/baseNumber)
	}
	var adjustedDataMatrix [][]float64

	for _, rowData := range inputDataSheet.Data {
		//for each gene:
		var adjustedGene []float64
		for i := 0; i < len(rowData); i++ {
			adjustedGene = append(adjustedGene, rowData[i]*adjustingFactors[i])
		}
		adjustedDataMatrix = append(adjustedDataMatrix, adjustedGene)
		fmt.Printf("adjusted data matrix: %v\n", adjustedGene)
	}

	return &ObjectModule.OutputDataSheet{
		RowTitles:    inputDataSheet.RowTitles,
		ColumnTitles: inputDataSheet.DataColumnTitles,
		Data:         adjustedDataMatrix,
	}

}

//Given a data matrix and a row number, this method will:
//1. calculate the Se/Mean from all other rows according to this row.
//2. among all the calculated Se/Mean, found the minimal one and return
func calculateMinSeMeanAccordingToRow(rawInput [][]float64, index int) float64 {
	minSeMean := math.MaxFloat64

	var calculationDataSheet = ObjectModule.CalculationDataSheet{
		Data:                  divideAccordingToRow(rawInput, index),
		CurrentDividingTarget: index,
	}
	calculationDataSheet = fillInSeArrayAndMeanArray(calculationDataSheet)
	calculationDataSheet = fillInSeMeanArray(calculationDataSheet)

	for _, seMean := range calculationDataSheet.SeMean {
		if seMean < minSeMean {
			minSeMean = seMean
		}
	}

	return minSeMean
}

func divideAccordingToRow(sheet [][]float64, targetIdx int) [][]float64 {
	var result [][]float64
	base := sheet[targetIdx]

	for idx, row := range sheet {

		var currentResult []float64

		for j, d := range row {
			newVal := 0.0
			if idx != targetIdx {
				newVal = d / base[j]
			}
			currentResult = append(currentResult, newVal)
		}

		result = append(result, currentResult)
	}

	return result
}

func fillInSeMeanArray(sheet ObjectModule.CalculationDataSheet) ObjectModule.CalculationDataSheet {
	for idx, se := range sheet.Se {
		mean := sheet.Mean[idx]
		sheet.SeMean = append(sheet.SeMean, se/mean)
	}
	return sheet
}

func fillInSeArrayAndMeanArray(sheet ObjectModule.CalculationDataSheet) ObjectModule.CalculationDataSheet {
	for _, val := range sheet.Data {
		se, mean := calculateSeAndMean(val)
		sheet.Se = append(sheet.Se, se)
		sheet.Mean = append(sheet.Mean, mean)
	}
	return sheet
}

// given an array, calculate the Se and mean of it.
func calculateSeAndMean(arr []float64) (float64, float64) {
	sum := 0.0
	for _, val := range arr {
		sum += val
	}
	mean := sum / float64(len(arr))
	sd := 0.0

	for _, val := range arr {
		sd += math.Pow(val-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(arr)-1))
	se := sd / math.Sqrt(float64(len(arr)))
	return se, mean
}
