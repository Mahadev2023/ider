package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/barasher/go-exiftool"
)

// declaring struct
type Student struct {

	// defining struct fields
	Name  string
	Marks int
	Id    string
}

func getPath() string {
	Cwd, ee := os.Getwd()
	if ee != nil {
		fmt.Println(ee)
	}
	return Cwd

}

// // defining struct instance
// std1 := Student{"Vani", 94, "20024"}
// // Parsing the required html
// // file in same directory
// t, err := template.ParseFiles("index.html")

// // standard output to print merged data
// err = t.Execute(os.Stdout, std1)

// http.Handle("/", http.FileServer(http.Dir("./static")))
// http.ListenAndServe(":3000", nil)

func main() {
	// fs := http.FileServer(http.Dir("./public"))
	// http.Handle("/image/", http.StripPrefix("/image", fs))
	// http.HandleFunc("/", uploadHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/image/", http.StripPrefix("/image", fs))
	http.HandleFunc("/", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		renderForm(w)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read the image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the image data", http.StatusInternalServerError)
		return
	}

	// Perform any image processing or property extraction here
	imageName := header.Filename
	imageProperty := "Example Property"

	// Store the image file on the server
	// err = ioutil.WriteFile(header.Filename, fileBytes, 0644)

	// Pathimg := path.Join(getPath(), "public", "image", header.Filename)

	TempName := "temp" + filepath.Ext(header.Filename)
	Pathimgtemp := path.Join(getPath(), "public", "image", TempName)

	err = ioutil.WriteFile(header.Filename, fileBytes, 0644)
	err = ioutil.WriteFile(Pathimgtemp, fileBytes, 0644)

	err = ioutil.WriteFile(header.Filename, fileBytes, 0644)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("problem cwd")
	}

	exiftoolPath := path.Join(cwd, "exiftool.exe")

	et, err := exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(exiftoolPath))

	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(header.Filename)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			fmt.Printf("[%v] %v\n", k, v)
		}
	}

	if err != nil {
		http.Error(w, "Failed to store the image file", http.StatusInternalServerError)
		return
	}
	// Render the HTML response with the image
	renderHTML(w, imageName, imageProperty, header.Filename, TempName)

}

func renderForm(w http.ResponseWriter) {
	html := `<html>
		<body>
			<form action="/" method="post" enctype="multipart/form-data">
				<input type="file" name="image" />
				<input type="submit" value="Upload" />
			</form>
		</body>
	</html>`

	fmt.Fprint(w, html)
}

type ImageInfo struct {
	ImageName     string
	ImageProperty string
	Filename      string
}

func renderHTML(w http.ResponseWriter, imageName, imageProperty, filename string, TempName string) {
	htmlTemplate := `<html>
		<body>
			<h1>Image Information:</h1>
			<p>Image Name: {{.ImageName}}</p>
			<p>Image Property: {{.ImageProperty}}</p>
			<img src="/image/{{.Filename}}" alt="Uploaded Image" />
		</body>
	</html>`

	imgSrc := "/image/" + TempName

	tmpl, err := template.New("imageInfo").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Failed to render the HTML template", http.StatusInternalServerError)
		return
	}

	data := ImageInfo{
		ImageName:     imageName,
		ImageProperty: imageProperty,
		Filename:      imgSrc,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render the HTML template", http.StatusInternalServerError)
		return
	}
	err = os.Remove(filename)
	if err != nil {
		log.Printf("Failed to delete the image file: %v\n", err)
	}

}

// 	// Delete the image file from the server

// 	// err = os.Remove(filename)
// 	// if err != nil {
// 	// 	log.Printf("Failed to delete the image file: %v\n", err)
// 	// }
//
