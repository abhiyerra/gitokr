package main

type Validation struct {
	Key       string
	Validated bool
	Learning  string
}

type Canvas struct {
	Name string

	Problem struct {
		Validations          []Validation
		ExistingAlternatives []Validation
	}
	Solution               []Validation
	KeyMetrics             []Validation
	UniqueValueProposition struct {
		Validations      []Validation
		HighLevelConcept []Validation
	}
	UnfairAdvantage  []Validation
	Channels         []Validation
	CustomerSegments struct {
		Validations   []Validation
		EarlyAdopters []Validation
	}
	CostStructure  []Validation
	RevenueStreams []Validation
}
