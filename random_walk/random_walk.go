package main


import (
  stress "github.com/amarburg/go-lazycache-benchmarking"
  "flag"
  "net/url"
  "math/rand"
  "time"
)

var hostFlag = flag.String("host", "127.0.0.1:5000", "Host to query")
var pathFlag = flag.String("path", "/org/oceanobservatories/rawdata/files/", "Root path to query")

var countFlag = flag.Int("count", 50, "Number of queries to make")
var parallelismFlag = flag.Int("parallelism", 5, "Parallelism of  queries")

func main() {

  flag.Parse()

  baseUrl := url.URL{
  	Scheme: "http",
  	Host: *hostFlag,
  	Path: *pathFlag,
  }

  rand.Seed( time.Now().UTC().UnixNano())


  err := stress.RandomWalk( stress.SetCount( *countFlag ),
 															stress.SetParallelism( *parallelismFlag ),
														  stress.AddUrl( baseUrl.String() ) )

	if err != nil {
		panic(err)
	}

}
