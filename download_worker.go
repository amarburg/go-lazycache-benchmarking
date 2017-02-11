package lazycache_benchmarking

import (
  "sync"
  "fmt"
  "io/ioutil"
  "net/http"
  "time"
  "encoding/json"
)

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
	 ioutil.ReadAll(resp.Body)
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
