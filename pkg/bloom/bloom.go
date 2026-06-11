package bloom

import (
	bloomfilter "github.com/bits-and-blooms/bloom/v3"
)

type Filter struct {
	bf *bloomfilter.BloomFilter
}

func New(capacity uint, fpRate float64) *Filter {
	return &Filter{
		bf: bloomfilter.NewWithEstimates(capacity, fpRate),
	}
}

func (f *Filter) Add(key string) {
	f.bf.Add([]byte(key))
}

func (f *Filter) AddBatch(keys []string) {
	for _, k := range keys {
		f.bf.Add([]byte(k))
	}
}

func (f *Filter) Contains(key string) bool {
	return f.bf.Test([]byte(key))
}
