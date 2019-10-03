package asyncparser

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"gonum.org/v1/gonum/stat"
)

type Formatter struct {
	Sizes    []int
	SizeData [][]float64
}

func NewFormatter() *Formatter {
	var f Formatter

	return &f
}

func (f *Formatter) AddSizeResults(size int, results []float64) {
	f.Sizes = append(f.Sizes, size)
	f.SizeData = append(f.SizeData, results)
}

func (f *Formatter) FormatSizes() {
	data := [][]string{}
	for i, size := range f.Sizes {
		// Data needs to be sorted for the statistics functions to work
		sort.Float64s(f.SizeData[i])
		data = append(data, []string{
			fmt.Sprintf("%d", size),
			fmt.Sprintf("%.2fms", stat.Mean(f.SizeData[i], nil)),
			fmt.Sprintf("%.2fms", stat.Quantile(0.5, stat.Empirical, f.SizeData[i], nil)),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Size", "Mean", "Median"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
