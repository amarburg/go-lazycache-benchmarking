package lazycache_benchmarking


import (
  "sync"
  "fmt"
  "bytes"
    "net/http"
    "time"
    "encoding/json"
)

type UrlProducer interface {
  Url( i int ) string
}

func RepeatedDownload( settings *StressOptions,
                        producer UrlProducer   ) error {

  parallelism := settings.Parallelism()
  count       := settings.Count()

  if parallelism > count {
    fmt.Printf("Count %d less than parallelism %d, setting parallelism to %d\n", count, parallelism, count )
    parallelism = count
  }

  var wg sync.WaitGroup

    var urls = make(chan string)
    var results = make(chan []byte, count )

    wg.Add( parallelism )

  	for i := 0; i < parallelism; i++ {
  		go DownloadWorker(urls, results, &wg)
  	}

    start := time.Now()
    for i := 0; i < count; i ++ {
      urls <- producer.Url(i)
    }

    close(urls)
    fmt.Printf("Waiting for workers to finish...\n")
    wg.Wait()

    duration := time.Since( start )

    fmt.Printf("Pulling results\n")
    fmt.Println("{ Data: [")
    for i := 0; i < count; i ++ {
      result := <- results
      fmt.Printf("%s,\n", bytes.NewBuffer(result).String())
    }
    fmt.Println("],")
    fmt.Printf("TotalTime: %f }\n", duration.Seconds() )

    return nil
  }



  func DownloadWorker(urls chan string, results chan []byte , wg *sync.WaitGroup) {
  	//fmt.Println("In random walker")
  	for {

      url,ok := <- urls

      if !ok {
        fmt.Println("Channel closed, quitting")
        wg.Done()
        return
      }

  		fmt.Println("Download worker Querying URL", url)

      start := time.Now()
  		resp, err := http.Get(url)
  		if err != nil {
  			fmt.Printf("%d: ERROR: %s\n", url, err)
  			fmt.Printf("Error making request: %s\n", err.Error())
        buf,_ := json.Marshal( struct {
          Url         string
          StatusCode  int
        }{
        url,
        resp.StatusCode,
      } )
      results <- buf

  		}

    fmt.Println(resp)

      //
  		defer resp.Body.Close()

      var bytesRead int64
      mb := make( []byte, 1024*1024 )
      for {
        n,err := resp.Body.Read(mb)
        bytesRead += int64(n)
        if err != nil { break }
      }
      duration := time.Since(start)

      // fmt.Println("Done fetching ", url, " elapsed ", duration )

      buf,_ := json.Marshal( struct {
        Url string
        ContentLength int64
        BytesRead      int64
        ElapsedTime   float64
        }{
        Url: url,
        ContentLength: resp.ContentLength,
        BytesRead:  bytesRead,
        ElapsedTime: duration.Seconds(),
      } )
      results <- buf

  	}
  }
