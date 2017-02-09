package lazycache_benchmarking


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
