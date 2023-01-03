package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Config is the configuration for the font downloader tool.
type Config struct {
	Token string
	Font  string
}

func loadConfig() (Config, error) {
	if len(os.Args[0]) < 2 {
		return Config{}, errors.New("path to key file must be given as an argument")
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// GetAllWithSubstr will download all fonts that contain substr.
// All fonts are downloaded to a folder with name "substr."
func GetAllWithSubstr(token string, substr string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	saveDir := filepath.Join(homeDir, ".fonts", "opentype", substr)
	_, err = os.Stat(saveDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(saveDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	url := "https://www.googleapis.com/webfonts/v1/webfonts?key=" + string(token)
	content, err := urlToString(url)
	if err != nil {
		return err
	}

	var dict map[string]interface{}
	if err := json.Unmarshal([]byte(content), &dict); err != nil {
		return err
	}

	entries := dict["items"].([]interface{})
	for _, entry := range entries {
		font := entry.(map[string]interface{})
		family := font["family"].(string)
		if !strings.Contains(strings.ToLower(family), strings.ToLower(substr+" ")) {
			continue
		}
		downloads := font["files"].(map[string]interface{})
		downloadURL, contains := downloads["regular"]
		if !contains {
			return errors.New("regular version of font not found")
		}

		filepath := filepath.Join(saveDir, family+".ttf")
		err := downloadFile(filepath, downloadURL.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func urlToString(url string) (string, error) {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}

	return "", errors.New("error")
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = GetAllWithSubstr(config.Token, config.Font)
	if err != nil {
		log.Fatal(err)
	}
}
