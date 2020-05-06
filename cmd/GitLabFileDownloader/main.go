package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dhcgn/GitLabFileDownloader/cmd/GitLabFileDownloader/updater"
)

const (
	AppName = "GitLab File Downloader"

	flagToken         = "token"
	flagOutPath       = "outPath"
	flagBranch        = "branch"
	flagUrl           = "url"
	flagProjectNumber = "projectNumber"
	flagRepoFilePath  = "repoFilePath"
	flagUpdate        = "update"
)

var (
	version = "undef"

	flagTokenPtr   = flag.String(flagToken, ``, `Private-Token with access right for "api" and "read_repository"`)
	flagOutPathPtr = flag.String(flagOutPath, ``, "Path to write file to disk")
	flagBranchPtr  = flag.String(flagBranch, `master`, "Branch")

	flagUrlPtr           = flag.String(flagUrl, ``, "Url to Api v4, like https://my-git-lab-server.local/api/v4/")
	flagProjectNumberPtr = flag.Int(flagProjectNumber, 0, "The Project ID from your project")
	flagRepoFilePathPar  = flag.String(flagRepoFilePath, ``, "File path in repo, like src/main.go")

	flagUpdatePtr = flag.Bool(flagUpdate, false, "Update executable from equinox.io")

	exitCode int
	Args     = os.Args
)

type GitLapFile struct {
	FileName      string `json:"file_name"`
	ContentSha256 string `json:"content_sha256"`
	Content       string
}

type settings struct {
	PrivateToken        string
	OutFile             string
	Branch              string
	ApiUrl              string
	ProjectNumber       string
	RepoFilePathEscaped string
}

func main() {
	log.Println(AppName, "Version:", version)
	log.Println(`Project: https://github.com/dhcgn/GitLabFileDownloader/`)

	if len(Args) == 2 && Args[1] == "update" {
		updater.EquinoxUpdate()
		Exit(2)
		return
	}

	flag.Parse()

	if *flagUpdatePtr == true {
		updater.EquinoxUpdate()
		Exit(2)
		return
	}

	settings := getSettings()
	isValid, args := isSettingsValid(settings)
	if !isValid {
		log.Println("Arguments are missing:", args)
		Exit(-1)
		return
	}

	exists, dir := testTargetFolder(settings.OutFile)
	if !exists {
		log.Println("Folder", dir, "doesn't exists.")
		Exit(-1)
		return
	}

	err, statusCode, status, bodyData := callApi(settings)
	if err != nil {
		log.Println("Error:", err)
		Exit(-1)
		return
	}

	if statusCode != 200 {
		log.Println("Error from API call:", statusCode, status)
		Exit(-1)
		return
	}

	gitLapFile, err := createGitLapFile(err, bodyData)
	if err != nil {
		log.Println("Error:", err)
		Exit(-1)
		return
	}

	fileData, err := base64.StdEncoding.DecodeString(gitLapFile.Content)
	if err != nil {
		log.Println("Error:", err)
		Exit(-1)
		return
	}

	isEqual, err := isOldFileEqual(gitLapFile, settings)
	if err != nil {
		log.Println("Error:", err)
		Exit(-1)
		return
	}

	if isEqual {
		log.Println("No diff, nothing to do.")
		Exit(1)
		return
	} else {
		log.Println("File from disk differs")
	}

	err = ioutil.WriteFile(*flagOutPathPtr, fileData, 0644)
	if err != nil {
		log.Println("Error:", err)
		Exit(-1)
		return
	}

	log.Println("New file was copied")
	Exit(0)
}

func Exit(code int) {
	log.Println("Program will exit!")
	exitCode = code
	// only deployed version should terminate here
	if version != "undef" {
		os.Exit(code)
	}
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
	if settings.RepoFilePathEscaped == "" {
		missingArgs = append(missingArgs, flagRepoFilePath)
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
	apiUrl := *flagUrlPtr + `projects/` + settings.ProjectNumber + `/repository/files/` + settings.RepoFilePathEscaped + `?ref=` + settings.Branch

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err, 0, "", nil
	}
	req.Header.Add("Private-Token", settings.PrivateToken)
	req.Header.Add("User-Agent", AppName+" "+version)

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
	return settings{
		PrivateToken:        *flagTokenPtr,
		OutFile:             *flagOutPathPtr,
		Branch:              *flagBranchPtr,
		ApiUrl:              *flagUrlPtr,
		ProjectNumber:       strconv.Itoa(*flagProjectNumberPtr),
		RepoFilePathEscaped: url.PathEscape(*flagRepoFilePathPar),
	}
}
