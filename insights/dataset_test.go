package insights

import (
	"testing"
)

/*
	This file contains the tests for the dataset
*/

func TestNewDataset(t *testing.T) {
	dataset := NewDataset()
	if dataset.DataF == nil {
		t.Fatal("Failed NewDataset. Got DataF as nil with NewDataset")
	}
	if dataset.DataS == nil {
		t.Fatal("Failed NewDataset. Got DataS as nil with NewDataset")
	}
	if dataset.Metrics == nil {
		t.Fatal("Failed NewDataset. Got Metrics as nil with NewDataset")
	}
}

type dAddMetricTC struct {
	ID          string
	Description string
	Dataset     Dataset
	Metric      Metric
	Data        interface{}
	Expected    error
}

var dAddMetricTCs = []dAddMetricTC{
	{"1", "Normal case Float", NewDataset(),
		Metric{Name: "Cars", DataType: Float}, []float64{0}, nil},
	{"2", "Normal case String", NewDataset(),
		Metric{Name: "Cars", DataType: String}, []string{"Tesla"}, nil},
	{"3", "Metric data type given as float but provided data is different",
		NewDataset(), Metric{Name: "Cars", DataType: Float}, []float32{0},
		&Error{ErrMDAddMetricFloatMismatch, ErrCDataTypeMismatch}},
	{"4", "Metric data type given as string but provided data is different",
		NewDataset(), Metric{Name: "Cars", DataType: String}, []float32{0},
		&Error{ErrMDAddMetricStringMismatch, ErrCDataTypeMismatch}},
	{"5", "Metric data type given is unsupported",
		NewDataset(), Metric{Name: "Cars", DataType: "notfloat"}, []int{0},
		&Error{ErrMDAddMetricUnsupportedType + "notfloat",
			ErrCUnsupportedDataType}},
	{"6", "Metric data has unequal no. of records for string",
		Dataset{Length: 2, Metrics: map[string]Metric{"Company": {
			Name: "Company"}}},
		Metric{Name: "Cars", DataType: String}, []string{"CX100"},
		&Error{ErrMMetricsDatasizeIncorrect, ErrCMetricSizeMismatch}},
	{"7", "Metric data has unequal no. of records for float",
		Dataset{Length: 2, Metrics: map[string]Metric{"Company": {
			Name: "Company"}}},
		Metric{Name: "Cars", DataType: Float}, []float64{1.1},
		&Error{ErrMMetricsDatasizeIncorrect, ErrCMetricSizeMismatch}},
}

func TestDataset_AddMetric(t *testing.T) {
	for _, v := range dAddMetricTCs {
		t.Run(v.ID, func(t *testing.T) {
			err := v.Dataset.AddMetric(v.Metric, v.Data)
			if err == v.Expected {
				return
			}
			if (err == nil && v.Expected != nil) || (err != nil &&
				v.Expected == nil) || err.Error() != v.Expected.Error() {
				t.Error("Failed", v.ID, "Expected:", v.Expected, "Got:", err)
			}
		})
	}
}

type dCorrelationTC struct {
	ID          string
	Description string
	Dataset     Dataset
	Metrics     []Metric
	Datas       []interface{}
	Var1        string
	Var2        string
	Weight      []float64
	ExpectedC   float64
	ExpectedErr error
}

var dCorrelationTCs = []dCorrelationTC{
	{"1", "Normal case", NewDataset(), []Metric{
		{Name: "age", DataType: Float},
		{Name: "height", DataType: Float},
	}, []interface{}{
		[]float64{10, 20, 30},
		[]float64{140, 178, 190},
	}, "age", "height", []float64{1, 1, 1},
		0.9577677079477441, nil},
	{"2", "Variable doesn't exists", NewDataset(), []Metric{
		{Name: "age", DataType: Float},
	}, []interface{}{
		[]float64{10, 20, 30},
	}, "age", "height", []float64{1, 1, 1},
		0.0, &Error{ErrMDCorrelationNoVaraible, ErrCGeneric}},
	{"3", "Variables have different data types", NewDataset(), []Metric{
		{Name: "age", DataType: Float},
		{Name: "height", DataType: String},
	}, []interface{}{
		[]float64{10, 20, 30},
		[]string{"140", "178", "190"},
	}, "age", "height", []float64{1, 1, 1},
		0.0, &Error{ErrMDCorrelationDatatypeMismatch + "age(" + Float +
			") and height(" + String + ")", ErrCDataTypeMismatch}},
	{"4", "String data types are given for test", NewDataset(), []Metric{
		{Name: "age", DataType: String},
		{Name: "height", DataType: String},
	}, []interface{}{
		[]string{"10", "20", "30"},
		[]string{"140", "178", "190"},
	}, "age", "height", []float64{1, 1, 1},
		0.0, &Error{ErrMDCorrelationNonFloat + String, ErrCUnsupportedDataType}},
	{"5", "Corrupt data", Dataset{
		Metrics: map[string]Metric{"age": {Name: "age", DataType: Float, Index: 2}},
		Length:  3,
	}, []Metric{
		{Name: "height", DataType: Float},
	}, []interface{}{
		[]float64{140, 178, 190},
	}, "age", "height", []float64{1, 1, 1},
		0.0, nil},
}

func TestDataset_Correlation(t *testing.T) {
	for _, v := range dCorrelationTCs {
		t.Run(v.ID, func(t *testing.T) {
			d := v.Dataset
			for i, m := range v.Metrics {
				err := d.AddMetric(m, v.Datas[i])
				if err != nil {
					t.Fatal("Error while adding metric in testcase with ID", v.ID, err)
				}
			}
			crr, errC := d.Correlation(v.Var1, v.Var2, v.Weight)
			if errC == v.ExpectedErr && crr == v.ExpectedC {
				return
			}
			if (errC == nil && v.ExpectedErr != nil) || (errC != nil && v.ExpectedErr == nil) {
				t.Fatal("Failed", v.ID, "Expected Error:", v.ExpectedErr, "Got Error:", errC)
			}
			if errC != nil && v.ExpectedErr != nil && errC.Error() != v.ExpectedErr.Error() {
				t.Fatal("Failed", v.ID, "Expected Error:", v.ExpectedErr, "Got Error:", errC)
			}
			if crr != v.ExpectedC {
				t.Fatal("Failed", v.ID, "Expected Correlation:", v.ExpectedC, "Got Correlation:", crr)
			}
		})
	}
}
