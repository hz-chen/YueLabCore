package DataAccess

import (
	"YusLabCore/src/ObjectModule"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const GENE_NAME_TITLE = "gene"

func WriteToCsv(outputDataSheet *ObjectModule.OutputDataSheet, fileName string, appendBenchmarkGene bool) {
	file, err := os.Create(fileName)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range outputDataSheet.ToPrintableFormat() {
		//fmt.Printf("printing line %v\n", value)
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}

	//write base gene name and index
	if appendBenchmarkGene {
		// Our index start with 0. In the output, excel rows start at 1, and the first row is the title row.
		// Therefore the exact row number should be index + 2
		err = writer.Write([]string{"BaseGene", outputDataSheet.BaseGeneA.GeneName, fmt.Sprintf("%v", outputDataSheet.BaseGeneA.Index + 2)})
		checkError("Failed to write base gene to file", err)
	}

}

func ReadFromCsv(path string) *ObjectModule.InputDataSheet {
	csvFile, err := os.Open(path)

	checkError("Failure to open input file "+path, err)
	// Parse the file
	r := csv.NewReader(csvFile)

	dataSheet := ObjectModule.InputDataSheet{
		DataColumnTitles: nil,
		RowTitles:        nil,
		Data:             nil,
	}

	// read & process first line
	titleLine, err := r.Read()

	checkError("Failure to read input file "+path, err)

	geneTitleIdx := 0

	//fmt.Printf("processing first line as title: %s\n", titleLine)
	if isTitleLine(titleLine) {
		dataSheet.DataColumnTitles, err = processTitleRow(titleLine)
		checkError("Failure to process title line ", err)
	}

	recordIdx := 0
	// Iterate through the records
	for {
		// Read each eachLine from csv
		eachLine, err := r.Read()
		if err == io.EOF {
			break
		}
		checkError("Failure to process data line "+strings.Join(eachLine, ","), err)

		//log.Printf("processing data line: %v\n", eachLine)

		// read row title
		dataSheet.RowTitles = append(dataSheet.RowTitles, ObjectModule.RowTitle{
			GeneName: eachLine[geneTitleIdx],
			Index:    recordIdx,
		})
		recordIdx++
		//fmt.Printf("updated row title %v\n", dataSheet.RowTitles)

		// read data cells
		var currentDataRow []float64

		for _, s := range eachLine[1:] {

			f, err := strconv.ParseFloat(s, 32)
			if err != nil {
				log.Fatalf("failure to convert %v to float number", s)
			}
			currentDataRow = append(currentDataRow, f)
		}
		dataSheet.Data = append(dataSheet.Data, currentDataRow)
	}

	return &dataSheet
}

func isTitleLine(line []string) bool {
	return strings.Contains(line[0], GENE_NAME_TITLE)
}

//always assume line[0] contains "gene"
func processTitleRow(line []string) ([]string, error) {

	if !isTitleLine(line) {
		return nil, errors.New("input file must have at least one title column")
	}

	return line[1:], nil
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
