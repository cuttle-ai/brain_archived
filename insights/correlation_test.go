package insights

import (
	"reflect"
	"testing"

	"github.com/cuttle-ai/brain/visualizations"
)

/*
	This file contains the tests for the correlation insight
*/

func TestCorrelation_New(t *testing.T) {
	d := NewDataset()
	m := Metric{Name: "car", DataType: String}
	d.AddMetric(m, []string{"Tesla", "Audi", "Suzuki"})
	ci := (&Correlation{}).New(d, []Metric{m})
	c, ok := ci.(*Correlation)
	//checking the insight given the New function is correlation itself
	if !ok {
		t.Fatal("Expected a correlation. Got", reflect.TypeOf(ci))
	}
	//checking the length of metric array of the dfatatset of correlation
	if len(c.dt.Metrics) != 1 {
		t.Fatal("Expected only 1 metric. Got", len(c.dt.Metrics))
	}
	//Checking the length of the datatset of c
	if c.dt.Length != 3 {
		t.Fatal("Expected length of dataset is 3. Got", c.dt.Length)
	}
	//Checking the length of metrics in the correlation
	if len(c.ms) != 1 {
		t.Fatal("Expected length of metrics in the correlation is 1. Got",
			len(c.ms))
	}
}

func TestCorrelation_Visual(t *testing.T) {
	t.Run("Testing visual for nil", func(t *testing.T) {
		cr := &Correlation{}
		if cr.Visual() != nil {
			t.Fatal("Expected visual as nil. Got it as",
				reflect.TypeOf(cr.Visual()))
		}
	})

	t.Run("Testing visual for not nil", func(t *testing.T) {
		cr := &Correlation{visual: visualizations.ScatterPlot{}}
		if cr.Visual() == nil {
			t.Fatal("Expected visual as scatter plot. Got it as nil")
		}
	})

}

func TestCorrelation_Type(t *testing.T) {
	cr := &Correlation{}
	if cr.Type() != CORRELATION {
		t.Fatal("Expected insight type is", CORRELATION, "Got", cr.Type())
	}
}

func TestCorrelation_Relevant(t *testing.T) {
	t.Run("Testing relevance for false", func(t *testing.T) {
		cr := &Correlation{}
		if cr.Relevant() {
			t.Fatal("Expected relevance to be false. Got true")
		}
	})

	t.Run("Testing relevance fo true", func(t *testing.T) {
		cr := &Correlation{relevant: true}
		if !cr.Relevant() {
			t.Fatal("Expected relevance to be true. Got false")
		}
	})
}

func TestCorrelation_FSFA(t *testing.T) {
	t.Run("Testing FSFA when metrics < 2", func(t *testing.T) {
		cr := &Correlation{ms: []Metric{{Name: "age"}}}
		cr.FSFA()
		if cr.Relevant() {
			t.Fatal("Expected correlation to be irrelevant with 1 metric. Got",
				"it as relevant")
		}
	})

	t.Run("Testing FSFA when metric data type not float", func(t *testing.T) {
		cr := &Correlation{ms: []Metric{
			{Name: "age", DataType: Float},
			{Name: "name", DataType: String},
		}}
		cr.FSFA()
		if cr.Relevant() {
			t.Fatal("Expected correlation to be irrelevant with not float",
				"data type. Got it as relevant")
		}
	})

	t.Run("Testing FSFA in normal conditions", func(t *testing.T) {
		cr := &Correlation{ms: []Metric{
			{Name: "age", DataType: Float},
			{Name: "height", DataType: Float},
		}}
		cr.FSFA()
		if !cr.Relevant() {
			t.Fatal("Expected correlation to be relevant with normal",
				"conditions Got it as irrelevant")
		}
	})
}

type cGenerateTC struct {
	ID              string
	Description     string
	Dataset         Dataset
	Metrics         []Metric
	Datas           []interface{}
	GenerateMetrics []Metric
	Relevance       bool
}

