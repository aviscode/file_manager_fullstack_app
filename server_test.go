package main

import (
	"encoding/json"
	"github.com/avicode/go_fullstack_app/file_model"
	"github.com/avicode/go_fullstack_app/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func Test_deleteFileByName(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		url      string
		expected interface{}
	}{
		{"Test good delete", "server.go", "/files", 200},
		{"Test non existing file delete", "server1.go", "/files", 404},
		{"Test bad delete without a file", "", "/files", 404},
		{"Test bad delete without a a url", "server.go", "", 404},
	}

	router := initServer()
	w := httptest.NewRecorder()
	req := utils.NewFileUploadRequest("server.go", "file", "/files")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("failed to put file for the test")
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", path.Join(tt.url, tt.fileName), nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expected, w.Code)
		})
	}
}

func Test_uploadFile(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		fieldName string
		url       string
		expected  interface{}
	}{
		{"Test upload good", "server.go", "file", "/files", 200},
		{"Test upload good url without a file", "", "file", "/files", 404},
		{"Test upload bad field name", "server.go", "image", "/files", 404},
		{"Test upload bad url", "server.go", "file", "/", 404},
	}
	router := initServer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := utils.NewFileUploadRequest(tt.fileName, tt.fieldName, tt.url)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expected, w.Code)
		})
	}
}

func Test_getFileByNameHandler(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		url      string
		expected interface{}
	}{
		{"Test get file good", "server.go", "/files", 200},
		{"Test get a non existing file", "server1.go", "/files", 404},
		{"Test bad get url", "server.go", "/", 404},
	}
	router := initServer()
	w := httptest.NewRecorder()
	req := utils.NewFileUploadRequest("server.go", "file", "/files")
	router.ServeHTTP(w, req)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", path.Join(tt.url, tt.fileName), nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expected, w.Code)
		})
	}
}

func Test_getFilesHandler(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		expected         interface{}
		expectedNumFiles interface{}
	}{
		{"Test get file good with 2 files in dir", "/files", 200, 2},
		{"Test get a files when non file in dir", "/files", 200, 0},
	}
	router := initServer()
	w := httptest.NewRecorder()
	var req = []*http.Request{
		utils.NewFileUploadRequest("server.go", "file", "/files"),
		utils.NewFileUploadRequest("server_test.go", "file", "/files"),
	}
	for _, i := range req {
		router.ServeHTTP(w, i)
	}

	for i, tt := range tests {
		if i == 1 {
			err := Fs.DeleteAllFileInDir()
			if err != nil {
				t.Errorf(err.Error())
				return
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expected, w.Code)
			files := &[]file_model.File{}
			err := json.NewDecoder(w.Body).Decode(files)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			assert.Equal(t, tt.expected, http.StatusOK)
			assert.Equal(t, tt.expectedNumFiles, len(*files))
		})
	}
}
