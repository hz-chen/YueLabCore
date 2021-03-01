package DataAccess

import (
	"YusLabCore/src/ObjectModule"
	"errors"
	"reflect"
	"strings"
	"testing"
)

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
			name: "has LibId",
			args: args{
				line: strings.Split("LibId 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: true,
		},
		{
			name: "has both",
			args: args{
				line: strings.Split("LibId gene 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: true,
		},
		{
			name: "has both reversed",
			args: args{
				line: strings.Split("gene LibId 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedResult: true,
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
			name: "with gene name only",
			args: args{
				line: strings.Split("gene 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedTitleString:   strings.Split("0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			expectedGeneTitleIdx:  0,
			expectedLibIdTitleIdx: -1,
		}, {
			name: "with libid only",
			args: args{
				line: strings.Split("LibId 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedTitleString:   strings.Split("0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			expectedGeneTitleIdx:  -1,
			expectedLibIdTitleIdx: 0,
		},
		{
			name: "with geneName and LibId",
			args: args{
				line: strings.Split("gene LibId 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedTitleString:   strings.Split("0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			expectedGeneTitleIdx:  0,
			expectedLibIdTitleIdx: 1,
		},
		{
			name: "with LibId and geneName",
			args: args{
				line: strings.Split("LibId gene 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedTitleString:   strings.Split("0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			expectedGeneTitleIdx:  1,
			expectedLibIdTitleIdx: 0,
		},
		{
			name: "with other titles",
			args: args{
				line: strings.Split("randomTitle 0um-16H 1um-16H 2um-16H 4um-16H 8um-16H", " "),
			},
			expectedTitleString:   nil,
			expectedGeneTitleIdx:  -1,
			expectedLibIdTitleIdx: -1,
			expectedError:         errors.New("input file must have at least one title column"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actualTitleString, actualGeneIdx, actualLibIdIdx, err := processTitleRow(tt.args.line);
				!reflect.DeepEqual(actualTitleString, tt.expectedTitleString) ||
					!reflect.DeepEqual(actualGeneIdx, tt.expectedGeneTitleIdx) ||
					!reflect.DeepEqual(actualLibIdIdx, tt.expectedLibIdTitleIdx) ||
					!reflect.DeepEqual(err, tt.expectedError) {
				t.Errorf("processTitleRow() = \t%v, %v, %v, %v,\n \t\t\t\t\t\t\t\twanted = \t%v, %v, %v, %v",
					actualTitleString, actualGeneIdx, actualLibIdIdx, err,
					tt.expectedTitleString, tt.expectedGeneTitleIdx, tt.expectedLibIdTitleIdx, tt.expectedError)
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
				path: "../../data/testcase1/input.csv",
			},
			want: &ObjectModule.InputDataSheet{
				DataColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
				RowTitles: []ObjectModule.RowTitle{
					{GeneName: "SPBC4F6.10"}, {GeneName: "SPBC4F6.11c"}, {GeneName: "SPBC18A7.01"},
					{GeneName: "SPBC530.04"}, {GeneName: "SPBC557.05"}, {GeneName: "SPBC365.20c"},
					{GeneName: "SPBC56F2.10c"}, {GeneName: "SPBC577.11"}, {GeneName: "SPBC577.14c"},
					{GeneName: "SPBC106.20"},
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
				t.Errorf("readFromCsv() \t= %v, \n" +
					"\t\t\t\t\t\t\twant \t= %v", got, tt.want)
			}
		})
	}
}


func Test_writeToCsv(t *testing.T) {
	type args struct {
		outputDataSheet *ObjectModule.OutputDataSheet
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
