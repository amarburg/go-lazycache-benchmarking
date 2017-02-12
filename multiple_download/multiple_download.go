package main

import (
	stress "github.com/amarburg/go-lazycache-benchmarking"
  "math/rand"
  "fmt"
  "os"
  flag "github.com/spf13/pflag"
)


type UrlList struct {
  urls     []string
}

func (list UrlList) Url( i int ) string {
  return list.urls[rand.Intn(len(list.urls))]
}

func main() {

  list := UrlList{
    []string{
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T000000Z.mov",
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T030000Z.mov",
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T060000Z.mov",
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T090000Z.mov",
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T120000Z.mov",
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T150000Z.mov",
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T180000Z.mov",
      "https://rawdata.oceanobservatories.org/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/01/CAMHDA301-20160701T210000Z.mov",
    },
  }

  opts := stress.NewSettings()

  set := flag.NewFlagSet("", flag.ExitOnError)
	stress.AddStressFlags(set)

  set.Parse( os.Args[1:])
  fmt.Println(set)
  opts.ApplyFlags(set)

	fmt.Println(opts)

	stress.RepeatedDownload(opts, list )
}
