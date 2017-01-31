package main

import (
	//"encoding/json"
	"fmt"
	"time"
	"math/rand"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	for {
		printStat()
		go fffff()
		time.Sleep(100 * time.Millisecond)
	}
}

func fffff() {
	for i := 0; i < 100000000;{
		i = i + rand.Int() % 5
	}
}

func printStat() {
	infos, err := mem.VirtualMemory()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(*infos)
	//infos.

	//info := infos[0]


	//fmt.Println("user: ", info.User / info.Total() * 100.0)
	//fmt.Println("system: ", info.System / info.Total() * 100.0)
	//fmt.Println("idle: ", info.Idle / info.Total() * 100.0)
}