package lazycache_benchmarking

import (
  "fmt"
  "net/http"
  "encoding/json"
  "math/rand"
  "sync"
)





func FrameStress( opts ...StressOption ) error {
  var wg sync.WaitGroup

  settings := NewSettings()
  if err := settings.Apply( opts... ); err != nil { return nil }

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
  var results = make( chan []byte, settings.count )

  wg.Add( settings.parallelism )

	for i := 0; i < settings.parallelism; i++ {
		go DownloadWorker(urls, results, &wg)
	}

  for i := 0; i < settings.count; i ++ {
    c := rand.Intn( len(settings.urls) )
    urls <- fmt.Sprintf("%s/frame/%d", settings.urls[c], rand.Intn(  lengths[c] ) )
  }

  // Drop results

  close(urls)
  fmt.Printf("Waiting for workers to finish...")
  wg.Wait()
  fmt.Printf("done\n")

  return nil
}
