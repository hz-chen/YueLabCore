package BusinessLogic

import (
	"YusLabCore/src/ObjectModule"
	"fmt"
	"math"
	"reflect"
	"testing"
)

const float64EqualityThreshold = 1e-6 // this value can't be too precise, as the excel calculation has different precision.

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
				Data: getInputMatrix(),
			},
		},
		want: 9,
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := Calculate(tt.args.inputDataSheet); got != tt.want {
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
				Data: getInputMatrix(),
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
		if got := Adjust(&tt.args.inputDataSheet, tt.args.baseGeneIndex); !almostEqualsOutputSheet(*got, tt.want) {
			t.Errorf("Adjust() \t= %v, \n" +
				"\t\t\t\t\t\twant\t= %v", *got, tt.want)
		}
	})
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
			arr: getDividingResultMatrix(matrix)[row],
		},
		se:   getSeMatrix(matrix)[row],
		mean: getMeanMatrix(matrix)[row],
	}

	t.Run(tt.name, func(t *testing.T) {
		if gotSe, gotMean := calculateSeAndMean(tt.args.arr);
			!almostEqualsFloat(gotSe, tt.se) || !almostEqualsFloat(gotMean, tt.mean) {
			t.Errorf("calculateSeAndMean() :\t se = %v, mean = %v,\n"+
				"\t\t\t\t want :\t se = %v, mean = %v", gotSe, gotMean, tt.se, tt.mean)
		}
	})

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
				Data:                  getDividingResultMatrix(matrix),
				Se:                    nil,
				Mean:                  nil,
				SeMean:                nil,
			},
		},
		want: ObjectModule.CalculationDataSheet{
			CurrentDividingTarget: matrix,
			Data:                  getDividingResultMatrix(matrix),
			Se:                    getSeMatrix(matrix),
			Mean:                  getMeanMatrix(matrix),
			SeMean:                nil,
		},
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := fillInSeArrayAndMeanArray(tt.args.sheet); !almostEqualsCalculationSheet(got, tt.want) {
			t.Errorf("fillInSeArrayAndMeanArray(%v) =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\t\t\twant =\t %v", tt.args.sheet.CurrentDividingTarget, got, tt.want)
		}
	})

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
				Data:                  getDividingResultMatrix(target),
				Se:                    getSeMatrix(target),
				Mean:                  getMeanMatrix(target),
				SeMean:                nil,
			},
		},
		want: ObjectModule.CalculationDataSheet{
			CurrentDividingTarget: target,
			Data:                  getDividingResultMatrix(target),
			Se:                    getSeMatrix(target),
			Mean:                  getMeanMatrix(target),
			SeMean:                getSeMeanMatrix(target),
		},
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := fillInSeMeanArray(tt.args.sheet); !almostEqualsCalculationSheet(got, tt.want) {
			t.Errorf("fillInSeMeanArray(%v) =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\twant =\t %v", tt.args.sheet.CurrentDividingTarget, got, tt.want)
		}
	})

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
			inputData: getInputMatrix(),
			targetIdx: targetId,
		},
		want: getDividingResultMatrix(targetId),
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := divideAccordingToRow(tt.args.inputData, tt.args.targetIdx); !almostEqualsMatrix(got, tt.want) {
			t.Errorf("divideAccordingToRow on target %v =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\t\t\t\twant =\t %v", targetId, got, tt.want)
		}
	})
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
			inputData: getInputMatrix(),
			targetIdx: targetId,
		},
		want: getMinSeMeanMatrix(targetId),
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := calculateMinSeMeanAccordingToRow(tt.args.inputData, tt.args.targetIdx); !almostEqualsFloat(got, tt.want) {
			t.Errorf("calculateMinSeMeanAccordingToRow on target %v =\t %v,"+
				"\n\t\t\t\t\t\t\t\t\twant =\t %v", targetId, got, tt.want)
		}
	})
}

