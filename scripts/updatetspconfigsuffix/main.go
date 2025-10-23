package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
			firstLine := getGoModFirstLine(path)
			arr := strings.Split(firstLine, "/")
			suffix := arr[len(arr)-1]
			lastChar := suffix[len(suffix)-1]
			if !isNumber(lastChar) {
				return nil
			}
			configFile := getConfigFilePath(path)
			fmt.Println("\"" + configFile + "\",")
			configFiledep := "D:/GoProject/azure-rest-api-specs" + "/" + configFile + "/tspconfig.yaml"
			module, _ := getPackageModule(configFiledep)
			arr2 := strings.Split(module, "/")
			if isNumber(arr2[len(arr2)-1][len(arr2[len(arr2)-1])-1]) {
				arr2[len(arr2)-1] = suffix
			} else {
				arr2 = append(arr2, suffix)
			}
			newModule := strings.Join(arr2, "/")
			updateModule(configFiledep, newModule)
			fmt.Println("Updated module to:", newModule)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
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

func getPackageModule(path string) (string, error) {
	data, err := ioutil.ReadFile(path) // Replace with your YAML file path
	if err != nil {
		return "", errors.New("failed to read file")
	}

	// Unmarshal into a generic map
	var root map[string]interface{}
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return "", errors.New("Failed to parse YAML")
	}

	// Navigate to options → @azure-tools/typespec-go
	options, ok := root["options"].(map[string]interface{})
	if !ok {
		return "", errors.New("Missing or invalid 'options' section")
	}

	goOptions, ok := options["@azure-tools/typespec-go"].(map[string]interface{})
	if !ok {
		return "", errors.New("Missing or invalid '@azure-tools/typespec-go' section")
	}
	// Extract specific fields
	return goOptions["module"].(string), nil
}

// Define minimal structure matching tspconfig.yaml
type TspConfig struct {
	Options map[string]map[string]interface{} `yaml:"options"`
}

func updateModule(path, newModule string) {

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	// Parse YAML
	var config TspConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(fmt.Errorf("failed to parse yaml: %w", err))
	}

	// Navigate to @azure-tools/typespec-go
	goOptions, ok := config.Options["@azure-tools/typespec-go"]
	if !ok {
		panic("no '@azure-tools/typespec-go' section found")
	}

	// Update module
	goOptions["module"] = newModule

	// Write back
	out, err := yaml.Marshal(&config)
	if err != nil {
		panic(fmt.Errorf("failed to marshal yaml: %w", err))
	}
	if err := os.WriteFile(path, out, 0644); err != nil {
		panic(fmt.Errorf("failed to write file: %w", err))
	}

	fmt.Println("Updated @azure-tools/typespec-go.module to:", newModule)
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

func getGoModFirstLine(path string) string {
	file, err := os.Open(path + "/go.mod")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		firstLine := scanner.Text()
		fmt.Println("First line:", firstLine)
		return firstLine
	} else if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	} else {
		fmt.Println("File is empty.")
	}
	return ""
}
