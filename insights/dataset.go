package insights

import (
	"log"

	"github.com/gonum/stat"
)

/*
	This file contains the structs and other data models required for the
	datasets.
*/

const (
	//Float is used to denote the variables with data type float
	Float = "float64"
	//String is used to denote the variables with data type string
	String = "string"
)

//Dataset stores a data in columnar form
type Dataset struct {
	//DataF contains the metrics in the dataset that are of data type Float
	DataF [][]float64
	//DataS contains the metrics in the dataset that are of data type String
	DataS [][]string
	//Metrics has the map of metrics in a data set mapped to their names
	Metrics map[string]Metric
	//Length is the no of records in the dataset. It is set after the first
	//metric is set.
	Length int64
}

//Metric is holds information about a metric in a Dataset
type Metric struct {
	Name string //Name is the name of the metric.
	//Index is the index at which metric is stored in the dataset.
	Index    int
	DataType string //DataType is the data type of the metric
	//DisplayName is a friendly name to be used for displaying the metric name.
	DisplayName string
}

//NewDataset returns an initialized Dataset.
//The data arrays are initialized in the Dataset that is returned.
func NewDataset() Dataset {
	return Dataset{[][]float64{}, [][]string{}, map[string]Metric{}, 0}
}

//AddMetric adds a metric to the dataset.
//We need to pass metric and the data corresponding to the metric.
//Depending upon the data type of the metric the data will added the data arrays
//of the dataset. It will return an error if the data type inferred from the
//metric is different from the one that is actually passed to the method.
//It will also return an error if the data type of the passed metric
//is not supported. If the metric is added successfully,
//the method will return nil.
func (d *Dataset) AddMetric(m Metric, data interface{}) error {
	/*
		First we will try to infer the data type from the metric.
		Then will check whether the no. of records provided is inmatch
		//with the dataset length.
		Then add it to the data arrays of the dataset.
		At the end we will check if the metric set is the first one,
		we set the Length property to the length of the data.
		If everything goes right we will return nil.
	*/
	//Based on the data type we will add the metric
	switch m.DataType {
	case Float:
		//We will try to do a type assertion for float
		df, ok := data.([]float64)
		if !ok {
			//The given array is not float
			return &Error{ErrMDAddMetricFloatMismatch, ErrCDataTypeMismatch}
		}

		//Checking the length of the metric
		if len(d.Metrics) != 0 && d.Length != int64(len(df)) {
			//The given metric has incorrect no. of records
			return &Error{ErrMMetricsDatasizeIncorrect, ErrCMetricSizeMismatch}
		}

		//Adding the metric to the dataset
		d.DataF = append(d.DataF, df)
		m.Index = len(d.DataF) - 1
		d.Metrics[m.Name] = m

		//Checking whether the metric was the first one
		if len(d.Metrics) == 1 {
			d.Length = int64(len(df))
		}
		break
	case String:
		//We will try to do a type assertion for string
		ds, ok := data.([]string)
		if !ok {
			//The given array is not string
			return &Error{ErrMDAddMetricStringMismatch, ErrCDataTypeMismatch}
		}
		//Checking the length of the metric
		if len(d.Metrics) != 0 && d.Length != int64(len(ds)) {
			//The given metric has incorrect no. of records
			return &Error{ErrMMetricsDatasizeIncorrect, ErrCMetricSizeMismatch}
		}

		//Adding the metric to the dataset
		d.DataS = append(d.DataS, ds)
		m.Index = len(d.DataS) - 1
		d.Metrics[m.Name] = m

		//Checking whether the metric was the first one
		if len(d.Metrics) == 1 {
			d.Length = int64(len(ds))
		}
		break
	default:
		return &Error{ErrMDAddMetricUnsupportedType + m.DataType,
			ErrCUnsupportedDataType}
	}

	//Everything went right
	return nil
}

//Correlation finds the correlation between two variables in the dataset.
//For finding the correlation between two variables, they must have same data
// types and their data type must be Float. In these cases correlation will
// be zero and an error will be returned
func (d Dataset) Correlation(var1, var2 string, weights []float64) (
	float64, error) {
	/*
		First we will check whether the variables exists in the dataset
		If the variables doesn't exist in the dataset,
		or their data types are mismatch or their data type isn't Float
		we will simply return 0.0 and an error
		Else we will find the correlation and return them.
	*/
	//Checking whether the variabls exist in the dataset
	m1, ok1 := d.Metrics[var1]
	m2, ok2 := d.Metrics[var2]

	//If the variables doesn't exist we will return 0.0
	if !ok1 || !ok2 {
		return float64(0.0), &Error{ErrMDCorrelationNoVaraible, ErrCGeneric}
	}
	//If the data types doesn't match
	if m1.DataType != m2.DataType {
		return float64(0.0), &Error{ErrMDCorrelationDatatypeMismatch + m1.Name +
			"(" + m1.DataType + ") and " + m2.Name + "(" +
			m2.DataType + ")", ErrCDataTypeMismatch}
	}
	//If the data types aren't Float
	if m1.DataType != Float {
		return float64(0.0), &Error{ErrMDCorrelationNonFloat + m1.DataType,
			ErrCUnsupportedDataType}
	}

	//Now we return the correlation between the data
	defer func() {
		//We are adding a panic recover as the Correlation function can panic
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	return stat.Correlation(d.DataF[m1.Index], d.DataF[m2.Index], weights), nil
}
