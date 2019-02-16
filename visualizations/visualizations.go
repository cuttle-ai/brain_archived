package visualizations

/*
	This file contains the structs and other utilities required for
	creating visualizations.
*/

const (
	//SCATTERPLOT is the string storing the name type of the
	//scatter plot visualization.
	SCATTERPLOT = "SCATTERPLOT"
)

//Visual is the interface to be implemented by any visualization
type Visual interface {
	Type() string //Type is th type of visual
	//Metrics is the metrics involved in the visual
	Metrics() []Metric
	Title() string                  //Title is the title of the visual
	Description() string            //Description for the visual
	Data() []map[string]interface{} //Data to be visualized
}

//Metric has the basic information about a variable in the dataset
//that can be used for rendering a visual
type Metric struct {
	Name string //Name of the metric
	//DisplayName is the name to be shown for the metric
	//Often the real variable will have weird names. For representation purposes
	//it is better to show a different name.
	DisplayName string
	//DataType is the data type of the variable
	DataType string
	//Dimension is the dimension that can be used to identify how to use the
	//metric in the visualization. For example for a bar graph, metric appearing
	//on the x axis should have dimension 0 and y axis should have dimension 1
	Dimension int
	//PreplacmentUnit is the unit to be used as a pre for displaying the values
	//of the metric. Example for showing price $70, $ is the preplacement unit
	PreplacementUnit string
	//PostplacementUnit is the unit to be used as a post for for displaying the
	//values of the metric. Example for showing temprature as 30C, C is the
	//post placement unit.
	PostplacementUnit string
}
