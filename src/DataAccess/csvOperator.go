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
const LIBID_NAME_TITLE = "LibId"

func WriteToCsv(outputDataSheet *ObjectModule.OutputDataSheet, fileName string) {
	file, err := os.Create(fileName)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range outputDataSheet.ToPrintableFormat() {
		fmt.Printf("printing line %v\n", value)
		err := writer.Write(value)
		checkError("Cannot write to file", err)
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

	geneTitleIdx := -1

	fmt.Printf("processing first line as title: %s\n", titleLine)
	if isTitleLine(titleLine) {
		dataSheet.DataColumnTitles, geneTitleIdx, _, err = processTitleRow(titleLine)
		checkError("Failure to process title line ", err)
	}

	if geneTitleIdx == -1 {
		log.Fatal("no gene found in title index, unsupported file format!")
		return nil
	}

	// Iterate through the records
	for {
		// Read each eachLine from csv
		eachLine, err := r.Read()
		if err == io.EOF {
			break
		}
		checkError("Failure to process data line "+strings.Join(eachLine, ","), err)

		log.Printf("processing data line: %v\n", eachLine)

		// read row title
		dataSheet.RowTitles = append(dataSheet.RowTitles, ObjectModule.RowTitle{
			GeneName: eachLine[geneTitleIdx],
		})

		fmt.Printf("updated row title %v\n", dataSheet.RowTitles)

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
	return strings.Contains(line[0], GENE_NAME_TITLE) || strings.Contains(line[0], LIBID_NAME_TITLE)
}

// first return value is the processed column titles for data columns only;
// second return value is the index of gene title index
// third return value is the index of the libId title index
func processTitleRow(line []string) ([]string, int, int, error) {

	geneTitleIdx := getIndex(line, GENE_NAME_TITLE)
	libIdIdx := getIndex(line, LIBID_NAME_TITLE)

	hasGeneTitle := geneTitleIdx != -1
	hasLibIdTitle := libIdIdx != -1

	var dataColumnStartingIdx = 0
	if hasGeneTitle && hasLibIdTitle {
		dataColumnStartingIdx = 2
	} else if hasGeneTitle || hasLibIdTitle {
		dataColumnStartingIdx = 1
	} else {
		return nil, -1, -1, errors.New("input file must have at least one title column")
	}

	return line[dataColumnStartingIdx:], geneTitleIdx, libIdIdx, nil
}

func getIndex(s []string, e string) int {
	for idx, a := range s {
		if strings.Contains(a, e) {
			return idx
		}
	}
	return -1
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
