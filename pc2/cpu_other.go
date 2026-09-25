//go:build !windows

package main

import "runtime/metrics"

func processCPUSeconds() float64 {
	s := []metrics.Sample{{Name: "/cpu/classes/total:cpu-seconds"}}
	metrics.Read(s)
	if s[0].Value.Kind() == metrics.KindFloat64 {
		return s[0].Value.Float64()
	}
	return 0
}
