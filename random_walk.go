package lazycache_benchmarking

import (
	"errors"
	"fmt"
	"io"
	//"io/ioutil"
	"encoding/json"
	"math/rand"
	"net/http"
)

func RandomWalk(opts StressOptions, baseurl string) error {

	count := opts.Count()
	parallelism := opts.Parallelism()

	if parallelism > count {
		fmt.Printf("Count %d less than parallelism %d, setting parallelism to %d\n", count, parallelism, count)
		parallelism = count
	}

	fmt.Printf("Random walk testing, count %d, parallelism %d\n", count, parallelism)

	var urls = make(chan string, parallelism)
	defer close(urls)

	var out = make(chan string)

	for i := 0; i < parallelism; i++ {
		go RandomWalkQuery(urls, out, baseurl)
		urls <- baseurl
	}

	//urls <- fmt.Sprintf("http://%s/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/", host )

	i := parallelism
	for {
		fmt.Printf("Wait for task %d to complete ...", i)
		new_url := <-out // wait for one task to complete

		// Always seed the channel with another url, just in case
		//urls <- fmt.Sprintf("http://%s/org/oceanobservatories/rawdata/files/",host)

		i++

		if i >= count {
			return nil
		} else if len(new_url) == 0 {
			return errors.New("Error from child")
		} else {
			urls <- new_url
		}
	}

}

func RandomWalkQuery(urls chan string, out chan string, baseurl string) {
	fmt.Println("In random walker")
	for url := range urls {

		fmt.Println("Random walker Querying URL", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%d: ERROR: %s\n", url, err)
			fmt.Printf("Error making request: %s\n", err.Error())
			out <- ""
			return

		}

		defer resp.Body.Close()

		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Printf("RESPONSE: %v\n%s\n", resp, body)

		// // Parse response
		decoder := json.NewDecoder(resp.Body)
		var listing struct {
			Path        string
			Files       []string
			Directories []string
		}

		//  var bar interface{}

		if err := decoder.Decode(&listing); err != nil && err != io.EOF {

			fmt.Printf("Query to %s; Error reading response (%d): %s\n", url, resp.StatusCode, err.Error())
			// body, _ := ioutil.ReadAll(resp.Body)
			// fmt.Printf("RESPONSE: %v\n%s\n", resp, body)
			out <- ""
			return
		}

		// fmt.Println(listing)
		// fmt.Printf("Has %d directories\n", len(listing.Directories))

		//    fmt.Println("Good response")

		if len(listing.Directories) > 0 {

			out <- url + listing.Directories[rand.Intn(len(listing.Directories))] + "/"
			//urls <- url + listing.Directories[rand.Intn(len(listing.Directories))]
		} else {
			out <- baseurl
		}
	}
}
