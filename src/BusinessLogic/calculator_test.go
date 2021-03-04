package BusinessLogic

import (
	"YusLabCore/src/ObjectModule"
	"YusLabCore/src/TestUtils"
	"fmt"
	"testing"
)

func BenchmarkCalculate(b *testing.B) {
	input := ObjectModule.InputDataSheet{
		DataColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
		RowTitles: []ObjectModule.RowTitle{
			{GeneName: "SPBC4F6.10"}, {GeneName: "SPBC4F6.11c"}, {GeneName: "SPBC18A7.01"},
			{GeneName: "SPBC530.04"}, {GeneName: "SPBC557.05"}, {GeneName: "SPBC365.20c"},
			{GeneName: "SPBC56F2.10c"}, {GeneName: "SPBC577.11"}, {GeneName: "SPBC577.14c"},
			{GeneName: "SPBC106.20"},
		},
		Data: TestUtils.GetInputMatrix(),
	}
	for i := 0; i < b.N; i++ {
		_ = Calculate(input)
	}
}

func TestCalculate(t *testing.T) {
	type args struct {
		inputDataSheet ObjectModule.InputDataSheet
	}
	tt := struct {
		name string
		args args
		want int
	}{
		name: "based on only input data",
		args: args{
			inputDataSheet: ObjectModule.InputDataSheet{
				DataColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
				RowTitles: []ObjectModule.RowTitle{
					{GeneName: "SPBC4F6.10"}, {GeneName: "SPBC4F6.11c"}, {GeneName: "SPBC18A7.01"},
					{GeneName: "SPBC530.04"}, {GeneName: "SPBC557.05"}, {GeneName: "SPBC365.20c"},
					{GeneName: "SPBC56F2.10c"}, {GeneName: "SPBC577.11"}, {GeneName: "SPBC577.14c"},
					{GeneName: "SPBC106.20"},
				},
				Data: TestUtils.GetInputMatrix(),
			},
		},
		want: 9,
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := Calculate(tt.args.inputDataSheet); got.Index != tt.want {
			t.Errorf("Calculate() = %v, want %v", got, tt.want)
		}
	})

}

func TestAdjust(t *testing.T) {
	type args struct {
		inputDataSheet ObjectModule.InputDataSheet
		baseGeneIndex  int
	}
	tt := struct {
		name string
		args args
		want ObjectModule.OutputDataSheet
	}{
		name: "based on only input data",
		args: args{
			inputDataSheet: ObjectModule.InputDataSheet{
				DataColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
				RowTitles: []ObjectModule.RowTitle{
					{GeneName: "SPBC4F6.10"}, {GeneName: "SPBC4F6.11c"}, {GeneName: "SPBC18A7.01"},
					{GeneName: "SPBC530.04"}, {GeneName: "SPBC557.05"}, {GeneName: "SPBC365.20c"},
					{GeneName: "SPBC56F2.10c"}, {GeneName: "SPBC577.11"}, {GeneName: "SPBC577.14c"},
					{GeneName: "SPBC106.20"},
				},
				Data: TestUtils.GetInputMatrix(),
			},
			baseGeneIndex: 9,
		},
		want: ObjectModule.OutputDataSheet{
			ColumnTitles: []string{"0um-16H", "1um-16H", "2um-16H", "4um-16H", "8um-16H"},
			RowTitles: []ObjectModule.RowTitle{
				{GeneName: "SPBC4F6.10"}, {GeneName: "SPBC4F6.11c"}, {GeneName: "SPBC18A7.01"},
				{GeneName: "SPBC530.04"}, {GeneName: "SPBC557.05"}, {GeneName: "SPBC365.20c"},
				{GeneName: "SPBC56F2.10c"}, {GeneName: "SPBC577.11"}, {GeneName: "SPBC577.14c"},
				{GeneName: "SPBC106.20"},
			},
			Data: [][]float64{
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
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := Adjust(&tt.args.inputDataSheet, tt.args.baseGeneIndex); !TestUtils.AlmostEqualsOutputSheet(*got, tt.want) {
			t.Errorf("Adjust() \t= %v, \n"+
				"\t\t\t\t\t\twant\t= %v", *got, tt.want)
		}
	})
}

func BenchmarkCalculateSeAndMean(b *testing.B) {
	input := TestUtils.GetDividingResultMatrix(0)[0]
	for i := 0; i < b.N; i++ {
		_, _ = calculateSeAndMean(input)
	}
}

func TestCalculateSeAndMean(t *testing.T) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			innerTestCalculateSeAndMean(t, i, j)
		}
	}
}

