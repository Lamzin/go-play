package main

import (
	"log"
	"time"

	clientv2 "github.com/influxdata/influxdb/client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var client clientv2.Client

func main() {

	influxDbInit()

	for {
		sendStat()
		log.Println("Success send metrics")
		//time.Sleep(200 * time.Millisecond)
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
		Database: "asus",
	})
	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoint(makeCpuPoint())
	bp.AddPoint(makeRAMPoint())
	bp.AddPoint(makeSwapPoint())

	// Write the batch
	if err := client.Write(bp); err != nil {
		log.Fatal(err)
	}
}

func getRAMStat() (total, used, available int64, usedPercent float64) {
	m, _ := mem.VirtualMemory()
	return int64(m.Total), int64(m.Used), int64(m.Free), m.UsedPercent
}

func getSwapStat() (total, used, available int64, usedPercent float64) {
	m, _ := mem.SwapMemory()
	return int64(m.Total), int64(m.Used), int64(m.Free), m.UsedPercent
}

func getCpuStat() (float64, []float64) {
	total, _ := cpu.Percent(100 * time.Millisecond, false)
	all, _ := cpu.Percent(100 * time.Millisecond, true)
	return total[0], all
}

func makeRAMPoint() *clientv2.Point {
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

func makeSwapPoint() *clientv2.Point {
	total, used, available, usedPercent := getSwapStat()

	tags := map[string]string{
		"resource": "swap",
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

func makeCpuPoint() *clientv2.Point {
	total, all := getCpuStat()

	tags := map[string]string{
		"resource": "cpu",
	}
	fields := map[string]interface{}{
		"cpu":  total,
		"cpu0": all[0],
		"cpu1": all[1],
		"cpu2": all[2],
		"cpu3": all[3],
	}
	point, err := clientv2.NewPoint("resources", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	return point
}
