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
)

type hostLatency struct {
	IP  uint32 `json:"ip"`
	TTL int    `json:"ttl"`
}

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

	data := make([]hostLatency, 0)

	for _, fin := range in {
		log.Println("Reading file " + fin)
		raw, err := ioutil.ReadFile(fin)
		if err != nil {
			panic(err)
		}

		log.Println("Parsing file " + fin)
		fdata := make([]hostLatency, 0)
		err = json.Unmarshal(raw, &fdata)
		if err != nil {
			panic(err)
		}

		log.Println("File " + fin + " read successfully")
		data = append(data, fdata...)
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

	log.Println("Grouping latencies by country")
	odata := make(map[string][]int, 0)
	for _, o := range data {
		ip := o.IP
		ttl := o.TTL
		found := false
		for _, g := range geo {
			if g.Start < ip && ip < g.End {
				odata[g.CC] = append(odata[g.CC], ttl)
				found = true
				break
			}
		}
		if !found {
			odata["Unknown"] = append(odata["Unknown"], ttl)
		}
	}

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
