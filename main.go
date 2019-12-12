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

	flagToken         = "token"
	flagOutPath       = "outPath"
	flagBranch        = "branch"
	flagUrl           = "url"
	flagProjectNumber = "projectNumber"
	flagReproFilePath = "reproFilePath"
)

var (
	version = "undef"

	flagTokenPtr   = flag.String(flagToken, ``, "Private-Token with access right for api and read_repository")
	flagOutPathPtr = flag.String(flagOutPath, ``, "Path to write file to disk")
	flagBranchPtr  = flag.String(flagBranch, `master`, "Branch")

	flagUrlPtr           = flag.String(flagUrl, ``, "Url to Api v4, like https://my-git-lab-server.local/api/v4/")
	flagProjectNumberPtr = flag.Int(flagProjectNumber, 0, "The Project ID from your project")
	flagReproFilePathPtr = flag.String(flagReproFilePath, ``, "File path in repro, like src/main.go")
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
	MachineId            string
}

func main() {
	fmt.Println(AppName, "Version:", version)
	fmt.Println(`Project: https://github.com/dhcgn/GitLabFileDownloader/`)

	if len(os.Args) == 2 && os.Args[1] == "update" {
		equinoxUpdate()
	}

	flag.Parse()

	settings := getSettings()
	isValid, args := isSettingsValid(settings)
	if !isValid {
		fmt.Println("Arguments are missing:", args)
		fmt.Println("Program will exit!")
		os.Exit(-1)
	}

	exists, dir := testTargetFolder(settings.OutFile)
	if !exists {
		fmt.Println("Folder", dir, "doesn't exists. Program will exit.")
		os.Exit(-1)
	}

	err, statusCode, status, bodyData := callApi(settings)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	if statusCode != 200 {
		fmt.Println("Error from API call:", statusCode, status)
		os.Exit(-1)
	}

	gitLapFile, err := createGitLapFile(err, bodyData)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fileData, err := base64.StdEncoding.DecodeString(gitLapFile.Content)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	isEqual, err := isOldFileEqual(gitLapFile, settings)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	if isEqual {
		fmt.Println("No diff, nothing to do. Program will exit.")
		os.Exit(1)
	} else {
		fmt.Println("File from disk differs.")
	}

	err = ioutil.WriteFile(*flagOutPathPtr, fileData, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("New file was copied.")
	os.Exit(0)
}

func isSettingsValid(settings settings) (bool, []string) {
	var missingArgs []string

	if settings.PrivateToken == "" {
		missingArgs = append(missingArgs, flagToken)
	}
	if settings.OutFile == "" {
		missingArgs = append(missingArgs, flagOutPath)
	}
	if settings.Branch == "" {
		missingArgs = append(missingArgs, flagBranch)
	}
	if settings.ReproFilePathEscaped == "" {
		missingArgs = append(missingArgs, flagReproFilePath)
	}
	if settings.ApiUrl == "" {
		missingArgs = append(missingArgs, flagUrl)
	}

	return len(missingArgs) == 0, missingArgs
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
	apiUrl := *flagUrlPtr + `projects/` + settings.ProjectNumber + `/repository/files/` + settings.ReproFilePathEscaped + `?ref=` + settings.Branch

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

func testTargetFolder(outFile string) (exists bool, dir string) {
	dir = filepath.Dir(outFile)
	if _, err := os.Stat(dir); err == nil {
		return true, dir
	} else if os.IsNotExist(err) {
		return false, dir
	}
	return false, dir
}

func getSettings() settings {
	id, err := machineid.ProtectedID(AppName)
	if err != nil {
		log.Fatal(err)
	}
	return settings{
		PrivateToken:         *flagTokenPtr,
		OutFile:              *flagOutPathPtr,
		Branch:               *flagBranchPtr,
		ApiUrl:               *flagUrlPtr,
		ProjectNumber:        strconv.Itoa(*flagProjectNumberPtr),
		ReproFilePathEscaped: url.PathEscape(*flagReproFilePathPtr),
		MachineId:            id,
	}
}
