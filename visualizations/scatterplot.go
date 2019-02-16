package visualizations

/*
	This file has the struct and utlities required for the scatter
	plot visualization
*/

//ScatterPlot is the scatter plot visualization
//It is used to plot two continuous/ discrete variables.
//It is often used to visualize the correlation between two variables.
type ScatterPlot struct {
	//M stores the metrics involved in rendering a scatterplot
	M []Metric `json:"Metrics"`
	//T is the title of the scatter plot
	T string `json:"Title"`
	//D is the description of the scatter plot
	D string `json:"Description"`
	//Dt stores the data to be plotted in the scatter plot
	Dt []map[string]interface{} `json:"Data"`
}

//Type returns the scatter plot's type string
func (s ScatterPlot) Type() string {
	return SCATTERPLOT
}

//Metrics returns the metrics involved for creating the scatter plot
func (s ScatterPlot) Metrics() []Metric {
	return s.M
}

//Title returns the title of the scatter plot
func (s ScatterPlot) Title() string {
	return s.T
}

//Description returns the description for the scatter plot
func (s ScatterPlot) Description() string {
	return s.D
}

//Data returns the data to be plotted in th escatter plot visualization
func (s ScatterPlot) Data() []map[string]interface{} {
	return s.Dt
}
