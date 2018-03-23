package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var flagIn = flag.String("in", "", "input file")
var flagOut = flag.String("out", "country.csv", "output file")

func main() {
	flag.Parse()
	in := *flagIn
	out := *flagOut

	data := make(map[string][]int, 0)

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

	log.Println("Converting to csv")
	odata := make([][]string, 0)
	for cc, ttls := range data {
		for _, ttl := range ttls {
			odata = append(odata, []string{
				cc,
				strconv.Itoa(ttl),
			})
		}
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
