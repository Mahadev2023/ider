package main

import (
	"fmt"
	"os"
	"path"

	"github.com/barasher/go-exiftool"
)

func main() {
	fmt.Println("Hello")
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("problem cwd")
	}

	exiftoolPath := path.Join(cwd, "exiftool.exe")

	fmt.Println(exiftoolPath)
	//if set environement variable
	// et, err := exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath("C:/Program Files/Exiftool/exiftool.exe"))

	// et, err := exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(cwd + "\\exiftool.exe"))

	et, err := exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(exiftoolPath))

	// C:\Users\godbo\OneDrive\Desktop\project\exiftool
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(cwd, "bg.jpg")

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			fmt.Printf("[%v] %v\n", k, v)
		}
	}
}