var cGenerateTCs = []cGenerateTC{
	{"1", "Failure in correlation due to unsupported data type", NewDataset(),
		[]Metric{
			{Name: "age", DataType: Float},
			{Name: "fheight", DataType: Float},
		}, []interface{}{
			[]float64{10, 20},
			[]float64{1.5, 1.8},
		}, []Metric{
			{Name: "age", DataType: Float, Index: 0},
			{Name: "height", DataType: String, Index: 1},
		}, false},
	{"2", "Failure in correlation due to weaker correlation", NewDataset(),
		[]Metric{
			{Name: "age", DataType: Float},
			{Name: "height", DataType: Float},
		}, []interface{}{
			[]float64{10, 20, 21},
			[]float64{1.5, 0.8, 100},
		}, []Metric{
			{Name: "age", DataType: Float, Index: 0},
			{Name: "height", DataType: Float, Index: 1},
		}, false},
	{"3", "Correlated metrics", NewDataset(),
		[]Metric{
			{Name: "age", DataType: Float},
			{Name: "height", DataType: Float},
		}, []interface{}{
			[]float64{10, 20},
			[]float64{1.5, 1.8},
		}, []Metric{
			{Name: "age", DataType: Float, Index: 0},
			{Name: "height", DataType: Float, Index: 1},
		}, true},
}

func TestCorrelation_Generate(t *testing.T) {
	t.Run("Testing generate when correlation is irrelevant",
		func(t *testing.T) {
			cr := &Correlation{}
			cr.Generate()
			if cr.Relevant() {
				t.Fatal("Expected generation to be irrelvant when the insight",
					"is irrelvant. Got it relevant")
			}
		})

	t.Run("Testing generate when correlation when insufficient metrics",
		func(t *testing.T) {
			cr := &Correlation{relevant: true}
			cr.Generate()
			if cr.Relevant() {
				t.Fatal("Expected generation to be irrelvant there is",
					"insufficient metrics. Got it relevant")
			}
		})

	//iterating through the testcases
	for _, v := range cGenerateTCs {
		t.Run(v.ID, func(t *testing.T) {
			d := v.Dataset
			for i, m := range v.Metrics {
				err := d.AddMetric(m, v.Datas[i])
				if err != nil {
					t.Fatal("Error while adding metric in testcase", v.ID, err)
				}
			}
			c := &Correlation{ms: v.GenerateMetrics, dt: d, relevant: true}
			c.Generate()

			if v.Relevance != c.Relevant() {
				t.Fatal("Expected relevance of insight", v.Relevance, "Got",
					c.Relevant(), v.ID)
			}

			if v.Relevance && c.Visual() == nil {
				t.Fatal("Expected relevance but didn't generate a visual", v.ID)
			}
		})
	}
}

type cProposeTC struct {
	ID          string
	Description string
	Dataset     Dataset
	Metrics     []Metric
	Data        []interface{}
	Expected    []ProposedInsight
}

var cProposeTCs = []cProposeTC{
	{"1", "Normal case", NewDataset(), []Metric{
		{Name: "age", DataType: Float},
		{Name: "height", DataType: Float},
	}, []interface{}{
		[]float64{1, 2, 3},
		[]float64{1, 2, 3},
	}, []ProposedInsight{
		{&Correlation{}, []Metric{
			{Name: "age", DataType: Float},
			{Name: "height", DataType: Float},
		}},
	},
	},
	{"2", "Non float metrics", NewDataset(), []Metric{
		{Name: "age", DataType: Float},
		{Name: "height", DataType: String},
	}, []interface{}{
		[]float64{1, 2, 3},
		[]string{"1", "2", "3"},
	}, []ProposedInsight{},
	},
}

func TestCorrelation_Propose(t *testing.T) {
	for _, v := range cProposeTCs {
		t.Run(v.ID, func(t *testing.T) {
			d := v.Dataset
			for i := range v.Metrics {
				err := d.AddMetric(v.Metrics[i], v.Data[i])
				if err != nil {
					t.Fatal("Error while adding metric for", v.ID, err)
				}
			}
			c := &Correlation{}
			pro := c.Propose(d)
			if len(pro) != len(v.Expected) {
				t.Fatal("Expected proposals", v.Expected, "Got", pro)
			}
		})
	}
}
