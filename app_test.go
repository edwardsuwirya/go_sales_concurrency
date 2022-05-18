package main

import "testing"

// to run
// go test -bench . -count 5 -benchtime=15000x
func BenchmarkSerialCalculation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		serialCalculation()
	}
}

func BenchmarkConcurrentCalculation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concurrentCalculation()
	}
}
func BenchmarkChannelCalculation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		channelCalculation()
	}
}
