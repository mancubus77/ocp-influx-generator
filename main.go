package main

import (
	"fmt"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	client := influxdb2.NewClientWithOptions("http://vminsert-example-vmcluster-persistent:8480/insert/0/influx/", "my-token",
		influxdb2.DefaultOptions().SetBatchSize(1000))
	// Get non-blocking write client
	writeAPI := client.WriteAPI("my-org", "my-bucket")
	// write some points
	for {
		for i := 0; i < 10; i++ {
			// create point
			p := influxdb2.NewPoint(
				"system",
				map[string]string{
					"id":       fmt.Sprintf("rack_%v", i%10),
					"vendor":   "AWS",
					"hostname": fmt.Sprintf("host_%v", i%100),
				},
				map[string]interface{}{
					"temperature": rand.Float64() * 80.0,
					"disk_free":   rand.Float64() * 1000.0,
					"disk_total":  (i/10 + 1) * 1000000,
					"mem_total":   (i/100 + 1) * 10000000,
					"mem_free":    rand.Uint64(),
				},
				time.Now())
			// write asynchronously
			writeAPI.WritePoint(p)
		}

		// Force all unwritten data to be sent
		writeAPI.Flush()
	}
	// Ensures background processes finishes
	client.Close()
}
