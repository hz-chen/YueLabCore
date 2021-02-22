package ObjectModule

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
	// Lib id. For example, 10-A01
	LibId string
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

	// Gene A in the immutable pair
	immutableGeneA RowTitle
	// Gene A in the immutable pair
	immutableGeneB RowTitle
}
