// SwartzNotes standalone package for clients
// ------------------------------------------
// Written by Huan Truong <htruong@tnhh.net>
// Licensed under an Apache license

// The package provides minimal support for the browser extension to work
// Also serves as a bridge for the fact extractor java package to work

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
	//	"strings"
)

type FileUploadResponse struct {
	ID      string
	Result  int
	Message string
}

var randGen *rand.Rand
var workspaceDir string

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// Grinds the PDF that the extension sent and sends back the id
func grind(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		randID := fmt.Sprintf("%d", randGen.Int63())
		storageLocation := workspaceDir + randID + ".pdf"
		f, err := os.OpenFile(storageLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}

		io.Copy(f, file)
		f.Close()

		log.Printf("Executing java -jar factExtractor.jar -o %s -pdf %s", workspaceDir+fmt.Sprintf("%s.fact.json", randID), storageLocation)
		cmd := exec.Command("java", "-jar", "factExtractor.jar", "-pdf", storageLocation, "-o", workspaceDir+fmt.Sprintf("%s.fact.json", randID))

		err = cmd.Run()

		msg := FileUploadResponse{ID: randID, Result: 0, Message: "File uploaded successfully"}

		if err != nil {
			//log.Fatal(err)
			fmt.Println("Fact extractor crashed on ", randID)
			msg.Result = -1
			msg.Message = "Fact extractor crashed."
		}

		b, _ := json.Marshal(msg)

		w.Write(b)
	}
}

func main() {
	_, err := exec.LookPath("java")
	if err != nil {
		log.Fatal("We couldn't find an executable of Java, please make sure you have Java installed. Java can be downloaded at http://java.com/download")
	}

	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

	workspaceDir = UserHomeDir() + "/factsWorkspace/"

	os.MkdirAll(workspaceDir, os.ModeDir|0777)

	fs := http.FileServer(http.Dir(workspaceDir))
	log.Printf("Workspace directory is %s", workspaceDir)

	statics := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/grind", grind)
	http.Handle("/get/", http.StripPrefix("/get", fs))
	http.Handle("/static/", http.StripPrefix("/static", statics))

	err = http.ListenAndServe(":3333", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
