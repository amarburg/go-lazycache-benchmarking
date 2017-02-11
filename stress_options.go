package lazycache_benchmarking

import (
  "errors"
  "fmt"
  flag "github.com/spf13/pflag"
)


type StressSettings struct {
  urls    []string
  count, parallelism   int
}

type StressOption interface {
    Apply( *StressSettings )
}

type countSetter struct { count int }
func (c countSetter) Apply( settings *StressSettings ) {
  settings.count = c.count
}

func SetCount( count int ) StressOption {
  return countSetter{ count }
}

type parallelismSetter struct { parallel int }
func (c parallelismSetter) Apply( settings *StressSettings ) {
  settings.parallelism = c.parallel
}

func SetParallelism( par int ) StressOption {
  return parallelismSetter{ par }
}

type urlAdder struct { url string }
func (c urlAdder) Apply( settings *StressSettings ) {
   settings.urls = append(settings.urls, c.url)
}

func AddUrl( url string ) StressOption {
  return urlAdder{ url }
}


type urlsSetter struct { urls []string }
func (c urlsSetter) Apply( settings *StressSettings ) {
   settings.urls =   c.urls
}

func SetUrls( url []string ) StressOption {
  return urlsSetter{ url }
}

func NewSettings() *StressSettings {
  return &StressSettings {
  urls: make( []string, 0 ),
  count: 1,
  parallelism: 1,
  }
}

func (settings *StressSettings) Apply( opts ...StressOption ) error {
  for _,f := range opts {
    f.Apply(settings)
  }

  return settings.Validate()
}

func (settings *StressSettings) Validate() (error) {
  if len(settings.urls) < 1 {
    return errors.New("No valid image urls specified")
  }

  if settings.count < 1 {
    return errors.New("Count set to less than 1")
  }

  if settings.parallelism < 1 {
    return errors.New("Parallelism set to less than 1 ")
  }

return nil
}


func StressFlags(opts ...StressOption) (*StressSettings, error) {
  settings := NewSettings()
  settings.Apply( opts... )

  set := flag.NewFlagSet("", flag.ExitOnError)

  set.String("host", "127.0.0.1:5000", "Host to query")
    set.String("path", "/org/oceanobservatories/rawdata/files/", "Root path to query")

    set.Int("count", 1, "Number of queries to make")
    set.Int("parallelism", 5, "Parallelism of  queries")
set.Parse( flag.Args() )

  // if host,err := set.GetString( "host " ); err == nil {
  //   settings.Apply( SetHost( host ) )
  // }

  if count,err := set.GetInt( "count" ); err == nil {
    fmt.Println("Set count to ", count )
    settings.Apply( SetCount( count ) )
  }

  if par,err := set.GetInt( "parallelism" ); err == nil {
    settings.Apply( SetParallelism( par ) )
  }

  err := settings.Validate()
  return settings,err
}
