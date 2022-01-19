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
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/dhcgn/GitLabFileDownloader/internal"
	"github.com/dhcgn/GitLabFileDownloader/internal/api"
)

const (
	AppName = "GitLab File Downloader"
)

var (
	version  = "undef"
	commitID = "undef"

	flagTokenPtr = flag.String(internal.FlagNameToken, ``, `Private-Token with access right for "api" and "read_repository"`)

	flagOutPathPtr      = flag.String(internal.FlagNameOutPath, ``, "Path to write file to disk")
	flagRepoFilePathPar = flag.String(internal.FlagNameRepoFilePath, ``, "File path in repo, like src/main.go")

	flagBranchPtr = flag.String(internal.FlagNameBranch, `master`, "Branch")

	flagOutFolderPtr      = flag.String(internal.FlagNameOutFolder, ``, "Folder to write file to disk")
	flagRepoFolderPathPtr = flag.String(internal.FlagNameRepoFolderPathEscaped, ``, "Folder to write file to disk")

	flagUrlPtr           = flag.String(internal.FlagNameUrl, ``, "Url to Api v4, like https://my-git-lab-server.local/api/v4/")
	flagProjectNumberPtr = flag.Int(internal.FlagNameProjectNumber, 0, "The Project ID from your project")

	exitCode int
)

type GitLapFile struct {
	FileName      string `json:"file_name"`
	ContentSha256 string `json:"content_sha256"`
	Content       string
}

func main() {
	mainSub(os.Args)
}

func mainSub(args []string) {
	log.Println(AppName, "Version:", version, "Commit:", commitID)
	log.Println(`Project: https://github.com/dhcgn/GitLabFileDownloader/`)

	flag.Parse()

	settings := getSettingsFromFlags()
	isValid, args, msgs := settings.IsValid()
	if !isValid {
		log.Println("Arguments are missing:", args)
		log.Println("Messages:", msgs)
		flag.PrintDefaults()
		Exit(-1)
		return
	}

	switch settings.Mode() {
	case internal.ModeFile:
		log.Println("Mode: File")
		fileModeHandling(settings)
	case internal.ModeFolder:
		log.Println("Mode: Folder")
		folderModeHandling(settings)
	}

}

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func folderModeHandling(settings internal.Settings) {
	if !exists(settings.OutFolder) {
		err := os.Mkdir(settings.OutFolder, 0755)
		if err != nil {
			log.Println("Error:", err)
			Exit(-1)
			return
		}
	}

	files, err := api.GetFilesFromFolder(settings)
	if err != nil {
		log.Println("Error:", err)
		Exit(-1)
		return
	}

	log.Println("Sync", len(files), "files")

	for _, file := range files {
		fileSettings := settings

		outFile := path.Join(fileSettings.OutFolder, path.Base(file))
		fileSettings.OutFile = outFile
		fileSettings.RepoFilePath = file

		fileModeHandling(fileSettings)
	}
}

func fileModeHandling(settings internal.Settings) {
	new, err := fileModeHandlingInternal(settings)

	if err != nil {
		log.Println("Error at", settings.RepoFilePath, ":", err)
	}
	if new {
		log.Println("Wrote file:", settings.RepoFilePath, ", because is new or changed")
	} else {
		log.Println("Skip:", settings.RepoFilePath, ", because content is equal")
	}
}

func fileModeHandlingInternal(settings internal.Settings) (bool, error) {
	exists, dir := testTargetFolder(settings.OutFile)
	if !exists {
		return false, fmt.Errorf("Target folder %v doesn't exists", dir)
	}

	err, statusCode, status, bodyData := callApi(settings)
	if err != nil {
		return false, fmt.Errorf("API Call error: %v", err)
	}

	if statusCode != 200 {
		return false, fmt.Errorf("API Call status code: %v %v", statusCode, status)
	}

	gitLapFile, err := createGitLapFile(err, bodyData)
	if err != nil {
		return false, fmt.Errorf("createGitLapFile: %v", err)
	}

	fileData, err := base64.StdEncoding.DecodeString(gitLapFile.Content)
	if err != nil {
		return false, fmt.Errorf("DecodeString: %v", err)
	}

	isEqual, err := isOldFileEqual(gitLapFile, settings)
	if err != nil {
		return false, fmt.Errorf("isOldFileEqual: %v", err)
	}

	if isEqual {
		return false, nil
	}

	err = ioutil.WriteFile(settings.OutFile, fileData, 0644)
	if err != nil {
		return false, fmt.Errorf("WriteFile: %v", err)
	}
	return true, nil
}

func Exit(code int) {
	log.Println("Program will exit!")
	exitCode = code
	// only deployed version should terminate here
	if version != "undef" {
		os.Exit(code)
	}
}

func isOldFileEqual(gitLapFile GitLapFile, settings internal.Settings) (bool, error) {
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

func callApi(settings internal.Settings) (error, int, string, []byte) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	apiUrl := fmt.Sprintf("%vprojects/%v/repository/files/%v?ref=%v", *flagUrlPtr, settings.ProjectNumber, url.QueryEscape(settings.RepoFilePath), settings.Branch)

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

func getSettingsFromFlags() internal.Settings {
	return internal.Settings{
		PrivateToken:   *flagTokenPtr,
		OutFile:        *flagOutPathPtr,
		OutFolder:      *flagOutFolderPtr,
		Branch:         *flagBranchPtr,
		ApiUrl:         *flagUrlPtr,
		ProjectNumber:  strconv.Itoa(*flagProjectNumberPtr),
		RepoFilePath:   *flagRepoFilePathPar,
		RepoFolderPath: *flagRepoFolderPathPtr,
	}
}