func innerTestCalculateSeAndMean(t *testing.T, matrix int, row int) {
	type args struct {
		arr []float64
	}
	tt := struct {
		name string
		args args
		se   float64
		mean float64
	}{
		name: fmt.Sprintf("matrix %v row %v", matrix, row),
		args: args{
			arr: TestUtils.GetDividingResultMatrix(matrix)[row],
		},
		se:   TestUtils.GetSeMatrix(matrix)[row],
		mean: TestUtils.GetMeanMatrix(matrix)[row],
	}

	t.Run(tt.name, func(t *testing.T) {
		if gotSe, gotMean := calculateSeAndMean(tt.args.arr);
			!TestUtils.AlmostEqualsFloat(gotSe, tt.se) || !TestUtils.AlmostEqualsFloat(gotMean, tt.mean) {
			t.Errorf("calculateSeAndMean() :\t se = %v, mean = %v,\n"+
				"\t\t\t\t want :\t se = %v, mean = %v", gotSe, gotMean, tt.se, tt.mean)
		}
	})

}

func BenchmarkFillInSeArrayAndMeanArray(b *testing.B) {
	input := ObjectModule.CalculationDataSheet{
		CurrentDividingTarget: 0,
		Data:                  TestUtils.GetDividingResultMatrix(0),
		Se:                    nil,
		Mean:                  nil,
		SeMean:                nil,
	}
	for i := 0; i < b.N; i++ {
		_ = fillInSeArrayAndMeanArray(input)
	}
}

func TestFillInSeArrayAndMeanArray(t *testing.T) {
	for i := 0; i < 10; i++ {
		innerTestFillInSeArrayAndMeanArrayForMatrix(t, i)
	}
}

func innerTestFillInSeArrayAndMeanArrayForMatrix(t *testing.T, matrix int) {
	type args struct {
		sheet ObjectModule.CalculationDataSheet
	}
	tt := struct {
		name string
		args args
		want ObjectModule.CalculationDataSheet
	}{
		name: fmt.Sprintf("matrix %v", matrix),
		args: args{
			sheet: ObjectModule.CalculationDataSheet{
				CurrentDividingTarget: matrix,
				Data:                  TestUtils.GetDividingResultMatrix(matrix),
				Se:                    nil,
				Mean:                  nil,
				SeMean:                nil,
			},
		},
		want: ObjectModule.CalculationDataSheet{
			CurrentDividingTarget: matrix,
			Data:                  TestUtils.GetDividingResultMatrix(matrix),
			Se:                    TestUtils.GetSeMatrix(matrix),
			Mean:                  TestUtils.GetMeanMatrix(matrix),
			SeMean:                nil,
		},
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := fillInSeArrayAndMeanArray(tt.args.sheet); !TestUtils.AlmostEqualsCalculationSheet(got, tt.want) {
			t.Errorf("fillInSeArrayAndMeanArray(%v) =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\t\t\twant =\t %v", tt.args.sheet.CurrentDividingTarget, got, tt.want)
		}
	})

}

func BenchmarkFillInSeMeanArray(b *testing.B) {
	input := ObjectModule.CalculationDataSheet{
		CurrentDividingTarget: 0,
		Data:                  TestUtils.GetDividingResultMatrix(0),
		Se:                    TestUtils.GetSeMatrix(0),
		Mean:                  TestUtils.GetMeanMatrix(0),
		SeMean:                nil,
	}

	for i := 0; i < b.N; i++ {
		_ = fillInSeMeanArray(input)
	}
}

