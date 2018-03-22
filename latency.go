package main

import (
	"encoding/gob"
	"flag"
	"io/ioutil"
	"log"
	"strings"
    "runtime"
)

type openPort struct {
	IP        string     `json:"ip"`
	Timestamp string     `json:"timestamp"`
	Ports     []portInfo `json:"ports"`
}

type portInfo struct {
	Port   int    `json:"port"`
	Proto  string `json:"proto"`
	Status string `json:"status"`
	Reason string `json:"reason"`
	TTL    int    `json:"ttl"`
}

type hostLatency struct {
	IP  string `json:"ip"`
	TTL int    `json:"ttl"`
}

var flagIn = flag.String("in", "", "comma-seperated list of input files")
var flagOut = flag.String("out", "latency.json", "output file")

func main() {
	flag.Parse()
	in := strings.Split(*flagIn, ",")
	out := *flagOut

	data := make([]openPort, 0)

	for _, fin := range in {
        runtime.GC()
		log.Println("Reading file " + fin)
		raw, err := ioutil.ReadFile(fin)
		if err != nil {
			panic(err)
		}

		log.Println("Fixing file " + fin)
		raw = raw[:len(raw)-3]
		raw[len(raw)-1] = ']'

		log.Println("Parsing file " + fin)
		fdata := make([]openPort, 0)
		err = json.Unmarshal(raw, &fdata)
		if err != nil {
			panic(err)
		}

		log.Println("File " + fin + " read successfully")
		data = append(data, fdata...)
	}

	log.Println("Reducing data to IP and latency")
	odata := make([]hostLatency, 0)
	for _, o := range data {
		h := hostLatency{
			IP:  o.IP,
			TTL: o.Ports[0].TTL,
		}
		odata = append(odata, h)
	}

    data = nil
    runtime.GC()

	log.Println("Encoding json")
	raw, err := json.Marshal(odata)
	if err != nil {
		panic(err)
	}
	log.Println("Writing to file " + out)
	err = ioutil.WriteFile(out, raw, 0)
	if err != nil {
		panic(err)
	}
}
