package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"strings"
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
	IP  uint32 `json:"ip"`
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
		ip := binary.BigEndian.Uint32(net.ParseIP(o.IP)[12:])
		h := hostLatency{
			IP:  ip,
			TTL: o.Ports[0].TTL,
		}
		odata = append(odata, h)
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
