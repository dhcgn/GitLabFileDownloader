package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

var (
	version = "undef"
)

type GitLapFile struct {
	FileName string `json:"file_name"`
	ContentSha256 string `json:"content_sha256"`
	Content string
}

func main() {
	fmt.Println("Version:", version)

	tokenPtr := flag.String("token", `xxxxxxxxxxxxxxxxxxxx`, "Private-Token")
	fielPathPtr := flag.String("outPath", `my_file.json`, "Path to write the file")
	branchPtr := flag.String("branch", `master`, "Branch")

	urlPtr := flag.String("url", `https://my-git-lap-server.local/api/v4/`, "Url to Api v4")
	projectNumberPtr := flag.Int("projectNumber", 34, "Url to Api v4")
	gitLabFilePathInReproPtr  := flag.String("reproFilePath", `my_config.json`, "gitLabFilePathInReproPtr")

	flag.Parse()

	dir := filepath.Dir(*fielPathPtr)
	if _, err := os.Stat(dir); err == nil {
		fmt.Println("File will be copied to folder", dir, "with file name", filepath.Base(*fielPathPtr))

	} else if os.IsNotExist(err) {
		fmt.Println("Folder",dir,"doesn't exists. Program will exit." )
		os.Exit(-1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	url:=*urlPtr+`projects/`+strconv.Itoa(*projectNumberPtr)+`/repository/files/`+url.PathEscape(*gitLabFilePathInReproPtr)+`?ref=`+*branchPtr
	fmt.Println(url)
	req, err := http.NewRequest("GET", url , nil)
	req.Header.Add("Private-Token", *tokenPtr)

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	robots, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Error from API:",resp.StatusCode, resp.Status)
		os.Exit(-1)
	}

	//fmt.Printf("%s", robots)

	var gitLapFile GitLapFile
	err = json.Unmarshal([]byte(robots), &gitLapFile)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Println( gitLapFile)

	data, err := base64.StdEncoding.DecodeString(gitLapFile.Content)
	if err != nil {
		fmt.Println("error:", err)
	}

	if _, err := os.Stat(*fielPathPtr); err == nil {
		fmt.Println("File exists at",*fielPathPtr)
		file, err := os.Open(*fielPathPtr)
		if err != nil {
			fmt.Println("error:", err)
		}
		defer file.Close()
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			fmt.Println("error:", err)
		}
		hashhex := hex.EncodeToString(hash.Sum(nil))
		fmt.Println("OnDisk", hashhex)
		fmt.Println("GitLab", gitLapFile.ContentSha256)

		if hashhex == gitLapFile.ContentSha256 {
			fmt.Println("No diff, program will exit")
			os.Exit(1)
		}else{
			fmt.Println("File from disk differs to GitLab.")
		}
	}

	err = ioutil.WriteFile(*fielPathPtr, data, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("New file is copied")
	os.Exit(0)
}
