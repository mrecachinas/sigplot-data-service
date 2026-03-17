package cache

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkUrlToCacheFileName(b *testing.B) {
	urls := []struct {
		name string
		url  string
	}{
		{"short", "data/file.tmp"},
		{"with_params", "data/file.tmp?x=1&y=2&z=3"},
		{"long_path", "/api/v1/locations/local/files/very/deep/nested/path/to/data.blue?subsize=1024&outfmt=RGBA&colormap=Ramp%20Colormap"},
	}

	for _, tc := range urls {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				UrlToCacheFileName(tc.url)
			}
		})
	}
}

func BenchmarkPutGetRoundTrip(b *testing.B) {
	sizes := []struct {
		name  string
		bytes int
	}{
		{"1KB", 1024},
		{"100KB", 100 * 1024},
		{"1MB", 1024 * 1024},
	}

	for _, s := range sizes {
		tmpDir, err := os.MkdirTemp("", "cache_bench_*")
		if err != nil {
			b.Fatal(err)
		}

		c := &Cache{Location: tmpDir}
		data := make([]byte, s.bytes)
		for i := range data {
			data[i] = byte(i & 0xFF)
		}
		fileName := fmt.Sprintf("bench_%s", s.name)

		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				c.PutItemInCache(fileName, "bench/", data)
				c.GetDataFromCache(fileName, "bench/")
			}
		})

		os.RemoveAll(tmpDir)
	}
}
