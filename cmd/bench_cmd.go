package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	exec "golang.org/x/sys/execabs"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

var statType = map[string]bool{
	"iterations": true,
	"ns/op":      true,
	"b/op":       true,
	"allocs/op":  true,
}

type BenchCommand struct {
	Groups []string `arg:"" name:"groups" help:"Benchmark groups"`
}

func (c *BenchCommand) Run(ctx *Context) error {
	dir := "./graphs"
	if err := ensureGraphsDir(dir); err != nil {
		return err
	}

	cmd := exec.Command("go", "test", "-bench", ".", "-benchmem")
	out, err := cmd.Output()
	if err != nil {
		return errors.New(string(out))
	}

	r := bytes.NewReader(out)

	res, err := c.parse(r)
	if err != nil {
		return err
	}

	if err = renderStats("Iterations", "", res.it, true); err != nil {
		return err
	}

	if err = renderStats("Memory Per Operation", "b/op", res.memory, true); err != nil {
		return err
	}

	if err = renderStats("Memory Allocations", "allocs/op", res.alloc, true); err != nil {
		return err
	}

	if err = writeOutput(dir, out); err != nil {
		return err
	}

	if err = renderStats("Speed", "ns/op", res.speed, false); err != nil {
		return err
	}

	return nil
}

func writeOutput(dir string, b []byte) error {
	f, err := os.Create(dir + "/output.txt")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func renderStats(title, yLabel string, data map[string][]chart.Value, isInteger bool) error {
	for k, v := range data {
		if err := makeGraph(k+" "+title, yLabel, v, isInteger); err != nil {
			return err
		}
	}

	return nil
}

func makeGraph(title, vertTitle string, values []chart.Value, isInteger bool) error {
	graph := chart.BarChart{
		Title: addSpacesBeforeUpper(title),
		TitleStyle: chart.Style{
			FontSize: 20,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		Width:    -1,
		Height:   512,
		BarWidth: 60,
		Bars:     values,
	}

	if len(vertTitle) > 0 {
		graph.YAxis = chart.YAxis{
			Name: vertTitle,
			NameStyle: chart.Style{
				FontSize: 14,
				Padding: chart.Box{
					Left:  20,
					Right: 40,
				},
			},
			Style: chart.Style{
				Hidden:      false,
				StrokeColor: drawing.ColorBlack,
				FontColor:   drawing.ColorBlack,
			},
		}
	}

	if isInteger {
		graph.YAxis.ValueFormatter = chart.IntValueFormatter
	}

	f, err := os.Create("./graphs/" + title + ".png")
	if err != nil {
		return err
	}
	defer f.Close()
	graph.Render(chart.PNG, f)

	return nil
}

func (c *BenchCommand) parse(r io.Reader) (*BenchResults, error) {
	br := newBenchResults()

	s := bufio.NewScanner(r)

	// Parse OS info
	s.Scan()
	_, br.OS = parseKV(s.Text(), ": ")

	// Scan arch
	s.Scan()
	_, br.Arch = parseKV(s.Text(), ": ")

	// Skip package
	s.Scan()

	// Read CPU description
	s.Scan()
	_, br.CPU = parseKV(s.Text(), ": ")

	for s.Scan() {
		t := s.Text()
		if strings.Contains(t, "PASS") {
			break
		}

		var desc string

		var itCount float64

		var nsopAmount float64
		var nsopLabel string

		var bopAmount float64
		var bopLabel string

		var allocs float64
		var allocsLabel string

		fmt.Sscan(t, &desc, &itCount, &nsopAmount, &nsopLabel, &bopAmount, &bopLabel,
			&allocs, &allocsLabel)
		gr, name, err := parseGrAndName(desc, c.Groups)
		if err != nil {
			continue
			// return nil, err
		}

		br.addItStat(gr, chart.Value{
			Value: itCount,
			Label: name,
		})

		br.addSpeedStat(gr, chart.Value{
			Value: nsopAmount,
			Label: name,
		})

		br.addMemStat(gr, chart.Value{
			Value: bopAmount,
			Label: name,
		})

		br.addAllocStat(gr, chart.Value{
			Value: allocs,
			Label: name,
		})
	}

	return br, nil
}

func parseGrAndName(desc string, groups []string) (string, string, error) {
	for _, gr := range groups {
		if strings.Contains(desc, gr) {
			desc = strings.Replace(desc, "Benchmark", "", 1)
			desc = strings.Replace(desc, gr, "", 1)
			slashIndex := strings.Index(desc, "/")
			if slashIndex == -1 {
				slashIndex = strings.Index(desc, "-")
			}

			if slashIndex != -1 {
				desc = desc[0:slashIndex]
			}

			name := strings.ToLower(desc)

			return gr, name, nil
		}
	}

	return "", "", errors.New("could not extract bench group from: " + desc)
}

func parseKV(kvStr string, delimiter string) (string, string) {
	arr := strings.Split(kvStr, delimiter)
	if len(arr) != 2 {
		return "", ""
	}

	return arr[0], arr[1]
}

func addSpacesBeforeUpper(str string) string {
	sb := strings.Builder{}

	lastLower := true
	for _, v := range str {
		isUpper := unicode.IsUpper(v)
		if isUpper && lastLower {
			sb.WriteString(" ")
		}

		sb.WriteRune(v)

		lastLower = !isUpper
	}

	return sb.String()
}

func ensureGraphsDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, 0755)
	}

	return nil
}
