package lazycache_benchmarking


import (
  "sync"
  "fmt"
  "math/rand"
  "bytes"
)


func RepeatedDownload( settings *StressSettings ) error {

  if settings.parallelism > settings.count {
    fmt.Printf("Count %d less than parallelism %d, setting parallelism to %d\n", settings.count, settings.parallelism, settings.count )
    settings.parallelism = settings.count
  }

  var wg sync.WaitGroup

    var urls = make(chan string)
    var results = make(chan []byte, settings.count )

    wg.Add( settings.parallelism )

  	for i := 0; i < settings.parallelism; i++ {
  		go DownloadWorker(urls, results, &wg)
  	}

    for i := 0; i < settings.count; i ++ {
      c := rand.Intn( len(settings.urls) )
      urls <- settings.urls[c]
    }

    close(urls)
    fmt.Printf("Waiting for workers to finish...\n")
    wg.Wait()

    fmt.Printf("Pulling results\n")
    for i := 0; i < settings.count; i ++ {
      result := <- results
      fmt.Println(bytes.NewBuffer(result).String())
    }



    return nil
  }
