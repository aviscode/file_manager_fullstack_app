package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func NewFileUploadRequest(fileName, fieldName, url string) *http.Request {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		fmt.Errorf("writer.CreateFormFile() error = %v", err)
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	req, _ := http.NewRequest(http.MethodPut, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
