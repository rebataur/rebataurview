package main

import (
	"fmt"
	"strings"
	"net/http"
	"log"
)

var nwPath string

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cmd",cmdHandler)
	http.Handle("/app/", http.FileServer(http.Dir(nwPath)))
	http.ListenAndServe(":9999", nil)

}
func homeHandler(w http.ResponseWriter, r *http.Request){
	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	fmt.Fprintf(w, "Hi there, I love %s!", "hahah")

	// executeCommand("initpg D:\\uploads\\cc.csv")

}
func cmdHandler(w http.ResponseWriter, r *http.Request){
	cmd := r.FormValue("cmd")
	executeCommand(cmd)
	fmt.Fprintf(w, "%s", result )

}

func executeCommand(cmd string) {
	rootCmd.SetArgs(strings.Split(cmd," "))
	rootCmd.Execute()
}

func init(){
	config,err := getConfig()
	if err == nil{
		nwPath = config.NW.NWPath
	}else{
		log.Fatal("Error getting NW Path")
	}

}
