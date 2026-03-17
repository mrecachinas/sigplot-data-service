package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/spectriclabs/sigplot-data-service/internal/config"
	"github.com/stretchr/testify/assert"
)

var sdsConfigString string = `[{"location_name":"ServiceDir","location_type":"localFile","path":"./"},{"location_name":"ServiceDirData","location_type":"localFile","path":"./data"},{"location_name":"sdsdata","location_type":"localFile","path":"/data/sdsdata/"},{"location_name":"minio","location_type":"minio","minio_bucket":"sdsdata","location":"192.168.1.229:9000","minio_access_key":"minio","minio_secret_key":"miniostorage"}]`

func TestFS(t *testing.T) {
	var locationDetails []config.Location
	err := json.Unmarshal([]byte(sdsConfigString), &locationDetails)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}

	sdsConfig := config.Config{
		LocationDetails: locationDetails,
	}

	e := echo.New()
	req, err := http.NewRequest(http.MethodGet, "/sds/fs", nil)
	if err != nil {
		t.Fatalf("The request could not be created because of: %v", err)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sds/fs")
	a := NewSDSAPI(&sdsConfig)

	if assert.NoError(t, a.GetFileLocations(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, sdsConfigString, strings.TrimSpace(rec.Body.String()))
	}
}

func TestFSDir(t *testing.T) {
	sdsConfig := config.Config{
		LocationDetails: []config.Location{
			{LocationName: "TestDir", LocationType: "localFile", Path: "."},
		},
	}

	e := echo.New()
	req, err := http.NewRequest(http.MethodGet, "/sds/fs/TestDir/", nil)
	if err != nil {
		t.Fatalf("The request could not be created: %v", err)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sds/fs/:location/*")
	c.SetParamNames("location", "*")
	c.SetParamValues("TestDir", "")
	a := NewSDSAPI(&sdsConfig)

	if assert.NoError(t, a.GetFileOrDirectory(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "filename")
	}
}

func TestFSFile(t *testing.T) {
	sdsConfig := config.Config{
		LocationDetails: []config.Location{
			{LocationName: "TestDir", LocationType: "localFile", Path: "."},
		},
	}

	e := echo.New()
	req, err := http.NewRequest(http.MethodGet, "/sds/fs/TestDir/fs_test.go", nil)
	if err != nil {
		t.Fatalf("The request could not be created: %v", err)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sds/fs/:location/*")
	c.SetParamNames("location", "*")
	c.SetParamValues("TestDir", "fs_test.go")
	a := NewSDSAPI(&sdsConfig)

	if assert.NoError(t, a.GetFileOrDirectory(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, rec.Body.Len() > 0, "File contents should not be empty")
	}
}

func TestFSMinio(t *testing.T) {
	sdsConfig := config.Config{
		LocationDetails: []config.Location{
			{LocationName: "minio", LocationType: "minio", MinioBucket: "sdsdata"},
		},
	}

	e := echo.New()
	req, err := http.NewRequest(http.MethodGet, "/sds/fs/minio/", nil)
	if err != nil {
		t.Fatalf("The request could not be created: %v", err)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sds/fs/:location/*")
	c.SetParamNames("location", "*")
	c.SetParamValues("minio", "")
	a := NewSDSAPI(&sdsConfig)

	a.GetFileOrDirectory(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
