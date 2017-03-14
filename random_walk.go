package lazycache_benchmarking

import (
  "fmt"
  "errors"
  "net/http"
  // "encoding/json"
  // "math/rand"
)

func RandomWalk( opts StressOptions, baseurl string ) error {

  count := opts.Count()
  parallelism := opts.Parallelism()

  if parallelism > count {
    fmt.Printf("Count %d less than parallelism %d, setting parallelism to %d\n", count, parallelism, count )
    parallelism = count
  }

  fmt.Printf("Random walk testing, count %d, parallelism %d\n", count, parallelism )

  var urls = make(chan string, count )
  var out = make(chan bool)

	for i := 0; i < parallelism; i++ {
		go RandomWalkQuery(urls,out, baseurl)
    urls <- baseurl
	}


	//urls <- fmt.Sprintf("http://%s/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/", host )

	i := 1
	for {
		fmt.Printf("Wait for task %d to complete ...", i)
		resp := <-out // wait for one task to complete

		// Always seed the channel with another url, just in case
		//urls <- fmt.Sprintf("http://%s/org/oceanobservatories/rawdata/files/",host)

		i++

		if !resp {
			return errors.New("Error from child")
		} else if i >= count {
			return nil
		}
	}

}

func RandomWalkQuery(urls chan string, out chan bool, baseurl string) {
	fmt.Println("In random walker")
	for url := range urls {

		fmt.Println("Random walker Querying URL", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%d: ERROR: %s\n", url, err)
			fmt.Printf("Error making request: %s\n", err.Error())
			out <- false
			return

		}

		defer resp.Body.Close()
		// body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Printf("%d: RESPONSE: %v\n%s\n", i, resp, body)

    // TODO:  Fix this
		// // Parse response
		// decoder := json.NewDecoder(resp.Body)
		// var listing lazycache.DirListing
    //
		// if err := decoder.Decode(&listing); err != nil {
		// 	fmt.Println("Error reading response: %s\n", err.Error())
		// 	out <- false
		// 	return
		// }
		// //fmt.Printf("Has %d directories\n", len(listing.Directories))
    //
		// if len(listing.Directories) > 0 {
    //
		// 	urls <- url + listing.Directories[rand.Intn(len(listing.Directories))]
		// 	//urls <- url + listing.Directories[rand.Intn(len(listing.Directories))]
		// } else {
    //   urls <- baseurl
    // }

		fmt.Println("Good response")
		out <- true
	}
}
