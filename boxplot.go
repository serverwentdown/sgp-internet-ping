package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"gonum.org/v1/gonum/stat"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
)

var flagIn = flag.String("in", "", "input file")
var flagOut = flag.String("out", "boxplot.csv", "output file")

func main() {
	flag.Parse()
	in := *flagIn
	out := *flagOut

	data := make(map[string][]float64, 0)

	log.Println("Reading file " + in)
	raw, err := ioutil.ReadFile(in)
	if err != nil {
		panic(err)
	}

	log.Println("Parsing file " + in)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		panic(err)
	}

	log.Println("File " + in + " read successfully")

	log.Println("Computing params")
	odata := make([][]string, 0)
	for cc, ttls := range data {
		sort.Float64s(ttls)
		min := stat.Quantile(0.00, stat.Empirical, ttls, nil)
		q1 := stat.Quantile(0.25, stat.Empirical, ttls, nil)
		median := stat.Quantile(0.50, stat.Empirical, ttls, nil)
		q3 := stat.Quantile(0.75, stat.Empirical, ttls, nil)
		max := stat.Quantile(1.00, stat.Empirical, ttls, nil)
		odata = append(odata, []string{
			cc,
			strconv.Itoa(len(ttls)),
			strconv.FormatFloat(min, 'f', -1, 64),
			strconv.FormatFloat(q1, 'f', -1, 64),
			strconv.FormatFloat(median, 'f', -1, 64),
			strconv.FormatFloat(q3, 'f', -1, 64),
			strconv.FormatFloat(max, 'f', -1, 64),
		})
	}

	log.Println("Creating file " + out)
	outFile, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	log.Println("Encoding csv")
	w := csv.NewWriter(outFile)
	err = w.WriteAll(odata)
	if err != nil {
		panic(err)
	}
}
