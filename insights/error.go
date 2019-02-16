package insights

import "fmt"

/*
	This file contains the structs for the error handling in insights package
*/

const (
	//ErrCGeneric inidcates that the error is generic in nature and need to
	//rely on error message for the same
	ErrCGeneric = 0
	//ErrCDataTypeMismatch indicates that the error is due to the data
	//type mismatch
	ErrCDataTypeMismatch = 1
	//ErrCUnsupportedDataType indicates that the error is due to the unsupported
	//data type
	ErrCUnsupportedDataType = 2
	//ErrCMetricSizeMismatch indicates that the size of the metric provided is
	//mistmatch with that of the existing
	ErrCMetricSizeMismatch = 3
)

const (
	//ErrMDAddMetricFloatMismatch is the error message given by add metric
	//method of the dataset
	//when trying to add a float data metric.
	//But provided data has different data type
	ErrMDAddMetricFloatMismatch = "The given data is not of the type []" +
		Float + " while the metric data type is" + Float
	//ErrMDAddMetricStringMismatch is the error message given by add metric
	//method of the dataset
	//when trying to add a string data metric.
	//But provided data has different data type
	ErrMDAddMetricStringMismatch = "The given data is not of the type []" +
		String + " while the metric data type is" + String
	//ErrMDAddMetricUnsupportedType is the error message given by add metric of
	//the dataset when trying to add a metric of unsupported datatype.
	ErrMDAddMetricUnsupportedType = "Unsupported datatype. Got "
	//ErrMDCorrelationNoVaraible is the error message given by correlation
	//function when the given variabes doesn't exist in the dataset.
	ErrMDCorrelationNoVaraible = "Variables doesn't exist in the dataset"
	//ErrMDCorrelationDatatypeMismatch is the error message given when
	//the given variables has data type mismatch
	ErrMDCorrelationDatatypeMismatch = "Datatype mismatch Got "
	//ErrMDCorrelationNonFloat is the error message given the variables
	//doesn't have float fdata type.
	ErrMDCorrelationNonFloat = "Only " + Float + " datatype supported. Got"
	//ErrMMetricsDatasizeIncorrect is the error message informing the no. of
	//records in the /metric is != to that Length property of the dataset
	ErrMMetricsDatasizeIncorrect = "The no. of records provided in the " +
		"metric mismatch to that of the dataset"
)

//Error will be used to return errors in the insights package functions
type Error struct {
	Message string //Message is the error message
	//Code is the error code. This useful while parsing the logs or referring to
	//the error in forums
	Code int
}

//String returns the string form of the error
func (e *Error) String() string {
	return fmt.Sprintf("C-%d %s", e.Code, e.Message)
}

//Error returns the string form of the error. This implmentation qualifies
//the error struct to pass it as a error interface
func (e *Error) Error() string {
	return e.String()
}
