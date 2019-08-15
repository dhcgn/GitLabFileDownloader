package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

const (
	AppName = "GitLab File Downloader"
)

var (
	version = "undef"

	tokenPtr    = flag.String("token", ``, "Private-Token")
	fielPathPtr = flag.String("outPath", ``, "Path to write the file")
	branchPtr   = flag.String("branch", `master`, "Branch")

	urlPtr                   = flag.String("url", ``, "Url to Api v4, like https://my-git-lap-server.local/api/v4/")
	projectNumberPtr         = flag.Int("projectNumber", 0, "Url to Api v4")
	gitLabFilePathInReproPtr = flag.String("reproFilePath", ``, "File path in repro")
)

type GitLapFile struct {
	FileName      string `json:"file_name"`
	ContentSha256 string `json:"content_sha256"`
	Content       string
}

type settings struct {
	PrivateToken         string
	OutFile              string
	Branch               string
	ApiUrl               string
	ProjectNumber        string
	ReproFilePathEscaped string
	MachineId string
}

func main() {
	fmt.Println(AppName, "Version:", version)
	flag.Parse()

	settings := getSettings()

	checkFileLocation(settings.OutFile)

	err, statusCode, status, bodyData := callApi(settings)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(-1)
	}

	if statusCode != 200 {
		fmt.Println("Error from API:", statusCode, status)
		os.Exit(-1)
	}

	gitLapFile, err := createGitLapFile(err, bodyData)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(-1)
	}

	fileData, err := base64.StdEncoding.DecodeString(gitLapFile.Content)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(-1)
	}

	isEqual, err := isOldFileEqual(gitLapFile, settings)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(-1)
	}

	if isEqual {
		fmt.Println("No diff, program will exit")
		os.Exit(1)
	} else {
		fmt.Println("File from disk differs to GitLab.")
	}

	err = ioutil.WriteFile(*fielPathPtr, fileData, 0644)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(-1)
	}

	fmt.Println("New file is copied")
	os.Exit(0)
}

func isOldFileEqual(gitLapFile GitLapFile, settings settings) (bool, error) {
	if _, err := os.Stat(settings.OutFile); err == nil {
		file, err := os.Open(settings.OutFile)
		if err != nil {
			return false, err
		}
		defer file.Close()
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			return false, err
		}
		sha256Hex := hex.EncodeToString(hash.Sum(nil))

		return sha256Hex == gitLapFile.ContentSha256, nil
	}
	return false, nil
}

func createGitLapFile(err error, data []byte) (GitLapFile, error) {
	var gitLapFile GitLapFile
	err = json.Unmarshal(data, &gitLapFile)
	return gitLapFile, err
}

func callApi(settings settings) (error, int, string, []byte) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	apiUrl := *urlPtr + `projects/` + settings.ProjectNumber + `/repository/files/` + settings.ReproFilePathEscaped + `?ref=` + settings.Branch

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err, 0, "", nil
	}
	req.Header.Add("Private-Token", settings.PrivateToken)
	req.Header.Add("MachineId", settings.MachineId)

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err, 0, "", nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, 0, "", nil
	}
	err = resp.Body.Close()
	if err != nil {
		return err, 0, "", nil
	}
	return err, resp.StatusCode, resp.Status, body
}

func checkFileLocation(outFile string) {
	dir := filepath.Dir(outFile)
	if _, err := os.Stat(dir); err == nil {
		fmt.Println("File will be copied to folder", dir, "with file name", filepath.Base(outFile))

	} else if os.IsNotExist(err) {
		fmt.Println("Folder", dir, "doesn't exists. Program will exit.")
		os.Exit(-1)
	}
}

func getSettings() settings {
	id, err := machineid.ProtectedID(AppName)
	if err != nil {
		log.Fatal(err)
	}
	return settings{
		PrivateToken:         *tokenPtr,
		OutFile:              *fielPathPtr,
		Branch:               *branchPtr,
		ApiUrl:               *urlPtr,
		ProjectNumber:        strconv.Itoa(*projectNumberPtr),
		ReproFilePathEscaped: url.PathEscape(*gitLabFilePathInReproPtr),
		MachineId: id,
	}
}