func almostEqualsOutputSheet(a ObjectModule.OutputDataSheet, b ObjectModule.OutputDataSheet) bool {
	return reflect.DeepEqual(a.ColumnTitles, b.ColumnTitles) &&
		reflect.DeepEqual(a.RowTitles, b.RowTitles) &&
		almostEqualsMatrix(a.Data, b.Data)
}

func almostEqualsCalculationSheet(a ObjectModule.CalculationDataSheet, b ObjectModule.CalculationDataSheet) bool {
	return a.CurrentDividingTarget == b.CurrentDividingTarget &&
		almostEqualsArray(a.Mean, b.Mean) &&
		almostEqualsArray(a.Se, b.Se) &&
		almostEqualsArray(a.SeMean, b.SeMean) &&
		almostEqualsMatrix(a.Data, b.Data)
}

func almostEqualsMatrix(a [][]float64, b [][]float64) bool {
	if len(a) != len(b) {
		return false
	}

	if len(a[0]) != len(b[0]) {
		return false
	}

	equals := true
	for i, _ := range a {
		equals = equals && almostEqualsArray(a[i], b[i])
	}
	return equals
}

func almostEqualsArray(a []float64, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if !almostEqualsFloat(a[i], b[i]) {
			return false
		}
	}
	return true
}

func almostEqualsFloat(a float64, b float64) bool {
	if math.IsNaN(a) {
		return almostEqualsFloat(0.0, b)
	}

	if math.IsNaN(b) {
		return almostEqualsFloat(a, 0.0)
	}

	return math.Abs(a-b) <= float64EqualityThreshold
}

func getInputMatrix() [][]float64 {
	return [][]float64{
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
	}
}

func getSeMatrix(rowNum int) []float64 {
	return [][]float64{
		{0.0, 0.032451081, 0.487753954, 1.601793734, 0.741290993, 0.971019076, 0.851110879, 2.753418229, 5.677179904, 2.845795479},
		{1.699264547, 0.0, 3.565340096, 11.78167645, 5.307777313, 12.61723425, 9.23904757, 8.510340769, 26.18437034, 15.08476052},
		{0.030704828, 0.00388212, 0.0, 0.188732159, 0.10401625, 0.199456261, 0.166734128, 0.223083896, 0.627457964, 0.127975252},
		{0.012795214, 0.001532047, 0.018956462, 0.0, 0.039652153, 0.067744193, 0.042823987, 0.083410156, 0.289174552, 0.078797035},
		{0.057011689, 0.004829121, 0.111124584, 0.338007953, 0.0, 0.390940893, 0.191984687, 0.412801625, 1.654824435, 0.472659209},
		{0.019631613, 0.002797343, 0.039402746, 0.123384288, 0.068522108, 0.0, 0.107176705, 0.20959054, 0.645060599, 0.202751957},
		{0.018117622, 0.003585775, 0.048487755, 0.113486864, 0.062995348, 0.168098231, 0.0, 0.258726788, 0.66185136, 0.256586785},
		{0.011232174, 0.000604345, 0.015257092, 0.045563935, 0.030681814, 0.071657591, 0.048970901, 0.0, 0.162101975, 0.046233057},
		{0.00336747, 0.000291508, 0.005495132, 0.025991625, 0.013996101, 0.026245017, 0.02196804, 0.022056863, 0.0, 0.027734306},
		{0.007578798, 0.00069072, 0.005336863, 0.03056705, 0.01985959, 0.045516759, 0.032687591, 0.028441734, 0.142311729, 0.0},
	}[rowNum]
}

