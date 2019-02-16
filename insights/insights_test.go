package insights

import (
	"testing"

	"github.com/cuttle-ai/brain/visualizations"
)

/*
	This file contains the tests for the insights genertaion functions
*/

func TestInsights(t *testing.T) {
	ins := Insights()
	if len(ins) != 1 {
		t.Fatal("Expected to support 1 insight. But got", len(ins))
	}
}

type proposeTC struct {
	ID          string
	Description string
	Dataset     Dataset
	Metrics     []Metric
	Datas       []interface{}
	Expected    []ProposedInsight
}

var proposeTCs = []proposeTC{
	{"1", "Normal case", NewDataset(), []Metric{
		{
			Name:        "Age",
			DataType:    Float,
			DisplayName: "Age",
		},
		{
			Name:        "Height",
			DataType:    Float,
			DisplayName: "Height",
		},
	}, []interface{}{
		[]float64{10, 20, 30},
		[]float64{150, 180, 188},
	}, []ProposedInsight{
		{
			&Correlation{},
			[]Metric{
				{
					Name:        "Age",
					DataType:    Float,
					DisplayName: "Age",
				},
				{
					Name:        "Height",
					DataType:    Float,
					DisplayName: "Height",
				},
			},
		},
	}},
}

func TestPropose(t *testing.T) {
	for _, v := range proposeTCs {
		t.Run(v.ID, func(t *testing.T) {
			d := v.Dataset
			for i, m := range v.Metrics {
				err := d.AddMetric(m, v.Datas[i])
				if err != nil {
					t.Fatal("Error while adding metric in testcase", v.ID, err)
				}
			}
			ps := Propose(d)
			//Doing the length check
			if len(ps) != len(v.Expected) {
				t.Fatal("Expected", len(v.Expected), "proposals. Got", len(ps))
			}
		})
	}
}

type generateInsightsTC struct {
	ID          string
	Description string
	Dataset     Dataset
	Metrics     []Metric
	Datas       []interface{}
	Expected    []Insight
}

var generateInsightsTCs = []generateInsightsTC{
	{"1", "Normal case", NewDataset(), []Metric{
		{
			Name:        "Age",
			DataType:    Float,
			DisplayName: "Age",
		},
		{
			Name:        "Height",
			DataType:    Float,
			DisplayName: "Height",
		},
	}, []interface{}{
		[]float64{10, 20, 30},
		[]float64{150, 180, 188},
	}, []Insight{
		&Correlation{
			visual: visualizations.ScatterPlot{
				T: "Age and Height",
				M: []visualizations.Metric{
					{
						Name:        "Age",
						DisplayName: "Age",
						DataType:    Float,
						Dimension:   0,
					},
					{
						Name:        "Height",
						DisplayName: "Height",
						DataType:    Float,
						Dimension:   1,
					},
				},
			},
			relevant: true,
			ms: []Metric{
				{
					Name:        "Age",
					DataType:    Float,
					DisplayName: "Age",
				},
				{
					Name:        "Height",
					DataType:    Float,
					DisplayName: "Height",
				},
			},
		},
	}},
	{"2", "Unrelated Metrics", NewDataset(), []Metric{
		{
			Name:        "Age",
			DataType:    Float,
			DisplayName: "Age",
		},
		{
			Name:        "Height",
			DataType:    Float,
			DisplayName: "Height",
		},
	}, []interface{}{
		[]float64{10, 20, 30},
		[]float64{150, 0, 100},
	}, []Insight{},
	},
}

func TestGenerateInsights(t *testing.T) {
	for _, v := range generateInsightsTCs {
		t.Run(v.ID, func(t *testing.T) {
			d := v.Dataset
			for i, m := range v.Metrics {
				err := d.AddMetric(m, v.Datas[i])
				if err != nil {
					t.Fatal("Error while adding metric in testcase", m.Name, v.ID, err)
				}
			}
			ps := GenerateInsights(d)
			//Doing the length check
			if len(ps) != len(v.Expected) {
				t.Fatal("Expected", len(v.Expected), "proposals. Got", len(ps))
			}
		})
	}
}
