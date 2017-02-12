package main

import (
	flag "github.com/spf13/pflag"
	stress "github.com/amarburg/go-lazycache-benchmarking"
	"os"
	"math/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RandomVideoFrames struct {
	url       string
	numFrames int
}

func NewRandomVideoFrames(url string) (RandomVideoFrames, error) {
	framer := RandomVideoFrames{
		url:       url,
		numFrames: -1,
	}

fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return framer, err
	}
	defer resp.Body.Close()
	imgInfo := json.NewDecoder(resp.Body)

	var videoInfo struct {
		NumFrames int
	}
	err = imgInfo.Decode(&videoInfo)

	if err != nil {
		panic(fmt.Sprintf("Couldn't figure out number of frames: %s", err.Error()))
	} else if videoInfo.NumFrames < 1 {
		panic("Couldn't figure out number of frames")
	}
	fmt.Printf("Video has %d frames\n", videoInfo.NumFrames)

	framer.numFrames = videoInfo.NumFrames

	return framer, nil
}

func (video RandomVideoFrames) Url(i int) string {
	frameNum := rand.Intn(video.numFrames)
	return fmt.Sprintf("%s/frame/%d", video.url, frameNum)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	url := "%h/files/RS03ASHS/PN03B/06-CAMHDA301/2017/01/01/CAMHDA301-20170101T000500.mov"
	host := "https://ferrous-ranger-158304.appspot.com/org/oceanobservatories/rawdata"

	fullUrl := strings.Replace(url, "%h", host, 1)

	opts := stress.NewSettings()

	set := flag.NewFlagSet("", flag.ExitOnError)
	stress.AddStressFlags(set)

	set.Parse(os.Args[1:])
	fmt.Println(set)
	opts.ApplyFlags(set)

	fmt.Println(opts)

	framer,err := NewRandomVideoFrames(fullUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	stress.RepeatedDownload(opts, framer)
}
