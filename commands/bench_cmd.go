package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	exec "golang.org/x/sys/execabs"

	"github.com/wcharczuk/go-chart/v2"
)

const (
	DefaultGr = "default_group"
)

var statType = map[string]bool{
	"iterations": true,
	"ns/op":      true,
	"b/op":       true,
	"allocs/op":  true,
}

type BenchCommand struct {
	Groups []string `name:"groups" help:"Benchmark groups"`

	groupDef bool
	dir      string
}

func (c *BenchCommand) Run(ctx *Context) error {
	c.dir = "./benchx-graphs"
	if err := ensureGraphsDir(c.dir); err != nil {
		return err
	}

	if len(c.Groups) > 0 {
		c.groupDef = true
	}

	cmd := exec.Command("go", "test", "./...", "-bench", ".", "-benchmem")
	out, err := cmd.Output()
	if err != nil {
		return errors.New(string(out))
	}

	r := bytes.NewReader(out)

	res, err := c.parse(r)
	if err != nil {
		return err
	}

	if err = c.renderStats("Iterations", res.it, true); err != nil {
		return err
	}

	if err = c.renderStats("b/op", res.memory, true); err != nil {
		return err
	}

	if err = c.renderStats("allocs/op", res.alloc, true); err != nil {
		return err
	}

	if err = c.renderStats("ns/op", res.speed, false); err != nil {
		return err
	}

	if err = writeOutput(c.dir, out); err != nil {
		return err
	}

	return nil
}

func (c *BenchCommand) renderStats(statType string, data map[string][]chart.Value, isInteger bool) error {
	for k, v := range data {
		var title string

		if c.groupDef {
			title = fmt.Sprintf("%s %s", k, statType)
		} else {
			title = statType
		}

		if err := c.makeGraph(title, v, isInteger); err != nil {
			return err
		}
	}

	return nil
}

func (c *BenchCommand) makeGraph(title string, values []chart.Value, isInteger bool) error {
	graph := chart.BarChart{
		Title: title,
		TitleStyle: chart.Style{
			FontSize: 20,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 60,
			},
		},
		Width:    resolveChartWidth(len(values), 60),
		Height:   512,
		BarWidth: 60,
		Bars:     values,
	}

	if isInteger {
		graph.YAxis.ValueFormatter = chart.IntValueFormatter
	}

	fileName := strings.ReplaceAll(title, " ", "_")
	fileName = strings.ReplaceAll(fileName, "/", "_")
	path := fmt.Sprintf("%s/%s.png", c.dir, strings.ToLower(fileName))

	f, err := os.Create(path)
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

		group, name, err := c.parseGrAndName(desc)

		if err != nil {
			return nil, err
		}

		br.addItStat(group, chart.Value{
			Value: itCount,
			Label: name,
		})

		br.addSpeedStat(group, chart.Value{
			Value: nsopAmount,
			Label: name,
		})

		br.addMemStat(group, chart.Value{
			Value: bopAmount,
			Label: name,
		})

		br.addAllocStat(group, chart.Value{
			Value: allocs,
			Label: name,
		})
	}

	return br, nil
}

func writeOutput(dir string, b []byte) error {
	f, err := os.Create(dir + "/output.txt")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func (c *BenchCommand) parseGrAndName(desc string) (string, string, error) {
	var name string
	if !c.groupDef {
		fmt.Sscanf(desc, "Benchmark%s-", &name)
		name = strings.ToLower(name)
		return DefaultGr, name, nil
	}

	for _, gr := range c.Groups {
		if strings.Contains(desc, gr) {
			var name string
			fmt.Sscanf(desc, "Benchmark"+gr+"%s-", &name)
			name = strings.ToLower(name)

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

func ensureGraphsDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, 0755)
	}

	return nil
}

func resolveChartWidth(barCount, barWidth int) int {
	return barCount * (150 + barWidth)
}
