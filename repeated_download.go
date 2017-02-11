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

    for i := 0; i < count; i ++ {
      urls <- producer.Url(i)
    }

    close(urls)
    fmt.Printf("Waiting for workers to finish...\n")
    wg.Wait()

    fmt.Printf("Pulling results\n")
    for i := 0; i < count; i ++ {
      result := <- results
      fmt.Println(bytes.NewBuffer(result).String())
    }

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

      //
  		defer resp.Body.Close()

      mb := make( []byte, 1024*1024 )
      for {
        if _,err = resp.Body.Read(mb); err != nil { break }
      }
      duration := time.Since(start)

      // fmt.Println("Done fetching ", url, " elapsed ", duration )

      buf,_ := json.Marshal( struct {
        Url string
        ContentLength int64
        ElapsedTime   float64
        }{
        url,
        resp.ContentLength,
        duration.Seconds(),
      } )
      results <- buf

  	}
  }
