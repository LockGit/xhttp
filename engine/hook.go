package engine

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func toMegaBytes(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

func StatHook() (err error) {
	tc := time.NewTicker(600 * time.Second)
	defer tc.Stop()
	var ms runtime.MemStats
	for range tc.C {
		n := runtime.NumGoroutine()
		runtime.ReadMemStats(&ms)
		log.Println(fmt.Sprintf("HeapAlloc:%f", toMegaBytes(ms.HeapAlloc)))
		log.Println(fmt.Sprintf("TotalAlloc:%f", toMegaBytes(ms.TotalAlloc)))
		log.Println(fmt.Sprintf("HeapSys:%f", toMegaBytes(ms.HeapSys)))
		log.Println(fmt.Sprintf("HeapIdle:%f", toMegaBytes(ms.HeapIdle)))
		log.Println(fmt.Sprintf("HeapReleased:%f", toMegaBytes(ms.HeapReleased)))
		log.Println(fmt.Sprintf("Goroutine:%d", n))
	}
	return
}
