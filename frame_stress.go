package lazycache_benchmarking

import (
  "fmt"
  "net/http"
  "encoding/json"
  "math/rand"
  "sync"
  "errors"
  "io/ioutil"
)

var wg sync.WaitGroup




func FrameStress( opts ...StressOption ) error {

  settings := &stressSettings {
  urls: make( []string, 0 ),
  count: 0,
  parallelism: 0,
  }

  for _,f := range opts {
    f.Apply(settings)
  }

  if len(settings.urls) < 1 {
    return errors.New("No valid image urls specified")
  }


  lengths := make( []int, len(settings.urls) )

    for i,imgUrl := range( settings.urls ) {
    resp,err := http.Get( imgUrl )
    if err != nil {
      return err
    }
    defer resp.Body.Close()
    imgInfo := json.NewDecoder( resp.Body )

    var videoInfo struct {
      NumFrames   int
    }
    err = imgInfo.Decode( &videoInfo )

    if err != nil {
      panic( fmt.Sprintf("Couldn't figure out number of frames: %s", err.Error()))
    } else if videoInfo.NumFrames < 1  {
      panic("Couldn't figure out number of frames")
    }
    fmt.Println("Video has %d frames", videoInfo.NumFrames)

    lengths[i] = videoInfo.NumFrames
  }

  var urls = make(chan string )

  wg.Add( settings.parallelism )

	for i := 0; i < settings.parallelism; i++ {
		go FrameStressWorker(urls)
	}

  for i := 0; i < settings.count; i ++ {
    c := rand.Intn( len(settings.urls) )
    urls <- fmt.Sprintf("%s/frame/%d", settings.urls[c], rand.Intn(  lengths[c] ) )
  }

  close(urls)
  fmt.Printf("Waiting for workers to finish...")
  wg.Wait()
  fmt.Printf("done\n")

  return nil
}

func FrameStressWorker(urls chan string) {
	fmt.Println("In random walker")
	for {

    url,ok := <- urls

    if !ok {
      //fmt.Println("Channel closed, quitting")
      wg.Done()
      return
    }

		fmt.Println("Random walker Querying URL", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%d: ERROR: %s\n", url, err)
			fmt.Printf("Error making request: %s\n", err.Error())
		}

    //
		defer resp.Body.Close()
	 ioutil.ReadAll(resp.Body)


	}

}
