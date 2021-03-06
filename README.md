
<h1 align="center">
    sgp-internet-ping
</h1>

<p align="center">
    A quick analysis of Internet latency from Singapore to the rest of the world.
</p>

<p align="center">
    <img src="poster.png" alt="poster">
</p>

## Scanning the Internet

To scan the IPv4 Internet, I used the tool [`masscan`](https://github.com/robertdavidgraham/masscan) by security researcher Robert Graham. He has extensively made use of `masscan` in the security research he does. It provides latency measurements up to the millisecond. 

> WARNING: Only scan the Internet if your service provider approves of it. It can cause networking issues.

I wrote a configuration file to run `masscan` on the entire Internet with the included exclusion list, and scan the top 5 open ports according to [speedguide.net](https://www.speedguide.net/ports_common.php). It is stored as `scan.conf`. I started the scan with: 

```
masscan -c scan.conf
```

This produces the output file `scan.bin`. You might want to make use of shards to scan only a portion of the internet per file as scanning the entire internet will produce a huge file that cannot be parsed unless you have enough RAM available. 

## Counting latency

I need to reduce the data for the five ports per host into a single latency reading for each host.

But first, I had to convert the scan binary into JSON:

```
masscan --readscan scan.bin -oJ scan.json
go run latency.go -in scan.json -out latency.json
```

## Categorising readings by country

To associate an IP address with a country, a geolocation lookup database must be used. I used two sources of geolocation databases; Maxmind's GeoLite2, and Webnet77's IPToCountry; to test the accuracy of either database. 

Next, I wrote and used a Go script to group the scans by country:

```
go run country.go -in latency.json -db iptocountry -out country.json
```

## Plotting the latency from Singapore by country

I will start with a simple plot of latency to every country from Singapore. 

```
go run boxplot.go -in country.json -out boxplot.csv
```

## 
