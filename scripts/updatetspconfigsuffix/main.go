package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"updatetspconfigsuffix/model"

	"gopkg.in/yaml.v2"
)

func main() {
	rootDir := "D:/GoProject2/azure-sdk-for-go/sdk/resourcemanager"
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

		// Check for spec match
		/*
		if d.IsDir() && fileExists(path+"\\"+"tsp-location.yaml") && strings.Contains(path, "\\arm") && !strings.Contains(path, "fake") && !strings.Contains(path, "node_modules") {
			if strings.Contains(content, path) {
				// fmt.Println(path, " has been generated")
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
			configFiledep := "D:/GoProject2/azure-rest-api-specs" + "/" + configFile + "/tspconfig.yaml"
			if !fileExists(configFiledep) {
				return nil
			}
			module, err := getPackageModule(configFiledep)
			if err != nil {
				fmt.Println("Error getting package module:", err, "configfile path :", configFiledep)
				return nil
			} 
			arr2 := strings.Split(module, "/")
			if isNumber(arr2[len(arr2)-1][len(arr2[len(arr2)-1])-1]) {
				arr2[len(arr2)-1] = suffix
			} else {
				arr2 = append(arr2, suffix)
			}
			newModule := strings.Join(arr2, "/")
			fmt.Println("    Updated module to:", newModule)
			return nil
		}*/

		// Check for autorest match
		if d.IsDir() && fileExists(path+"\\"+"autorest.md") && strings.Contains(path, "\\arm") && !strings.Contains(path, "fake") && !strings.Contains(path, "node_modules") {
			if strings.Contains(content, path) {
				// fmt.Println(path, " has been generated")
				return nil
			}
			firstLine := getGoModFirstLine(path)
			arr := strings.Split(firstLine, "/")
			suffix := arr[len(arr)-1]
			lastChar := suffix[len(suffix)-1]
			if !isNumber(lastChar) {
				return nil
			}
			configFile := getAutorestConfigFilePath(path)
			if configFile == "" {
				fmt.Println("firstLine:", firstLine)
				return nil
			}
			configFiledep := "D:/GoProject2/azure-rest-api-specs" + "/" + configFile + "readme.go.md"
			// fmt.Println(configFiledep, "===================", configFile)
			if !fileExists(configFiledep) {
				fmt.Println("Folder changed===================", configFiledep, "suffix:", suffix)
				return nil
			}
			module, err := getReadmeGoPackageModule(configFiledep)
			if err != nil {
				fmt.Println("Error getting package module:", err, "configfile path :", configFiledep)
				return nil
			} 

			arr2 := strings.Split(module, "/")
			if isNumber(arr2[len(arr2)-1][len(arr2[len(arr2)-1])-1]) {
				fmt.Println("old module is-----:", strings.Join(arr2, "/"))
				arr2[len(arr2)-1] = suffix
			} else {
				arr2 = append(arr2, suffix)
			}
			newModule := strings.Join(arr2, "/")
			fmt.Println("    Updated module to-----:", newModule)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
	fmt.Println("scan successfully")
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func getConfigFilePath(path string) string {
	filePath := path + "/tsp-location.yaml" // replace with your file path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return ""
	}

	// Unmarshal into a map
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	return result["directory"].(string) // return the value for the specified key

}

func getAutorestConfigFilePath(path string) string {
	filePath := path + "/autorest.md" // replace with your file path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return ""
	}
	// Regex to extract path after 'blob/<commit>/'
	// re := regexp.MustCompile(`github\.com/.+?/blob/[^/]+/(specification/[^ ]+)`)
	re := regexp.MustCompile(`specification/[^/]+/resource-manager/readme\.go\.md`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	if len(matches) == 0 {
		// fmt.Println("No specification paths found")
		return ""
	}
	s := strings.Replace(matches[0][0], "readme.go.md", "", 1)
	return s
}

func getReadmeGoPackageModule(path string) (string, error) {
	// Read the markdown file
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Regex to capture module-name value
	re := regexp.MustCompile(`(?m)^\s*module-name:\s*([^\s]+)\s*$`)
	matches := re.FindSubmatch(data)

	if len(matches) < 2 {
		// fmt.Println("module-name not found")
		return "", fmt.Errorf("module-name not found")
	}

	moduleName := string(matches[1])
	// fmt.Println("module-name:", moduleName)
	return moduleName, nil
}


func getPackageModule(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
		return "", err
	}

	var cfg model.YamlConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
		return "", err
	}
	fmt.Println("Go SDK module path: →", cfg.Options.Go.Module)
	return cfg.Options.Go.Module, nil
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
		// fmt.Println("First line:", firstLine)
		return firstLine
	} else if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	} else {
		fmt.Println("File is empty.")
	}
	return ""
}
