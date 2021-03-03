package ObjectModule

import (
	"fmt"
)

// Each input will be represented in an InputDataSheet
// For an input sheet:
//  RowTitle.LibId 		| RowTitle.GeneName		| DataColumnTitles[0] 	| DataColumnTitles[1] 	| ...
//  RowTitles[0].LibId 	| RowTitles[0].GeneName | Data[0][0]			| Data[0][1]			| ...
//  RowTitles[1].LibId 	| RowTitles[1].GeneName | Data[1][0]			| Data[1][1]			| ...
//  RowTitles[2].LibId 	| RowTitles[2].GeneName | Data[2][0]			| Data[2][1]			| ...
//  ...
type InputDataSheet struct {
	// titles for data columns ONLY
	DataColumnTitles []string
	// titles for header columns.
	RowTitles []RowTitle
	// data matrix, [0][0] indicates the first data cell, excluding headers columns
	Data [][]float64
}

// including all possible titles of a sheet. Now only support LibId and GeneName
type RowTitle struct {
	// index
	Index int
	// Gene name. For example, SPBC4F6.10
	GeneName string
}

// for each division, build such a matrix
// For each sheet:
//  DataColumnTitles[0] 	| DataColumnTitles[1] 	| ...	| StandardError	| SE/Mean
//  Data[0][0]				| Data[0][1]			| ...	| Se[0]			| SeMean[0]
//  						| 						| ...	| 				| 			 <- CurrentDividingTarget
//  Data[2][0]				| Data[2][1]			| ...	| Se[2]			| SeMean[2]
//  ...
type CalculationDataSheet struct {
	CurrentDividingTarget int
	Data                  [][]float64
	Se                    []float64
	Mean                  []float64
	SeMean                []float64
}

// Result of calculation
// For an output sheet:
//  RowTitle.LibId 		| RowTitle.GeneName		| DataColumnTitles[0] 	| DataColumnTitles[1] 	| ...
//  RowTitles[0].LibId 	| RowTitles[0].GeneName | Data[0][0]			| Data[0][1]			| ...
//  RowTitles[1].LibId 	| RowTitles[1].GeneName | Data[1][0]			| Data[1][1]			| ...
//  RowTitles[2].LibId 	| RowTitles[2].GeneName | Data[2][0]			| Data[2][1]			| ...
//  ...
type OutputDataSheet struct {
	// titles for data columns ONLY
	ColumnTitles []string
	// titles for header columns.
	RowTitles []RowTitle
	// data matrix, [0][0] indicates the first data cell, excluding headers columns
	Data [][]float64

	// base gene used for calculation
	BaseGeneA RowTitle
}

func (o OutputDataSheet) ToPrintableFormat() [][]string {
	var output [][]string
	var titleRow []string

	titleRow = append(titleRow, "gene")
	titleRow = append(titleRow, o.ColumnTitles...)
	output = append(output, titleRow)

	for i := 0; i < len(o.Data); i++ {
		var eachRow []string
		eachRow = append(eachRow, o.RowTitles[i].GeneName)
		for _, val := range o.Data[i] {
			eachRow = append(eachRow, fmt.Sprintf("%f", val))
		}
		output = append(output, eachRow)
	}

	return output
}
