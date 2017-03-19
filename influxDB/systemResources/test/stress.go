package main

import (
	"log"
	"time"

	clientv2 "github.com/influxdata/influxdb/client/v2"
	"github.com/shirou/gopsutil/mem"
)

var client clientv2.Client

const (
	BatchSize  = 1000
	GoRoutines = 5
)

var count = 0

func main() {

	influxDbInit()

	for i := 0; i < GoRoutines; i++ {
		go func() {
			for {
				sendStat()
				count++
			}
		}()
	}

	prev := 0
	for {
		log.Printf("Write %d points, goroutine: %d, batch size: %d\n", (count - prev) * BatchSize, GoRoutines, BatchSize)
		prev = count
		time.Sleep(time.Second)
	}
}

func influxDbInit() {
	var err error
	client, err = clientv2.NewHTTPClient(clientv2.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("success connected to influxDB")
}

func sendStat() {
	// Create a new point batch
	bp, err := clientv2.NewBatchPoints(clientv2.BatchPointsConfig{
		Database: "test",
	})
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < BatchSize; i++ {
		bp.AddPoint(makePoint())
	}

	// Write the batch
	if err := client.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func makePoint() *clientv2.Point {
	total, used, available, usedPercent := getRAMStat()

	tags := map[string]string{
		"resource": "ram",
	}
	fields := map[string]interface{}{
		"total":        total,
		"used":         used,
		"available":    available,
		"used_percent": usedPercent,
	}
	point, err := clientv2.NewPoint("resources", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	return point
}

func getRAMStat() (total, used, available int64, usedPercent float64) {
	m, _ := mem.VirtualMemory()
	return int64(m.Total), int64(m.Used), int64(m.Free), m.UsedPercent
}
