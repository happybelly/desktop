// SwartzNotes standalone package for clients
// ------------------------------------------
// Written by Huan Truong <htruong@tnhh.net>
// Licensed under an Apache license

// The package provides minimal support for the browser extension to work
// Also serves as a bridge for the fact extractor java package to work

package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ResponseStruct struct {
	ID      string
	Result  int
	Message string
}

var randGen *rand.Rand
var workspaceDir string

// Returns the user home directory
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

// Downloads a file from a location on the net
func HTTPDownload(uri string) ([]byte, error) {
	fmt.Printf("HTTPDownload From: %s.\n", uri)

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	res, err := client.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ReadFile: Size of download: %d\n", len(d))
	return d, err
}

// Write to file dst with array of bytes
func WriteFile(dst string, d []byte) error {
	fmt.Printf("WriteFile: Size of download: %d\n", len(d))
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Download from url uri to save at dst
func DownloadToFile(uri string, dst string) {
	log.Printf("DownloadToFile From: %s.\n", uri)
	if d, err := HTTPDownload(uri); err == nil {
		log.Printf("Downloaded %s.\n", uri)
		if WriteFile(dst, d) == nil {
			log.Printf("Saved %s as %s\n", uri, dst)
		}
	}
}

// Grinds the PDF that has a fixed URL and sends back the id
func onlineGrind(w http.ResponseWriter, r *http.Request) {
	downloadURL := r.URL.Query().Get("url")

	h := md5.New()
	io.WriteString(h, downloadURL)
	randID := fmt.Sprintf("%x", h.Sum(nil))

	storageLocation := workspaceDir + randID + ".pdf"

	if _, err := os.Stat(storageLocation); os.IsNotExist(err) {

		DownloadToFile(downloadURL, storageLocation)
		ProcessPDFFile(storageLocation, randID, w, r)

	} else {
		http.Redirect(w, r, fmt.Sprintf("/get/%s.fact.json", randID), 303)
	}
}

// Given a file and its ID, call FactExtractor to extract the fact and send back the json file
func ProcessPDFFile(storageLocation string, fileID string, w http.ResponseWriter, r *http.Request) {
	log.Printf("Executing java -jar factExtractor.jar -o %s -pdf %s", workspaceDir+fmt.Sprintf("%s.fact.json", fileID), storageLocation)
	cmd := exec.Command("java", "-jar", "factExtractor.jar", "-pdf", storageLocation, "-o", workspaceDir+fmt.Sprintf("%s.fact.json", fileID))

	err := cmd.Run()
	log.Printf("File is ground.")
	msg := ResponseStruct{ID: fileID, Result: 0, Message: "File ground successfully."}

	if err != nil {
		//log.Fatal(err)
		fmt.Println("Fact extractor crashed on ", fileID)
		msg.Result = -1
		msg.Message = "Fact extractor crashed."

		b, _ := json.Marshal(msg)

		w.Write(b)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/get/%s.fact.json", fileID), 303)
	}

}

// Saves the PDF blob that is sent by the extension
func processPDFBlob(w http.ResponseWriter, r *http.Request) {

	rURI, _ := url.Parse(r.RequestURI)
	fileID := rURI.Query().Get("uri")

	storageLocation := workspaceDir + fileID + ".pdf"

	// We sniff if the pdf file is already there, if it were there previously,
	// the facts file has to be there too
	if _, err := os.Stat(storageLocation); os.IsNotExist(err) {
		file, _, err := r.FormFile("data")

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer file.Close()

		out, err := os.Create(storageLocation)
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}

		defer out.Close()

		// write the content from POST to the file
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		ProcessPDFFile(storageLocation, fileID, w, r)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/get/%s.fact.json", fileID), 303)
	}
}

// Returns a raw text file
func rawTextFileReturn(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	fn := workspaceDir + p[len(p)-1]
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		msg := ResponseStruct{ID: "0", Result: 404, Message: "File not found."}
		b, _ := json.Marshal(msg)

		w.Write(b)
	} else {
		f, _ := os.Open(fn)
		defer f.Close()
		io.Copy(w, f)
	}
}

// Add the default headers so that we calm firefox down
func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}

// Main function, basically it creates an API server that listens on port 3333
func main() {
	_, err := exec.LookPath("java")
	if err != nil {
		log.Fatal("We couldn't find an executable of Java, please make sure you have Java installed. Java can be downloaded at http://java.com/download")
	}

	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

	workspaceDir = UserHomeDir() + "/factsWorkspace/"

	os.MkdirAll(workspaceDir, os.ModeDir|0777)

	log.Printf("Workspace directory is %s", workspaceDir)

	http.HandleFunc("/onlinegrind", addDefaultHeaders(onlineGrind))
	http.Handle("/get/", addDefaultHeaders(rawTextFileReturn))
	http.Handle("/blobsub", addDefaultHeaders(processPDFBlob))

	// If we really are concerned about security, make this 127.0.0.1:3333
	// TODO: HUAN I'm not sure we should do the 127.0.0.1 or not.
	err = http.ListenAndServe("127.0.0.1:3333", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