func getSeMeanMatrix(rowNum int) []float64 {
	return [][]float64{
		{0.0, 0.261637909, 0.132008213, 0.143847313, 0.191439749, 0.127559468, 0.126625284, 0.187832753, 0.148428542, 0.156719867},
		{0.174283543, 0.0, 0.107975169, 0.118730993, 0.156802875, 0.179323966, 0.149378295, 0.06806639, 0.078092366, 0.094723771},
		{0.107117789, 0.121567008, 0.0, 0.062769328, 0.101583744, 0.095679675, 0.090194712, 0.057615977, 0.060659658, 0.026342757},
		{0.131732688, 0.141943463, 0.056197606, 0.0, 0.115110093, 0.096922043, 0.069580309, 0.064233103, 0.082835055, 0.048250316},
		{0.189713026, 0.148564724, 0.108767767, 0.110575016, 0.0, 0.18074028, 0.103329608, 0.105006216, 0.15427189, 0.095806587},
		{0.138972601, 0.173513108, 0.079654093, 0.083478785, 0.134150022, 0.0, 0.11701229, 0.108801155, 0.124792062, 0.084298314},
		{0.114781313, 0.196657034, 0.086864342, 0.068523012, 0.111738255, 0.143757145, 0.0, 0.119300395, 0.114063647, 0.094620703},
		{0.146887903, 0.074042516, 0.058283278, 0.058285528, 0.114900306, 0.130359507, 0.101168244, 0.0, 0.060263828, 0.036555109},
		{0.119876822, 0.09486399, 0.056071258, 0.088121889, 0.137674013, 0.127306477, 0.120541473, 0.058498432, 0.0, 0.058285449},
		{0.126845513, 0.105622458, 0.025856281, 0.049444202, 0.094453101, 0.105606965, 0.085793702, 0.035783692, 0.066669249, 0.0},
	}[rowNum]
}

func getMeanMatrix(rowNum int) []float64 {
	return [][]float64{
		{0.0, 0.124030501, 3.694875817, 11.13537473, 3.872189542, 7.612285403, 6.721492375, 14.65888235, 38.24857298, 18.15848584},
		{9.75, 0.0, 33.02, 99.23, 33.85, 70.36, 61.85, 125.03, 335.3, 159.25},
		{0.286645464, 0.031933997, 0.0, 3.006757655, 1.023945827, 2.084625198, 1.848602029, 3.871910306, 10.34390878, 4.858081088},
		{0.09713014, 0.010793361, 0.337317956, 0.0, 0.344471559, 0.698955477, 0.615461295, 1.298554047, 3.490968327, 1.633088489},
		{0.300515415, 0.032505164, 1.021668337, 3.056820291, 0.0, 2.162998162, 1.8579833, 3.931211339, 10.72667507, 4.933472961},
		{0.141262472, 0.016121794, 0.494673205, 1.478031673, 0.510787152, 0.0, 0.915944002, 1.926363179, 5.169083592, 2.405172152},
		{0.1578447, 0.01823365, 0.558200915, 1.656186156, 0.563776013, 1.169320877, 0.0, 2.168700178, 5.802474114, 2.711740414},
		{0.076467661, 0.008162139, 0.261774782, 0.781736689, 0.267029873, 0.549692098, 0.484054074, 0.0, 2.68987186, 1.26474953},
		{0.028091082, 0.003072903, 0.098002644, 0.294950838, 0.101661168, 0.20615618, 0.182244664, 0.377050503, 0.0, 0.47583584},
		{0.059748252, 0.006539515, 0.206404883, 0.618213029, 0.210258741, 0.431001482, 0.381002216, 0.794823889, 2.13459325, 0.0},
	}[rowNum]
}

