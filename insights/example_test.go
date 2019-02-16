package insights_test

import (
	"fmt"

	"github.com/cuttle-ai/brain/insights"
)

/*
 This file contains the examples for the functions in the insights package
*/

func ExampleDataset_AddMetric() {
	d := insights.NewDataset()
	m := insights.Metric{Name: "age", DataType: insights.Float}
	data := []float64{10, 20, 90}
	err := d.AddMetric(m, data)
	fmt.Println(err)
	// Output:
	// <nil>
}
