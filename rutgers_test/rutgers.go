package main

import (
	"fmt"
	"github.com/amarburg/go-lazycache-benchmarking"
	"os"
)

func main() {

	settings, err := lazycache_benchmarking.StressFlags(
                      lazycache_benchmarking.SetUrls( []string{"https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T000000Z.mp4"} ))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	lazycache_benchmarking.RepeatedDownload(settings)
}
