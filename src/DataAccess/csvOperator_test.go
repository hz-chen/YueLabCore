package DataAccess

import (
	"YusLabCore/src/ObjectModule"
	"YusLabCore/src/TestUtils"
	"errors"
	"reflect"
	"strings"
	"testing"
)

const InputPath = "../../data/testcase1/input.csv"
const OutputPath = "../../data/testcase1/output.csv"

func Test_isTitleLine(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name           string
		args           args
		expectedResult bool
	}{
		{
			name: "has gene name",
			args: args{
				line: strings.Split("gene 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: true,
		},
		{
			name: "has gene_name",
			args: args{
				line: strings.Split("gene_name 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: true,
		},
		{
			name: "has LibId",
			args: args{
				line: strings.Split("LibId 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: false,
		},
		{
			name: "gene not first column",
			args: args{
				line: strings.Split("LibId gene 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: false,
		},
		{
			name: "has none",
			args: args{
				line: strings.Split("0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTitleLine(tt.args.line); !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("isTitleLine() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

func Test_processTitleRole(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name                  string
		args                  args
		expectedTitleString   []string
		expectedGeneTitleIdx  int
		expectedLibIdTitleIdx int
		expectedError         error
	}{
		{
			name: "with gene name",
			args: args{
				line: strings.Split("gene 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedTitleString: strings.Split("0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
		}, {
			name: "with gene_name",
			args: args{
				line: strings.Split("gene_name 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedTitleString: strings.Split("0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
		},
		{
			name: "with other titles",
			args: args{
				line: strings.Split("randomTitle 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedError: errors.New("input file must have at least one title column"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actualTitleString, err := processTitleRow(tt.args.line);
				!reflect.DeepEqual(actualTitleString, tt.expectedTitleString) ||
					!reflect.DeepEqual(err, tt.expectedError) {
				t.Errorf("processTitleRow() = \t%v, %v,\n \t\t\t\t\t\t\t\twanted = \t%v, %v",
					actualTitleString, err,
					tt.expectedTitleString, tt.expectedError)
			}
		})
	}
}

func Test_readFromCsv(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want *ObjectModule.InputDataSheet
	}{
		{
			name: "test parse title with gene name only",
			args: args{
				path: InputPath,
			},
			want: &ObjectModule.InputDataSheet{
				DataColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
				RowTitles: []ObjectModule.RowTitle{
					{GeneName: "SPBC4F6.10", Index: 0},
					{GeneName: "SPBC4F6.11c", Index: 1},
					{GeneName: "SPBC18A7.01", Index: 2},
					{GeneName: "SPBC530.04", Index: 3},
					{GeneName: "SPBC557.05", Index: 4},
					{GeneName: "SPBC365.20c", Index: 5},
					{GeneName: "SPBC56F2.10c", Index: 6},
					{GeneName: "SPBC577.11", Index: 7},
					{GeneName: "SPBC577.14c", Index: 8},
					{GeneName: "SPBC106.20", Index: 9},
				},
				Data: [][]float64{
					{50, 51, 8, 17, 54},
					{5, 4, 2, 2, 4},
					{143, 159, 44, 67, 165},
					{417, 476, 123, 249, 431},
					{110, 124, 49, 83, 201},
					{249, 452, 81, 142, 310},
					{265, 254, 63, 161, 323},
					{557, 555, 197, 283, 540},
					{1795, 1537, 477, 648, 1483},
					{670, 743, 227, 347, 758},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadFromCsv(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readFromCsv() \t= %v, \n"+
					"\t\t\t\t\t\t\twant \t= %v", got, tt.want)
			}
		})
	}
}

func Test_writeToCsv(t *testing.T) {
	type args struct {
		outputDataSheet *ObjectModule.OutputDataSheet
		outputFileName  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "comparing output file",
			args: args{
				outputDataSheet: &ObjectModule.OutputDataSheet{
					ColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
					RowTitles: []ObjectModule.RowTitle{
						{GeneName: "SPBC4F6.10"}, {GeneName: "SPBC4F6.11c"}, {GeneName: "SPBC18A7.01"},
						{GeneName: "SPBC530.04"}, {GeneName: "SPBC557.05"}, {GeneName: "SPBC365.20c"},
						{GeneName: "SPBC56F2.10c"}, {GeneName: "SPBC577.11"}, {GeneName: "SPBC577.14c"},
						{GeneName: "SPBC106.20"},
					},
					Data: [][]float64{ // from Haijie's hand calculation
						{50, 45.98923284, 23.6123348, 32.82420749, 47.73087071},
						{5, 3.606998654, 5.9030837, 3.86167147, 3.535620053},
						{143, 143.3781965, 129.8678414, 129.3659942, 145.8443272},
						{417, 429.2328398, 363.0396476, 480.778098, 380.9630607},
						{110, 111.8169583, 144.6255507, 160.259366, 177.6649077},
						{249, 407.5908479, 239.0748899, 274.1786744, 274.0105541},
						{265, 229.0444145, 185.9471366, 310.8645533, 285.5013193},
						{557, 500.4710633, 581.4537445, 546.426513, 477.3087071},
						{1795, 1385.989233, 1407.885463, 1251.181556, 1310.831135},
						{670, 670, 670, 670, 670},
					},
				},
				outputFileName: "/tmp/3398ad58-7a44-11eb-ba44-367dda9c6a1b.csv",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteToCsv(tt.args.outputDataSheet, tt.args.outputFileName, false)
			want := ReadFromCsv("../../data/testcase1/output.csv")
			got := ReadFromCsv(tt.args.outputFileName)
			if !TestUtils.AlmostEqualsInputSheet(*want, *got) {
				t.Errorf("WriteToCsv() \t= %v, \n"+
					"\t\t\t\t\t\t\twant \t= %v", *got, *want)
			}
		})
	}
}

func BenchmarkReadFromCsv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ReadFromCsv(InputPath)
	}
}

func BenchmarkWriteToCsv(b *testing.B) {
	outputDataSheet := ObjectModule.OutputDataSheet{
		ColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
		RowTitles: []ObjectModule.RowTitle{
			{GeneName: "SPBC4F6.10", Index: 0},
			{GeneName: "SPBC4F6.11c", Index: 1},
			{GeneName: "SPBC18A7.01", Index: 2},
			{GeneName: "SPBC530.04", Index: 3},
			{GeneName: "SPBC557.05", Index: 4},
			{GeneName: "SPBC365.20c", Index: 5},
			{GeneName: "SPBC56F2.10c", Index: 6},
			{GeneName: "SPBC577.11", Index: 7},
			{GeneName: "SPBC577.14c", Index: 8},
			{GeneName: "SPBC106.20", Index: 9},
		},
		Data: [][]float64{ // from Haijie's hand calculation
			{50, 45.98923284, 23.6123348, 32.82420749, 47.73087071},
			{5, 3.606998654, 5.9030837, 3.86167147, 3.535620053},
			{143, 143.3781965, 129.8678414, 129.3659942, 145.8443272},
			{417, 429.2328398, 363.0396476, 480.778098, 380.9630607},
			{110, 111.8169583, 144.6255507, 160.259366, 177.6649077},
			{249, 407.5908479, 239.0748899, 274.1786744, 274.0105541},
			{265, 229.0444145, 185.9471366, 310.8645533, 285.5013193},
			{557, 500.4710633, 581.4537445, 546.426513, 477.3087071},
			{1795, 1385.989233, 1407.885463, 1251.181556, 1310.831135},
			{670, 670, 670, 670, 670},
		},
		BaseGeneA: ObjectModule.RowTitle{GeneName: "SPBC106.20", Index: 9},
	}

	for i := 0; i < b.N; i++ {
		WriteToCsv(&outputDataSheet, OutputPath, false)
	}
}

func BenchmarkToPrintableFormat(b *testing.B) {
	outputDataSheet := ObjectModule.OutputDataSheet{
		ColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
		RowTitles: []ObjectModule.RowTitle{
			{GeneName: "SPBC4F6.10", Index: 0},
			{GeneName: "SPBC4F6.11c", Index: 1},
			{GeneName: "SPBC18A7.01", Index: 2},
			{GeneName: "SPBC530.04", Index: 3},
			{GeneName: "SPBC557.05", Index: 4},
			{GeneName: "SPBC365.20c", Index: 5},
			{GeneName: "SPBC56F2.10c", Index: 6},
			{GeneName: "SPBC577.11", Index: 7},
			{GeneName: "SPBC577.14c", Index: 8},
			{GeneName: "SPBC106.20", Index: 9},
		},
		Data: [][]float64{ // from Haijie's hand calculation
			{50, 45.98923284, 23.6123348, 32.82420749, 47.73087071},
			{5, 3.606998654, 5.9030837, 3.86167147, 3.535620053},
			{143, 143.3781965, 129.8678414, 129.3659942, 145.8443272},
			{417, 429.2328398, 363.0396476, 480.778098, 380.9630607},
			{110, 111.8169583, 144.6255507, 160.259366, 177.6649077},
			{249, 407.5908479, 239.0748899, 274.1786744, 274.0105541},
			{265, 229.0444145, 185.9471366, 310.8645533, 285.5013193},
			{557, 500.4710633, 581.4537445, 546.426513, 477.3087071},
			{1795, 1385.989233, 1407.885463, 1251.181556, 1310.831135},
			{670, 670, 670, 670, 670},
		},
		BaseGeneA: ObjectModule.RowTitle{GeneName: "SPBC106.20", Index: 9},
	}

	for i := 0; i < b.N; i++ {
		ToPrintableFormat(outputDataSheet)
	}
}
