package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
    "math/rand"
)

type geoRange struct {
	Start uint32
	End   uint32
	CC    string
}

var flagIn = flag.String("in", "", "comma-seperated list of input files")
var flagOut = flag.String("out", "country.json", "output file")
var flagDB = flag.String("db", "iptocountry", "geographical database. either geolite or iptocountry")

func main() {
	flag.Parse()
	in := strings.Split(*flagIn, ",")
	out := *flagOut

	data := make(map[uint32]int, 0)

	for _, fin := range in {
		log.Println("Reading file " + fin)
		raw, err := ioutil.ReadFile(fin)
		if err != nil {
			panic(err)
		}

		log.Println("Parsing file " + fin)
		fdata := make(map[uint32]int, 0)
		err = json.Unmarshal(raw, &fdata)
		if err != nil {
			panic(err)
		}

		log.Println("File " + fin + " read successfully")
		for k, v := range fdata {
			data[k] = v
		}
	}

	geo := make([]geoRange, 0)
	if *flagDB == "geolite" { /*
			// incomplete
			log.Println("Reading geolite files")
			rangeFile, err := ioutil.ReadFile("data/GeoLite2-Country-Blocks-IPv4.csv")
			if err != nil {
				panic(err)
			}
			countryFile, err := ioutil.ReadFile("data/GeoLite2-Country-Locations-en.csv")
			if err != nil {
				panic(err)
			}

			countryLines := bytes.Split(countryFile, []byte{'\n'})
			for i, line := range countryLines {

			}*/
	}
	if *flagDB == "iptocountry" {
		log.Println("Read iptocountry files")
		countryFile, err := os.Open("data/Webnet77-IpToCountry.csv")
		if err != nil {
			panic(err)
		}

		r := csv.NewReader(countryFile)
		records, err := r.ReadAll()
		if err != nil {
			panic(err)
		}

		for _, record := range records {
			cc := record[4]
			start, err := strconv.ParseUint(record[0], 10, 32)
			if err != nil {
				panic(err)
			}
			end, err := strconv.ParseUint(record[1], 10, 32)
			if err != nil {
				panic(err)
			}

			g := geoRange{
				Start: uint32(start),
				End:   uint32(end),
				CC:    cc,
			}
			geo = append(geo, g)
		}
	}

    // TODO: verify correctness
	log.Println("Grouping latencies by country")
	odata := make(map[string][]int, 0)
	j := 0
	for ip, ttl := range data {
		for geo[j].End <= ip {
			j += 1
		}
		if ip < geo[j].Start {
			j -= 1
			odata["Unknown"] = append(odata["Unknown"], ttl)
			continue
		}
		cc := geo[j].CC

        if rand.Intn(100) == 0 {
		    odata[cc] = append(odata[cc], ttl)
        }
	}

	total := 0
	for k, v := range odata {
		log.Println("Country " + k + " has " + strconv.Itoa(len(v)) + " samples")
		total += len(v)
	}
	log.Println("Total: " + strconv.Itoa(total))

	log.Println("Encoding json")
	raw, err := json.Marshal(odata)
	if err != nil {
		panic(err)
	}
	log.Println("Writing to file " + out)
	err = ioutil.WriteFile(out, raw, 0644)
	if err != nil {
		panic(err)
	}
}
