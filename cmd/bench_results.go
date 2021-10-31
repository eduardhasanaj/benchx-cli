package cmd

import "github.com/wcharczuk/go-chart/v2"

type BenchResults struct {
	OS   string
	Arch string
	CPU  string

	it     map[string][]chart.Value
	speed  map[string][]chart.Value
	memory map[string][]chart.Value
	alloc  map[string][]chart.Value
}

func newBenchResults() *BenchResults {
	return &BenchResults{
		it:     make(map[string][]chart.Value),
		speed:  make(map[string][]chart.Value),
		memory: make(map[string][]chart.Value),
		alloc:  make(map[string][]chart.Value),
	}
}

func (br *BenchResults) addSpeedStat(gr string, v chart.Value) {
	slic, ok := br.speed[gr]
	if !ok {
		slic = make([]chart.Value, 0, 2)
	}

	slic = append(slic, v)

	br.speed[gr] = slic
}

func (br *BenchResults) addMemStat(gr string, v chart.Value) {
	slic, ok := br.memory[gr]
	if !ok {
		slic = make([]chart.Value, 0, 2)
	}

	slic = append(slic, v)

	br.memory[gr] = slic
}

func (br *BenchResults) addAllocStat(gr string, v chart.Value) {
	slic, ok := br.alloc[gr]
	if !ok {
		slic = make([]chart.Value, 0, 2)
	}

	slic = append(slic, v)

	br.alloc[gr] = slic
}

func (br *BenchResults) addItStat(gr string, v chart.Value) {
	slic, ok := br.it[gr]
	if !ok {
		slic = make([]chart.Value, 0, 2)
	}

	slic = append(slic, v)

	br.it[gr] = slic
}
