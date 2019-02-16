package insights

import (
	"log"
	"strconv"

	"github.com/cuttle-ai/brain/visualizations"
)

/*
	This file contains the utilities and structs required for correlation
	insights
*/

//Correlation is the correlation insight.
//It states whether to variables are correlated or not
type Correlation struct {
	//visual has the visualization to be used for showing
	//the correlation betweeen metrics in a dataset. By default
	//scatter plot is preferred.
	visual visualizations.Visual
	//relevant stores the information whether the insight is relevant or not.
	//This property is updated after running methods like FSFA and Generate
	relevant bool
	dt       Dataset //dt is the dataset to be used for the corelation
	//ms is the list of metrics on which correlation has to be found
	ms []Metric
}

//New returns a new instance of the Correlation with
//initializations done for the given dataset
func (c *Correlation) New(d Dataset, ms []Metric) Insight {
	return &Correlation{dt: d, ms: ms}
}

//Visual returns the visualization to be used for visualizing the
//correlation between two variables of dataset
func (c *Correlation) Visual() visualizations.Visual {
	return c.visual
}

//Type returns the type string for th correlation type of insight
func (c *Correlation) Type() string {
	return CORRELATION
}

//Relevant returns whether the insight is relevant or not for the given dataset.
func (c *Correlation) Relevant() bool {
	return c.relevant
}

//FSFA does the fast statistical feasibilty analysis over the dataset
//with the given metrics whether the correlation is statistically
//possible between the metrics. Note this function is still under development.
//Not ready to use.
//Plese update this documentation when FSFA is production ready.
func (c *Correlation) FSFA() {
	/*
		Will check whether the length of the metrics array is
		2. Can check correlation between only two variables.
		Then it will check whether the data types of the variables
		are float.
	*/
	//Checking the length between the metrics
	if len(c.ms) != 2 {
		c.relevant = false
		return
	}

	//checking the data types of the metrics
	if c.ms[0].DataType != Float || c.ms[1].DataType != Float {
		c.relevant = false
		return
	}

	//Everything is fine
	c.relevant = true
}

//Generate generates the correlation insight for the datatset associated with
//it for the provided variables.
//This method can only be run after running the FSFA.
//Else the insight won't be generated
func (c *Correlation) Generate() {
	/*
		If the correlation is not relevant we won't event bother
		to go forward.
		We won't go forward if the length of metric < 2 or if the length of
		dataset's float array < the metric indices
		Then we will create a weights array with weight = 1.
		Then will run the correlation on the dataset with the
		given variables.
		Now we will create the visualization for the corelation.
		Then add data to the visualization data.
	*/
	//Checking whether the existing relevance of the insight
	if !c.relevant {
		return
	}

	//Checking whether there are suffcient metrics and dataset arrays
	if len(c.ms) < 2 || len(c.dt.DataF) < c.ms[0].Index ||
		len(c.dt.DataF) < c.ms[1].Index {
		c.relevant = false
		return
	}

	//Creating the weights array
	weights := make([]float64, len(c.dt.DataF[c.ms[0].Index]))
	for i := 0; i < len(weights); i++ {
		weights[i] = float64(1)
	}

	//Running the correlation
	corr, err := c.dt.Correlation(c.ms[0].Name, c.ms[1].Name, weights)
	if err != nil {
		//Error while generating the correlation between the variables
		log.Println(err)
		c.relevant = false
		return
	}
	if corr < 0.7 {
		//Don't bother to look for the correlation
		c.relevant = false
		return
	}

	//Now we have a correlation.
	//Will create the visual for the same.
	c.relevant = true
	visual := visualizations.ScatterPlot{
		T: c.ms[0].DisplayName + " and " + c.ms[1].DisplayName,
		D: "have a correlation of " + strconv.FormatFloat(corr, 'f', -1, 64),
		M: []visualizations.Metric{
			{
				Name:        c.ms[0].Name,
				DisplayName: c.ms[0].DisplayName,
				DataType:    Float,
				Dimension:   0,
			},
			{
				Name:        c.ms[1].Name,
				DisplayName: c.ms[1].DisplayName,
				DataType:    Float,
				Dimension:   1,
			},
		},
	}

	data := make([]map[string]interface{}, len(c.dt.DataF[c.ms[0].Index]))

	//Now we need to add data to the records
	//Will iterate through each metric's data array and add them to the
	//visualization data. First we will add the data from the first metric.
	//Then from the second metric.
	for i, v := range c.dt.DataF[c.ms[0].Index] {
		data[i] = map[string]interface{}{}
		data[i][c.ms[0].Name] = v
	}
	for i, v := range c.dt.DataF[c.ms[1].Index] {
		data[i][c.ms[1].Name] = v
	}
	visual.Dt = data
	c.visual = visual
}

//Propose suggests the possible insights from the domain knowledge.
//This is a WIP. Please dont use it in production
func (c *Correlation) Propose(d Dataset) []ProposedInsight {
	/*
		We will iterate through the metrics in the data set and select the
		variables that have float data type.
		We will get the combination of all the selected metrics as the group of
		two. Then add it as proposed insight.
	*/
	//variable for storing the result
	result := []ProposedInsight{}
	//variable to store the float variables
	svars := []Metric{}

	//iterating through the dataset metrics to select the ones with float
	//datatype.
	for _, v := range d.Metrics {
		if v.DataType != Float {
			//We need to select only the float metrics
			continue
		}
		svars = append(svars, v)
	}

	//iterating through the variables to create the proposals with the
	//combination of all the float variables
	for i := 0; i < len(svars)-1; i++ {
		//We need to iterate from zeroth element to the penaltimate element
		//Then again iterate to the elements starting from i+1 th element to the
		//last element in the selected metric array.
		for j := i + 1; j < len(svars); j++ {
			//Now we have a combination we may now propose a insight
			metrics := []Metric{svars[i], svars[j]}
			result = append(result, ProposedInsight{
				c.New(d, metrics),
				metrics,
			})
		}
	}
	//Returning the resultset
	return result
}
