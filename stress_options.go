package lazycache_benchmarking

import (
  "errors"
)


type stressSettings struct {
  urls    []string
  count, parallelism   int
}

type StressOption interface {
  Apply( *stressSettings )
}

type countSetter struct { count int }
func (c countSetter) Apply( settings *stressSettings ) {
  settings.count = c.count
}

func SetCount( count int ) StressOption {
  return countSetter{ count }
}

type parallelismSetter struct { parallel int }
func (c parallelismSetter) Apply( settings *stressSettings ) {
  settings.parallelism = c.parallel
}

func SetParallelism( par int ) StressOption {
  return parallelismSetter{ par }
}

type urlAdder struct { url string }
func (c urlAdder) Apply( settings *stressSettings ) {
   settings.urls = append(settings.urls, c.url)
}

func AddUrl( url string ) StressOption {
  return urlAdder{ url }
}


func NewSettings() *stressSettings {
  return &stressSettings {
  urls: make( []string, 0 ),
  count: 1,
  parallelism: 1,
  }
}

func (settings *stressSettings) Apply( opts ...StressOption ) error {
  for _,f := range opts {
    f.Apply(settings)
  }

  return settings.Validate()
}

func (settings *stressSettings) Validate() (error) {
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