func TestFillInSeMeanArray(t *testing.T) {
	for i := 0; i < 10; i++ {
		innerTestFillInSeMeanArray(t, i)
	}
}

func innerTestFillInSeMeanArray(t *testing.T, target int) {
	type args struct {
		sheet ObjectModule.CalculationDataSheet
	}
	tt := struct {
		name string
		args args
		want ObjectModule.CalculationDataSheet
	}{
		name: fmt.Sprintf("target %v", target),
		args: args{
			sheet: ObjectModule.CalculationDataSheet{
				CurrentDividingTarget: target,
				Data:                  TestUtils.GetDividingResultMatrix(target),
				Se:                    TestUtils.GetSeMatrix(target),
				Mean:                  TestUtils.GetMeanMatrix(target),
				SeMean:                nil,
			},
		},
		want: ObjectModule.CalculationDataSheet{
			CurrentDividingTarget: target,
			Data:                  TestUtils.GetDividingResultMatrix(target),
			Se:                    TestUtils.GetSeMatrix(target),
			Mean:                  TestUtils.GetMeanMatrix(target),
			SeMean:                TestUtils.GetSeMeanMatrix(target),
		},
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := fillInSeMeanArray(tt.args.sheet); !TestUtils.AlmostEqualsCalculationSheet(got, tt.want) {
			t.Errorf("fillInSeMeanArray(%v) =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\twant =\t %v", tt.args.sheet.CurrentDividingTarget, got, tt.want)
		}
	})

}

func BenchmarkDivideAccordingToRow(b *testing.B) {
	input := TestUtils.GetInputMatrix()
	for i := 0; i < b.N; i++ {
		_ = divideAccordingToRow(input, 0)
	}
}

func TestDivideAccordingToRow(t *testing.T) {
	for i := 0; i < 10; i++ {
		innerTestDivideAccordingToRow(t, i)
	}
}

func innerTestDivideAccordingToRow(t *testing.T, targetId int) {
	type args struct {
		inputData [][]float64
		targetIdx int
	}
	tt := struct {
		name string
		args args
		want [][]float64
	}{
		name: fmt.Sprintf("target %v", targetId),
		args: args{
			inputData: TestUtils.GetInputMatrix(),
			targetIdx: targetId,
		},
		want: TestUtils.GetDividingResultMatrix(targetId),
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := divideAccordingToRow(tt.args.inputData, tt.args.targetIdx); !TestUtils.AlmostEqualsMatrix(got, tt.want) {
			t.Errorf("divideAccordingToRow on target %v =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\t\t\t\twant =\t %v", targetId, got, tt.want)
		}
	})
}

func BenchmarkCalculateMinSeMeanAccordingToRow(b *testing.B) {
	input := TestUtils.GetInputMatrix()
	for i := 0; i < b.N; i++ {
		_ = calculateMinSeMeanAccordingToRow(input, 0)
	}
}

func TestCalculateMinSeMeanAccordingToRow(t *testing.T) {
	for i := 0; i < 10; i++ {
		innerTestCalculateMinSeMeanAccordingToRow(t, i)
	}
}

func innerTestCalculateMinSeMeanAccordingToRow(t *testing.T, targetId int) {
	type args struct {
		inputData [][]float64
		targetIdx int
	}
	tt := struct {
		name string
		args args
		want float64
	}{
		name: fmt.Sprintf("target %v", targetId),
		args: args{
			inputData: TestUtils.GetInputMatrix(),
			targetIdx: targetId,
		},
		want: TestUtils.GetMinSeMeanMatrix(targetId),
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := calculateMinSeMeanAccordingToRow(tt.args.inputData, tt.args.targetIdx); !TestUtils.AlmostEqualsFloat(got, tt.want) {
			t.Errorf("calculateMinSeMeanAccordingToRow on target %v =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\twant =\t %v", targetId, got, tt.want)
		}
	})
}
