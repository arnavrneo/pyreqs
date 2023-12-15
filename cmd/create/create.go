/*
Copyright © 2023 ARNAV RAINA <r.arnav@icloud.com>
*/
package create

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	dirPath  string
	venvPath string
)

// TODO: improve regex
var importRegEx string = `^(?m)[\s\S]*import\s+(\w+)(\s+as\s+\w+)?\s*$` // m: m is multiline flag
// var fromRegEx string = `^(?m)[\s\S]*from\s+(\w+(\.\w+)*)\s+import\s+(\*\s+as\s+\w+|\w+(\s+as\s+\w+)?)\s*$`
var myClient = &http.Client{Timeout: 1000 * time.Second}

// error check
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// return the path of .py files in the dir
func getPaths(dirs string) ([]string, map[os.DirEntry]struct{}) {
	var pyFiles []string
	fileList := []string{}
	dirList := map[os.DirEntry]struct{}{}

	// ignoreDirs := []string {".hg", ".svn", ".git", ".tox", "__pycache__", "env", "venv"}
	err := filepath.WalkDir(dirs, func(path string, f os.DirEntry, err error) error {
		if f.IsDir() && (f.Name() == "venv" || f.Name() == "env" || f.Name() == "__pycache__") {
			return filepath.SkipDir
		}
		fileList = append(fileList, path)
		return nil
	})
	check(err)

	for _, file := range fileList {
		if strings.Contains(file, ".py") {
			pyFiles = append(pyFiles, file)
		}
	}

	return pyFiles, dirList
}

// return all the imported libraries
func readImports(pyFiles []string, dirList map[os.DirEntry]struct{}) map[string]struct{} {
	importSet := map[string]struct{}{}    // set
	importStruct := map[string]struct{}{} // for implementing a set

	// gives a subset of the whole file
	// code till the line which has the last import
	for _, v := range pyFiles {
		data, err := os.ReadFile(v)
		check(err)
		deleteThis := strings.Split(string(data), "\n") // now done a split at \n followed by regex search
		r, _ := regexp.Compile(importRegEx)
		for _, value := range deleteThis {
			importSet[r.FindString(value)] = struct{}{} // maintained a set for imports
		}
		fmt.Println("==IMPORT==", importSet)
	}

	// search in that subset
	// TODO: a better way of getting imports; no \n or special characters
	for i := range importSet {
		arr := strings.Split(i, " ")
		for i, j := range arr {
			if strings.Contains(j, "import") {
				if i-2 >= 0 && strings.Contains(arr[i-2], "from") {
					if strings.Contains(arr[i-1], ".") {
						// storedImports = append(storedImports, strings.Split(arr[i-1], ".")[0])
						importStruct[strings.TrimSpace(strings.Split(arr[i-1], ".")[0])] = struct{}{}
					} else {
						importStruct[strings.TrimSpace(arr[i-1])] = struct{}{}
					}
				} else {
					if strings.Contains(arr[i+1], ".") {
						importStruct[strings.TrimSpace(strings.Split(arr[i+1], ".")[0])] = struct{}{}
					} else {
						importStruct[strings.TrimSpace(arr[i+1])] = struct{}{}
					}
				}
			}
		}
	}

	// ignore all the directory imports
	for k := range importStruct {
		for l := range dirList {
			if k == l.Name() {
				delete(importStruct, k)
			}
		}
	}

	// python inbuilt imports
	predefinedLib, err := os.ReadFile("cmd/stdlib")
	check(err)
	inbuiltImports := strings.Split(string(predefinedLib), " ")
	for j := range importStruct {
		for _, k := range inbuiltImports {
			if j == k {
				delete(importStruct, k)
			}
		}
	}
	return importStruct
}

// return local packages installed in .venv dir
func getLocalPackages(venvDir string) []string {
	var distPackages []string

	if venvDir == " " { // if not specified, then search for venv in current dir
		dir, err := os.Getwd()
		check(err)
		venvDir = dir
	}
	// TODO: try a search for lib64 dir also
	pattern := filepath.Join(venvDir, "venv/lib/*/*/*.dist-info")
	matches, err := filepath.Glob(pattern)
	check(err)

	for _, v := range matches {
		list := strings.Split(v, "/")
		lastEle := list[len(list)-1]
		distPackages = append(distPackages, strings.Split(lastEle, ".dist-info")[0])
	}
	return distPackages
}

// fetch and return info from pypi package server
func fetchPyPIServer(imp []string) map[string]string {
	type Content struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	type Infos struct {
		Info Content `json:"info"` // dont do this: []Content; response isnt a list
	}

	var info Infos
	pypiSet := make(map[string]string)
	mappings, err := os.ReadFile("cmd/mappings.txt")
	check(err)
	mappingsImports := strings.Split(string(mappings), "\n")
	//mappingSet := map[string]struct{ string }{}
	mappingImportsMap := make(map[string]bool)

	for _, value := range mappingsImports {
		mappingImportsMap[strings.Split(value, ":")[0]] = true
	}

	for _, j := range imp {
		// TODO: try to look for string formatting ways instead of "+"
		// TODO: do http error handling
		var name string
		// check for python libraries mapping
		if mappingImportsMap[j] {
			for _, key := range mappingsImports {
				if strings.Split(key, ":")[0] == j {
					name = strings.Split(key, ":")[1]
				}
			}
		} else {
			name = j
		}
		resp, err := myClient.Get("https://pypi.org/pypi/" + name + "/json")
		check(err)
		defer resp.Body.Close()

		// reading response struct using ioutil
		body, err := io.ReadAll(resp.Body)
		check(err)

		json.Unmarshal(body, &info)
		pypiSet[info.Info.Name] = info.Info.Version

	}
	return pypiSet
}

// write to requirements.txt
func writeRequirements(venvDir string, codesDir string) {
	// type store struct {
	// 	name string
	// 	ver string
	// }

	file, err := os.Create("requirements.txt")
	check(err)

	defer func() { // delays the execution until the surrounding function completes
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	pyFiles, dirList := getPaths(codesDir)
	importStruct := readImports(pyFiles, dirList)
	distPackages := getLocalPackages(venvDir)

	localSet := make(map[string]string)

	for m := range importStruct {
		for _, n := range distPackages {
			name := strings.Split(n, "-")[0]
			ver := strings.Split(n, "-")[1]
			if m == name {
				localSet[name] = ver
			}
		}
	}

	// imports from pypi server
	var pypiStore []string
	for i := range importStruct {
		cntr := 0
		for j := range localSet {
			if i == j {
				cntr += 1
			}
		}
		if cntr == 0 {
			// name, ver := fetchPyPIServer(i)
			pypiStore = append(pypiStore, i)
		}
	}
	pypiSet := fetchPyPIServer(pypiStore)

	importsInfo := make(map[string]string)
	maps.Copy(importsInfo, localSet)
	maps.Copy(importsInfo, pypiSet)

	for i, j := range importsInfo {
		if i != "" && j != "" { // TODO: check why should this be done; find workaround
			fullImport := i + "==" + j + "\n"
			if _, err := file.Write([]byte(fullImport)); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Created successfully!")
}

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a requirements.txt file",
	Long:  `LONG DESC`,
	Run: func(cmd *cobra.Command, args []string) {
		writeRequirements(venvPath, dirPath)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	CreateCmd.Flags().StringVarP(&dirPath, "dirPath", "d", "./", "directory to .py files")
	CreateCmd.Flags().StringVarP(&venvPath, "venvPath", "v", " ", "directory to venv")

}
