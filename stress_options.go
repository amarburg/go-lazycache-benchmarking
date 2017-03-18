package lazycache_benchmarking

import (
	// "errors"
	// "fmt"
	flag "github.com/spf13/pflag"
)

type StressSettings interface {
	Count() int
	setCount(int)
	Parallelism() int
	setParallelism(int)
}

type StressOptions struct {
	count, parallelism int
}

func (opt StressOptions) Count() int       { return opt.count }
func (opt StressOptions) Parallelism() int { return opt.parallelism }

func (opt *StressOptions) SetCount(i int)       { opt.count = i }
func (opt *StressOptions) SetParallelism(i int) { opt.parallelism = i }

func NewSettings() *StressOptions {
	return &StressOptions{
		count:       1,
		parallelism: 1,
	}
}

// func (settings StressOptions) Apply( opts ...StressOption ) error {
//   for _,f := range opts {
//     f.Apply(settings)
//   }
//
//   return settings.Validate()
// }

// func (settings StressOptions) Validate() (error) {
//   if len(settings.urls) < 1 {
//     return errors.New("No valid image urls specified")
//   }
//
//   if settings.count < 1 {
//     return errors.New("Count set to less than 1")
//   }
//
//   if settings.parallelism < 1 {
//     return errors.New("Parallelism set to less than 1 ")
//   }
//
// return nil
// }

func AddStressFlags(set *flag.FlagSet) *flag.FlagSet {

	set.String("host", "127.0.0.1:5000", "Host to query")
	set.String("path", "/org/oceanobservatories/rawdata/files/", "Root path to query")

	set.Int("count", 1, "Number of queries to make")
	set.Int("parallelism", 5, "Parallelism of  queries")

	return set
}

func (opt *StressOptions) ApplyFlags(set *flag.FlagSet) {

	// if host,err := set.GetString( "host " ); err == nil {
	//   settings.Apply( SetHost( host ) )
	// }
	c, _ := set.GetInt("count")
	opt.SetCount(c)

	p, _ := set.GetInt("parallelism")
	opt.SetParallelism(p)

}
