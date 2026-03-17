package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestUrlToCacheFileName(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "simple URL no query",
			url:  "data/file.tmp",
			want: "datafiletmp",
		},
		{
			name: "URL with query params",
			url:  "data/file.tmp?x=1&y=2",
			want: "datafiletmp_x1y2",
		},
		{
			name: "URL with dots and slashes",
			url:  "path/to/some.file.dat",
			want: "pathtosomefiledat",
		},
		{
			name: "URL with multiple special chars",
			url:  "a/b/c.d?key=val&foo=bar",
			want: "abcd_keyvalfoobar",
		},
		{
			name: "empty URL",
			url:  "",
			want: "",
		},
		{
			name: "URL with no slashes or dots",
			url:  "simplename",
			want: "simplename",
		},
		{
			name: "URL with only question mark",
			url:  "file?",
			want: "file_",
		},
		{
			name: "URL with equals and ampersand",
			url:  "f?a=1&b=2&c=3",
			want: "f_a1b2c3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UrlToCacheFileName(tt.url)
			if got != tt.want {
				t.Errorf("UrlToCacheFileName(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestPutAndGetDataFromCache(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := &Cache{Location: tmpDir}
	subDir := "/"
	fileName := "testfile"
	data := []byte("hello cache data")

	if err := c.PutItemInCache(fileName, subDir, data); err != nil {
		t.Fatalf("PutItemInCache error: %v", err)
	}

	got, err := c.GetDataFromCache(fileName, subDir)
	if err != nil {
		t.Fatalf("GetDataFromCache error: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("GetDataFromCache = %q, want %q", string(got), string(data))
	}
}

func TestPutAndGetItemFromCache(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := &Cache{Location: tmpDir}
	subDir := "/"
	fileName := "testitem"
	data := []byte("seekable cache data")

	if err := c.PutItemInCache(fileName, subDir, data); err != nil {
		t.Fatalf("PutItemInCache error: %v", err)
	}

	reader, err := c.GetItemFromCache(fileName, subDir)
	if err != nil {
		t.Fatalf("GetItemFromCache error: %v", err)
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if string(content) != string(data) {
		t.Errorf("GetItemFromCache content = %q, want %q", string(content), string(data))
	}

	// Close the file if possible
	if closer, ok := reader.(interface{ Close() error }); ok {
		closer.Close()
	}
}

func TestGetDataFromCacheMissingFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := &Cache{Location: tmpDir}
	_, err = c.GetDataFromCache("nonexistent", "/")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestGetItemFromCacheMissingFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := &Cache{Location: tmpDir}
	_, err = c.GetItemFromCache("nonexistent", "/")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestPutItemInCacheCreatesSubDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := &Cache{Location: tmpDir}
	subDir := "/newsubdir/"
	fileName := "testfile"
	data := []byte("data in subdir")

	if err := c.PutItemInCache(fileName, subDir, data); err != nil {
		t.Fatalf("PutItemInCache error: %v", err)
	}

	// Verify the subdirectory was created
	info, err := os.Stat(filepath.Join(tmpDir, "newsubdir"))
	if err != nil {
		t.Fatalf("subdirectory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory, got file")
	}
}
