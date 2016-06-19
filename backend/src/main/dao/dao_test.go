package dao

import (
	"testing"
)

func BenchmarkGetStreetName(b *testing.B) {

	b.StopTimer()
	Open()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_ = GetStreetName(-34.501812, 150.810264)
	}
	b.StopTimer()
	Close()
	b.StartTimer()

}
