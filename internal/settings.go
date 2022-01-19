package internal

import "fmt"

const (
	FlagNameToken                 = "token"
	FlagNameOutPath               = "outPath"
	FlagNameOutFolder             = "outFolder"
	FlagNameBranch                = "branch"
	FlagNameUrl                   = "url"
	FlagNameProjectNumber         = "projectNumber"
	FlagNameRepoFilePath          = "repoFilePath"
	FlagNameRepoFolderPathEscaped = "repoFolder"
)

type Settings struct {
	PrivateToken   string
	OutFile        string
	OutFolder      string
	Branch         string
	ApiUrl         string
	ProjectNumber  string
	RepoFilePath   string
	RepoFolderPath string
	UserAgent      string
}

type Mode int

const (
	ModeUndef = iota
	ModeFile
	ModeFolder
)

func (s Settings) Mode() Mode {
	if s.OutFile != "" && s.RepoFilePath != "" {
		return ModeFile
	}
	if s.OutFolder != "" && s.RepoFolderPath != "" {
		return ModeFolder
	}
	return ModeUndef
}

func (s Settings) IsValid() (bool, []string, []string) {
	var missingArgs []string
	var errors []string

	if s.PrivateToken == "" {
		missingArgs = append(missingArgs, FlagNameToken)
	}
	// Missmatch between outFile and outFolder is not allowed
	if s.OutFile == "" && s.OutFolder == "" {
		missingArgs = append(missingArgs, FlagNameOutPath)
	}
	if s.OutFolder == "" && s.OutFile == "" {
		missingArgs = append(missingArgs, FlagNameOutFolder)
	}
	if s.OutFolder != "" && s.OutFile != "" {
		errors = append(errors, fmt.Sprint("You can't use both ", FlagNameOutPath, " and ", FlagNameOutFolder))
	}
	if s.OutFolder != "" && s.RepoFolderPath == "" {
		missingArgs = append(missingArgs, FlagNameRepoFolderPathEscaped)
	}
	if s.RepoFolderPath != "" && s.OutFolder == "" {
		missingArgs = append(missingArgs, FlagNameOutFolder)
	}
	if s.Branch == "" {
		missingArgs = append(missingArgs, FlagNameBranch)
	}
	if s.RepoFilePath == "" && s.RepoFolderPath == "" {
		missingArgs = append(missingArgs, FlagNameRepoFilePath)
	}
	if s.ApiUrl == "" {
		missingArgs = append(missingArgs, FlagNameUrl)
	}

	return len(missingArgs) == 0, missingArgs, errors
}
