package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	rootDir := "D:/GoProject/azure-sdk-for-go/sdk/resourcemanager"
	filePath := "example.txt" // Replace with your file path
	content := ""
	data, err := os.ReadFile(filePath)
	if err == nil {
		content = string(data)
	}

	err = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // stop walking if error
		}

		// Check for exact match
		if d.IsDir() && fileExists(path+"\\"+"tsp-location.yaml") && strings.Contains(path, "\\arm") && !strings.Contains(path, "fake") && !strings.Contains(path, "node_modules") {
			if strings.Contains(content, path) {
				fmt.Println(path, " has been generated")
				return nil
			}
			configFile := getConfigFilePath(path)
			fmt.Println("\"" + configFile + "\",")
			return nil // continue walking even if there's an error

			configFiledep := "D:/GoProject/azure-rest-api-specs" + "/" + configFile + "/tspconfig.yaml"
			sp1, sp2, err := getPackageDirectory(configFiledep)
			if err != nil {
				// fmt.Println("Error getting package directory:", err, "configfile path :", configFiledep)
				return nil // continue walking even if there's an error
			}
			cmd1 := " release-v2  D:\\GoProject\\azure-sdk-for-go  D:\\GoProject\\azure-rest-api-specs " + sp1 + " " + sp2 + " --tsp-config " + "\"" + configFiledep + "\""
			fmt.Println("running ", cmd1)
			// return nil
			cmd := exec.Command("generator", cmd1)
			err = cmd.Run()
			if err != nil {
				fmt.Println("run fail", err.Error())
			} else {
				appendLine(path)
				fmt.Println("Line appended successfully to example.txt")
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func appendLine(line string) {
	filePath := "example.txt" // Replace with your target file
	// Open the file in append mode
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()
	// Append the line
	if _, err := f.WriteString(line); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	log.Println("Line appended successfully.")
}

func getConfigFilePath(path string) string {
	filePath := path + "/tsp-location.yaml" // replace with your file path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal into a map
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	return result["directory"].(string) // return the value for the specified key

}

func getPackageDirectory(path string) (string, string, error) {
	data, err := ioutil.ReadFile(path) // Replace with your YAML file path
	if err != nil {
		return "", "", errors.New("failed to read file")
	}

	// Unmarshal into a generic map
	var root map[string]interface{}
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return "", "", errors.New("Failed to parse YAML")
	}

	// Navigate to options → @azure-tools/typespec-go
	options, ok := root["options"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("Missing or invalid 'options' section")
	}

	goOptions, ok := options["@azure-tools/typespec-go"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("Missing or invalid '@azure-tools/typespec-go' section")
	}

	// Extract specific fields
	serviceDir := goOptions["service-dir"]
	packageDir := goOptions["package-dir"]
	arr1 := strings.Split(serviceDir.(string), "/")
	arr2 := strings.Split(packageDir.(string), "/")
	return arr1[len(arr1)-1], arr2[len(arr2)-1], nil // Return the last part of each path
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return false // Other error, treat as non-existent
}
