package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

// PackVersion Details of a pack version
type PackVersion struct {
	Identifier  string // Pack identifier is unique across all language packs. Example: ml-basic-1
	Version     int
	Description string
	Size        int
}

// Pack Details of a pack
type Pack struct {
	Identifier  string
	Name        string
	Description string
	LangCode    string
	Versions    []PackVersion
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Download pack from upstream
func downloadPackFile(langCode string, packVersionIdentifier string) error {
	var (
		fileURL  = fmt.Sprintf("%s/packs/%s/%s", varnamdConfig.upstream, langCode, packVersionIdentifier)
		filePath = path.Join(getPacksDir(), langCode, "a"+packVersionIdentifier)
	)

	resp, err := http.Get(fileURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 response (%s)", resp.Status)
	}

	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func getPackFilePath(langCode string, packVersionIdentifier string) (string, error) {
	pack, err := getPackInfo(langCode)
	if err != nil {
		return "", err
	}

	var packVersion *PackVersion

	for _, version := range pack.Versions {
		if version.Identifier == packVersionIdentifier {
			packVersion = &version
			break
		}
	}

	if packVersion == nil {
		return "", fmt.Errorf("pack version not found")
	}

	// Example: .varnamd/ml/ml-basic-1
	packFilePath := path.Join(getPacksDir(), langCode, packVersionIdentifier)

	if !fileExists(packFilePath) {
		return "", fmt.Errorf("pack file not found")
	}

	return packFilePath, nil
}

func getPackInfo(langCode string) (*Pack, error) {
	packs, err := getPacksInfo()
	if err != nil {
		return nil, err
	}

	for _, pack := range packs {
		if pack.LangCode == langCode {
			return &pack, nil
		}
	}

	return nil, fmt.Errorf("pack not found")
}

func getPacksInfo() ([]Pack, error) {
	if err := createPacksDir(); err != nil {
		return nil, fmt.Errorf("failed to create packs directory, err: %s", err.Error())
	}

	packsFilePath := getPacksDir() + "/packs.json"

	if !fileExists(packsFilePath) {
		err := fmt.Errorf("packs file doesn't exist")
		return nil, err
	}

	packsFile, _ := ioutil.ReadFile(packsFilePath)

	var packsInfo []Pack

	if err := json.Unmarshal(packsFile, &packsInfo); err != nil {
		return nil, fmt.Errorf("parsing packs JSON failed, err: %s", err.Error())
	}

	return packsInfo, nil
}

func createPacksDir() error {
	packsDir := getPacksDir()
	return os.MkdirAll(packsDir, 0750)
}

func getPacksDir() string {
	configDir := getConfigDir()
	return path.Join(configDir, "packs")
}
