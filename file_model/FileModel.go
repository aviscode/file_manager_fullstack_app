package file_model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type File struct {
	FileName   string `json:"filename"`
	Format     string `json:"format"`
	Size       int64  `json:"size"`
	UploadDate string `json:"uploaded"`
}

type FileDir struct {
	DirPath string `json:"dir_path"`
}

func (fs *FileDir) GetAllFiles() ([]File, error) {
	files, err := ioutil.ReadDir(fs.DirPath)
	if err != nil {
		return nil, err
	}
	var data []File
	for _, file := range files {
		format := strings.Split(file.Name(), ".")
		data = append(data, File{
			FileName:   file.Name(),
			Format:     format[len(format)-1],
			Size:       file.Size(),
			UploadDate: file.ModTime().Format("2006-01-02"),
		})
	}
	return data, nil
}

func (fs *FileDir) MakeDirIfNotExists() error {
	if _, err := os.Stat(fs.DirPath); os.IsNotExist(err) {
		return os.Mkdir(fs.DirPath, os.ModeDir|0755)
	}
	return nil
}

func (fs *FileDir) DeleteAllFileInDir() error {
	files, err := ioutil.ReadDir(fs.DirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := os.Remove(path.Join(fs.DirPath, file.Name()))
		if err != nil {
			return err
		}
	}
	if files, err = ioutil.ReadDir(fs.DirPath); err != nil {
		return err
	}
	if len(files) != 0 {
		return fmt.Errorf("Not all files ware deleted ")
	}
	return nil
}