func getDividingResultMatrix(rowNum int) [][]float64 {
	return [][][]float64{
		{
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.1, 0.078431373, 0.25, 0.117647059, 0.074074074},
			{2.86, 3.117647059, 5.5, 3.941176471, 3.055555556},
			{8.34, 9.333333333, 15.375, 14.64705882, 7.981481481},
			{2.2, 2.431372549, 6.125, 4.882352941, 3.722222222},
			{4.98, 8.862745098, 10.125, 8.352941176, 5.740740741},
			{5.3, 4.980392157, 7.875, 9.470588235, 5.981481481},
			{11.14, 10.88235294, 24.625, 16.64705882, 10},
			{35.9, 30.1372549, 59.625, 38.11764706, 27.46296296},
			{13.4, 14.56862745, 28.375, 20.41176471, 14.03703704},
		},
		{
			{10, 12.75, 4, 8.5, 13.5},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{28.6, 39.75, 22, 33.5, 41.25},
			{83.4, 119, 61.5, 124.5, 107.75},
			{22, 31, 24.5, 41.5, 50.25},
			{49.8, 113, 40.5, 71, 77.5},
			{53, 63.5, 31.5, 80.5, 80.75},
			{111.4, 138.75, 98.5, 141.5, 135},
			{359, 384.25, 238.5, 324, 370.75},
			{134, 185.75, 113.5, 173.5, 189.5},
		},
		{
			{0.34965035, 0.320754717, 0.181818182, 0.253731343, 0.327272727},
			{0.034965035, 0.025157233, 0.045454545, 0.029850746, 0.024242424},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{2.916083916, 2.993710692, 2.795454545, 3.71641791, 2.612121212},
			{0.769230769, 0.779874214, 1.113636364, 1.23880597, 1.218181818},
			{1.741258741, 2.842767296, 1.840909091, 2.119402985, 1.878787879},
			{1.853146853, 1.597484277, 1.431818182, 2.402985075, 1.957575758},
			{3.895104895, 3.490566038, 4.477272727, 4.223880597, 3.272727273},
			{12.55244755, 9.666666667, 10.84090909, 9.671641791, 8.987878788},
			{4.685314685, 4.672955975, 5.159090909, 5.179104478, 4.593939394},
		},
		{
			{0.119904077, 0.107142857, 0.06504065, 0.068273092, 0.125290023},
			{0.011990408, 0.008403361, 0.016260163, 0.008032129, 0.009280742},
			{0.342925659, 0.334033613, 0.357723577, 0.269076305, 0.382830626},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.263788969, 0.260504202, 0.398373984, 0.333333333, 0.466357309},
			{0.597122302, 0.949579832, 0.658536585, 0.570281124, 0.719257541},
			{0.635491607, 0.533613445, 0.512195122, 0.646586345, 0.749419954},
			{1.335731415, 1.165966387, 1.601626016, 1.136546185, 1.252900232},
			{4.304556355, 3.228991597, 3.87804878, 2.602409639, 3.440835267},
			{1.606714628, 1.56092437, 1.845528455, 1.393574297, 1.758700696},
		},
		{
			{0.454545455, 0.411290323, 0.163265306, 0.204819277, 0.268656716},
			{0.045454545, 0.032258065, 0.040816327, 0.024096386, 0.019900498},
			{1.3, 1.282258065, 0.897959184, 0.807228916, 0.820895522},
			{3.790909091, 3.838709677, 2.510204082, 3, 2.144278607},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{2.263636364, 3.64516129, 1.653061224, 1.710843373, 1.542288557},
			{2.409090909, 2.048387097, 1.285714286, 1.939759036, 1.606965174},
			{5.063636364, 4.475806452, 4.020408163, 3.409638554, 2.686567164},
			{16.31818182, 12.39516129, 9.734693878, 7.807228916, 7.378109453},
			{6.090909091, 5.991935484, 4.632653061, 4.180722892, 3.771144279},
		},
		{
			{0.200803213, 0.112831858, 0.098765432, 0.11971831, 0.174193548},
			{0.020080321, 0.008849558, 0.024691358, 0.014084507, 0.012903226},
			{0.574297189, 0.351769912, 0.543209877, 0.471830986, 0.532258065},
			{1.674698795, 1.053097345, 1.518518519, 1.753521127, 1.390322581},
			{0.441767068, 0.274336283, 0.604938272, 0.584507042, 0.648387097},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{1.064257028, 0.561946903, 0.777777778, 1.133802817, 1.041935484},
			{2.236947791, 1.227876106, 2.432098765, 1.992957746, 1.741935484},
			{7.208835341, 3.400442478, 5.888888889, 4.563380282, 4.783870968},
			{2.690763052, 1.64380531, 2.802469136, 2.443661972, 2.44516129},
		},
		{
			{0.188679245, 0.200787402, 0.126984127, 0.105590062, 0.167182663},
			{0.018867925, 0.015748031, 0.031746032, 0.01242236, 0.012383901},
			{0.539622642, 0.625984252, 0.698412698, 0.416149068, 0.510835913},
			{1.573584906, 1.874015748, 1.952380952, 1.546583851, 1.334365325},
			{0.41509434, 0.488188976, 0.777777778, 0.51552795, 0.622291022},
			{0.939622642, 1.779527559, 1.285714286, 0.881987578, 0.959752322},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{2.101886792, 2.18503937, 3.126984127, 1.757763975, 1.671826625},
			{6.773584906, 6.051181102, 7.571428571, 4.02484472, 4.591331269},
			{2.528301887, 2.92519685, 3.603174603, 2.155279503, 2.346749226},
		},
		{
			{0.089766607, 0.091891892, 0.040609137, 0.060070671, 0.1},
			{0.008976661, 0.007207207, 0.010152284, 0.007067138, 0.007407407},
			{0.256732496, 0.286486486, 0.223350254, 0.236749117, 0.305555556},
			{0.748653501, 0.857657658, 0.624365482, 0.879858657, 0.798148148},
			{0.197486535, 0.223423423, 0.248730964, 0.293286219, 0.372222222},
			{0.447037702, 0.814414414, 0.411167513, 0.501766784, 0.574074074},
			{0.475763016, 0.457657658, 0.319796954, 0.568904594, 0.598148148},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{3.222621185, 2.769369369, 2.421319797, 2.28975265, 2.746296296},
			{1.202872531, 1.338738739, 1.152284264, 1.22614841, 1.403703704},
		},
		{
			{0.027855153, 0.033181522, 0.016771488, 0.026234568, 0.036412677},
			{0.002785515, 0.002602472, 0.004192872, 0.00308642, 0.002697235},
			{0.079665738, 0.103448276, 0.092243187, 0.103395062, 0.111260958},
			{0.232311978, 0.309694209, 0.257861635, 0.384259259, 0.290627107},
			{0.061281337, 0.080676643, 0.102725367, 0.12808642, 0.135536076},
			{0.138718663, 0.294079375, 0.169811321, 0.219135802, 0.209035738},
			{0.147632312, 0.165256994, 0.132075472, 0.24845679, 0.217801753},
			{0.310306407, 0.361093038, 0.412997904, 0.436728395, 0.36412677},
			{0.0, 0.0, 0.0, 0.0, 0.0},
			{0.373259053, 0.483409239, 0.475890985, 0.535493827, 0.511126096},
		},
		{

			{0.074626866, 0.068640646, 0.035242291, 0.048991354, 0.071240106},
			{0.007462687, 0.00538358, 0.008810573, 0.005763689, 0.005277045},
			{0.213432836, 0.213997308, 0.193832599, 0.193083573, 0.2176781},
			{0.62238806, 0.64064603, 0.54185022, 0.717579251, 0.568601583},
			{0.164179104, 0.166890983, 0.215859031, 0.239193084, 0.265171504},
			{0.371641791, 0.608344549, 0.356828194, 0.409221902, 0.408970976},
			{0.395522388, 0.341857335, 0.27753304, 0.463976945, 0.426121372},
			{0.831343284, 0.746971736, 0.86784141, 0.81556196, 0.712401055},
			{2.679104478, 2.068640646, 2.101321586, 1.867435159, 1.95646438},
			{0.0, 0.0, 0.0, 0.0, 0.0},
		},
	}[rowNum]
}

func getMinSeMeanMatrix(rowNum int) float64 {
	return []float64{
		0.126625284, 0.06806639, 0.026342757, 0.048250316, 0.095806587, 0.079654093, 0.068523012, 0.036555109, 0.056071258, 0.025856281,
	}[rowNum]
}
