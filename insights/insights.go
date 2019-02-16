package insights

import "github.com/cuttle-ai/brain/visualizations"

/*
	This file contains the utilities for generating insights for a dataset
*/

const (
	//CORRELATION is the type string of the correlation type of insight
	CORRELATION = "CORRELATION"
)

//Insight is the interface that has to be implemented by the any type of insight
type Insight interface {
	//New will return a new instance of a insight with the initializations done
	//for the dataset
	New(Dataset, []Metric) Insight
	//Visual is the visualization form in which the insight has to to be
	//presented
	Visual() visualizations.Visual
	//Type of insight like Correlation etc.
	Type() string
	//Relevance will tell whether the insight is relevant or not
	Relevant() bool
	//FSFA is the fast statistical feasibility analysis whether
	//the insight is statistically feasible without running the
	//analysis over the entire data. Calling FSFA followed by Relevant
	//can avoid unnecessary calculations for insight generations
	//even if the domain knowledge proposal states other wise.
	FSFA()
	//Generate will generate the insights for the given data set.
	Generate()
	//Propose will propose a list of possible insight using domain knowledge
	Propose(Dataset) []ProposedInsight
}

//ProposedInsight is the proposal for a potential insight.
type ProposedInsight struct {
	I Insight  //I is the insight to be generated
	M []Metric //M is the list of metrics to be used for generating insights.
}

//GenerateInsights generates the relevant insights for a given dataset
func GenerateInsights(d Dataset) []Insight {
	/*
		We will first use the domain knowledge to propose possible
		metric and insight type combinations.
		Then we will further filter out the proposed insights with short
		//statistical analysis.
		Then we will generate the insights.
		Then we will run the statistical functions to check whether the insight
		is a relevant one or not.
		If valid we will add it to the return result set
	*/
	//variable to store the insights
	result := []Insight{}

	//getting the proposals for the insights
	ps := Propose(d)

	//now we iterate through the insight proposals to validate them.
	for i := range ps {
		//We run the FSFA for each proposal
		//Then check whether it's relevant
		ps[i].I.FSFA()
		if !ps[i].I.Relevant() {
			//The insight isn't relevant. So we will skip the same
			continue
		}
		//Insight is relevant we will generate the same and then check the
		//relevance again
		ps[i].I.Generate()
		if !ps[i].I.Relevant() {
			//Insight isn't relevant after genertaing the insight. Will skip it.
			continue
		}
		//insight is relevant so we will add it to the results
		result = append(result, ps[i].I)
	}

	//Now we ill return the result
	return result
}

//Propose will propose the possible insights for a given data set.
//It uses the domain knowledge for proposing the same.
//This function is WIP and should be used for production purposes.
func Propose(d Dataset) []ProposedInsight {
	/*
		We will get the list of the possible insight types.
		We will simply iterate through the insight types and make decisions
		based on the same.
	*/
	//variable to store the comboined insights
	result := []ProposedInsight{}

	//getting the insight type lists
	ins := Insights()

	//now iterating through each to produce the proposals
	for i := range ins {
		result = append(result, ins[i].Propose(d)...)
	}

	return result
}

//Insights returns the list of insight types available in the system
func Insights() []Insight {
	return []Insight{
		&Correlation{},
	}
}
