package main

import (
	"flag"
	"net/url"
	stress "github.com/amarburg/go-lazycache-benchmarking"
)

var hostFlag = flag.String("host", "127.0.0.1:5000", "Host to query")
var pathFlag = flag.String("path", "/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/2017/01/01/CAMHDA301-20170101T000500.mov", "Path to movie to query")

var countFlag = flag.Int("count", 50, "Number of queries to make")
var parallelismFlag = flag.Int("parallelism", 5, "Parallelism of  queries")

func main() {

	flag.Parse()

imageUrl := url.URL{
	Scheme: "http",
	Host: *hostFlag,
	Path: *pathFlag,
}

	err := stress.FrameStress( stress.SetCount( *countFlag ),
 															stress.SetParallelism( *parallelismFlag ),
														  stress.AddUrl( imageUrl.String() ) )

	if err != nil {
		panic(err)
	}

}
