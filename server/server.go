package main

import (
	"fmt"
	"github.com/avicode/server/go_fullstack_app/file_model"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var Fs = file_model.FileDir{}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/files", getAllFilesHandler)
	router.GET("/files/:fileName", getFileByNameHandler)
	router.PUT("/files", uploadFile)
	router.DELETE("/files/:fileName", deleteFileByName)
	return router
}

func initServer() *gin.Engine {
	Fs = file_model.FileDir{DirPath: "./files/"}
	if err := Fs.MakeDirIfNotExists(); err != nil {
		fmt.Println(err)
	}
	return setupRouter()
}
func main() {
	r := initServer()
	r.Run(":3500")
}

func deleteFileByName(c *gin.Context) {
	fileName := c.Param("fileName")
	file, err := os.Open(filepath.Join(Fs.DirPath, fileName))
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintln("os.Open() Error: ", err.Error()))
		return
	}
	file.Close()
	if err = os.Remove(path.Join(Fs.DirPath, fileName)); err != nil {
		c.String(http.StatusNotFound, fmt.Sprintln("os.Remove() Error: ", err.Error()))
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("'%s' deleted!", fileName))
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintln("c.FormFile() Error: ", err.Error()))
		return
	}
	if file.Filename == "" {
		c.String(http.StatusNotFound, fmt.Sprintln("no file was find"))
		return
	}
	err = c.SaveUploadedFile(file, filepath.Join(Fs.DirPath, file.Filename))
	if err != nil {
		return
	}
}

func getFileByNameHandler(c *gin.Context) {
	fileName := c.Param("fileName")
	fmt.Println("fileName")
	fmt.Println(fileName)
	if fileName == "" {
		c.String(http.StatusNotFound, fmt.Sprintln("no file was find"))
		return
	}
	filePath := filepath.Join(Fs.DirPath, fileName)
	c.Status(http.StatusOK)
	c.FileAttachment(filePath, fileName)
}

func getAllFilesHandler(c *gin.Context) {
	files, err := Fs.GetAllFiles()
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintln("Error: ", err.Error()))
		return
	}
	c.JSON(http.StatusOK, files)
}
