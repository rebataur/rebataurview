package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
)

import (
	"github.com/ranjanprj/rebataurview/cmds"
)

var nwPath,repositoryPath,uploadedFilePath string

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cmd", cmdHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/app/", http.FileServer(http.Dir(nwPath)))
	http.ListenAndServe(":9999", nil)

}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	fmt.Fprintf(w, "This is home")

	// executeCommand("initpg D:\\uploads\\cc.csv")

}
func cmdHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	executeCommand(cmd)
	fmt.Fprintf(w, "%s", cmds.Result)

}



func executeCommand(cmd string) {
	cmds.SetAndExecuteCmd(cmd)

}


func uploadHandler(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	if action == "upload" {

		var fileName string

		r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
		fhs := r.MultipartForm.File["file"]
		for _, fh := range fhs {
			fileName = fh.Filename

			f, err := fh.Open()
			if err == nil {
				if file, err := ioutil.ReadAll(f); err == nil {
					path := strings.Join([]string{repositoryPath,fh.Filename}, "")
					ioutil.WriteFile(path, file, 0644)
					// fmt.Fprintf(w, "Done")
					fmt.Println("Writing to file done")
					//http.Redirect(w, r, "/static/dojoui/rep.html", http.StatusFound)
					//fmt.Fprintf(w, "File uploaded to repository",http.StatusFound)

				} else {
					fmt.Fprintf(w, "ERROR IN : File uploaded to repository",http.StatusNotFound)
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}

		uploadedFilePath = strings.Join([]string{repositoryPath, fileName}, "")

		fmt.Println("PG COPY**************8",uploadedFilePath)
		cmds.LoadDataIntoPG(uploadedFilePath,true)
	}

}

func init() {
		fmt.Println("init")
	config, err := cmds.GetConfig()
	fmt.Println("getting config")
	if err == nil {
		nwPath = config.NW.NWPath
		repositoryPath = config.Repository.Path
		fmt.Println(repositoryPath)
	} else {
		log.Fatal("Error getting NW Path")
	}

}
