/*
Copyright © 2023 ARNAV <r.arnav@icloud.com>
*/

package create

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"pyreqs/cmd/utils"
	"regexp"
	"strings"
	"time"
)

var (
	savePath string
)

var importRegEx = `^import\s+([^\s]+)(\s+as\s+([^\s]+))?$`
var fromImportRegEx = `^from\s+(\w+\.)*(\w+)\s+import\s+.+$`
var httpClient = &http.Client{Timeout: 1000 * time.Second}

// bundling files

//go:embed stdlib
var stdlib embed.FS

//go:embed mappings.txt
var mappings embed.FS

// MAIN

// GetPaths returns the path of python files in the dir
func GetPaths(dirs string, ignoreDirs string) ([]string, []string) {
	var pyFiles []string
	var fileList []string
	var dirList []string
	var ignoreDirList []os.DirEntry

	// directory ignores
	if ignoreDirs != " " {
		ignoreDirsArray := strings.Split(ignoreDirs, ",")
		for _, j := range ignoreDirsArray {
			walk := filepath.WalkDir(j, func(path string, f os.DirEntry, err error) error {
				if f.IsDir() {
					ignoreDirList = append(ignoreDirList, f)
				}
				return nil
			})
			utils.Check(walk)
		}
	}

	walk := filepath.WalkDir(dirs, func(path string, f os.DirEntry, err error) error {
		if f.IsDir() {
			if f.Name() == "venv" || f.Name() == "env" ||
				f.Name() == "__pycache__" || f.Name() == ".git" ||
				f.Name() == ".tox" || f.Name() == ".svn" || f.Name() == ".hg" {
				return filepath.SkipDir
			}
			if ignoreDirList != nil {
				for _, l := range ignoreDirList {
					if f.Name() == l.Name() {
						// TODO: better debug info
						fmt.Println("Skipped: ", l)
						return filepath.SkipDir
					}
				}
			}
		}
		fileList = append(fileList, path)
		return nil
	})
	utils.Check(walk)

	for _, file := range fileList {
		if strings.Contains(file, ".py") {
			pyFiles = append(pyFiles, file)
		} else {
			dirList = append(dirList, file)
		}
	}
	return pyFiles, dirList
}

// ReadImports returns all the imported libraries
func ReadImports(pyFiles []string, dirList []string) map[string]struct{} {

	importLinesSet := map[string]struct{}{}     // for "import" case
	fromImportLinesSet := map[string]struct{}{} // for "from import" case
	imports := map[string]struct{}{}            // final imports

	for _, v := range pyFiles {
		data, err := os.ReadFile(v)
		utils.Check(err)
		codesLines := strings.Split(string(data), "\n")

		importReg, _ := regexp.Compile(importRegEx)
		fromImportReg, _ := regexp.Compile(fromImportRegEx)

		// CASE: import <name> (as <name>); and
		// CASE: from <name> import <name> (as <name>)
		for _, value := range codesLines {
			importLinesSet[strings.TrimSpace(importReg.FindString(strings.TrimSpace(value)))] = struct{}{}
			fromImportLinesSet[strings.TrimSpace(fromImportReg.FindString(strings.TrimSpace(value)))] = struct{}{}
		}
	}

	// SEARCH CASE: import <name> (as <name>)
	for i := range importLinesSet {
		if i != "" {
			importLinesArray := strings.Split(i, " ")
			if strings.Contains(i, " as ") {
				for idx := range importLinesArray {
					if importLinesArray[idx] == "as" { // sep for "as" because it is being included if not done a sep condition
						if idx-1 > 0 && strings.Contains(importLinesArray[idx-1], ".") {
							imports[strings.Split(importLinesArray[idx-1], ".")[0]] = struct{}{}
						} else if idx-1 > 0 {
							imports[importLinesArray[idx-1]] = struct{}{}
						}
					}
				}
			} else {
				for idx := range importLinesArray {
					if importLinesArray[idx] == "import" {
						if idx+1 < len(importLinesArray) && strings.Contains(importLinesArray[idx+1], ".") {
							imports[strings.Split(importLinesArray[idx+1], ".")[0]] = struct{}{}
						} else if idx+1 < len(importLinesArray) {
							imports[importLinesArray[idx+1]] = struct{}{}
						}
					}
				}
			}
		}
	}

	// SEARCH CASE: from <name> import <name> (as <name>)
	for i := range fromImportLinesSet {
		if i != "" {
			fromImportLinesArray := strings.Split(i, " ")
			for idx := range fromImportLinesArray {
				if fromImportLinesArray[idx] == "from" {
					if idx+1 < len(fromImportLinesArray) && strings.Contains(fromImportLinesArray[idx+1], ".") {
						imports[strings.Split(fromImportLinesArray[idx+1], ".")[0]] = struct{}{}
					} else if idx+1 < len(fromImportLinesArray) {
						imports[fromImportLinesArray[idx+1]] = struct{}{}
					}
				}
			}
		}
	}

	// ignore all the directory imports
	for k := range imports {
		for _, l := range dirList {
			if strings.Contains(l, k) {
				delete(imports, k)
			}
		}
	}

	// python inbuilt imports
	predefinedLib, err := stdlib.Open("stdlib")
	utils.Check(err)
	scanner := bufio.NewScanner(predefinedLib)
	inbuiltImportsSet := make(map[string]bool)

	for scanner.Scan() {
		inbuiltImportsSet[scanner.Text()] = true
	}

	for j := range imports {
		if inbuiltImportsSet[j] {
			delete(imports, j)
		}
	}
	return imports
}

// GetLocalPackages returns local packages present in .venv dir
func GetLocalPackages(venvDir string) []string {
	var distPackages []string

	if venvDir == " " { // if not specified, then search for venv in current dir
		cwdDir, err := os.Getwd() // will be the current dir of the build file
		utils.Check(err)
		venvDir = cwdDir
	}
	// TODO: try a search for lib64 dir also
	pattern := filepath.Join(venvDir, "venv/lib/*/*/*.dist-info")
	matches, err := filepath.Glob(pattern)
	utils.Check(err)

	for _, v := range matches {
		list := strings.Split(v, "/")
		lastEle := list[len(list)-1]
		distPackages = append(distPackages, strings.Split(lastEle, ".dist-info")[0])
	}
	return distPackages
}

// FetchPyPIServer fetches and returns info from pypi package server
func FetchPyPIServer(imp []string) map[string]string {
	type Content struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	type Infos struct {
		Info Content `json:"info"` // dont do this: []Content; response isn't a list
	}

	var info Infos
	pypiSet := make(map[string]string)

	predefinedMappings, err := mappings.ReadFile("mappings.txt")
	utils.Check(err)
	mappingsImports := strings.Split(string(predefinedMappings), "\n")
	mappingImportsMap := make(map[string]bool)

	for _, value := range mappingsImports {
		mappingImportsMap[strings.Split(value, ":")[0]] = true
	}

	for _, j := range imp {
		// TODO: do http error handling
		var name string
		// utils.Check for python libraries mapping
		if mappingImportsMap[j] {
			for _, key := range mappingsImports {
				if strings.Split(key, ":")[0] == j {
					name = strings.Split(key, ":")[1]
				}
			}
		} else {
			name = j
		}
		//fmt.Println(name)
		resp, err := httpClient.Get("https://pypi.org/pypi/" + name + "/json")
		utils.Check(err)

		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		// reading response struct using ioutil
		body, err := io.ReadAll(resp.Body)
		utils.Check(err)

		err = json.Unmarshal(body, &info)
		utils.Check(err)
		pypiSet[info.Info.Name] = info.Info.Version

	}
	return pypiSet
}

// write to requirements.txt
func writeRequirements(venvDir string, codesDir string, savePath string, print bool, ignoreDirs string) {
	// type store struct {
	// 	name string
	// 	ver string
	// }

	_, err := os.Stat(filepath.Join(savePath, "requirements.txt"))
	if !os.IsNotExist(err) {
		fmt.Printf("requirements.txt already exists. It will be overwritten.\n")
	}
	file, err := os.Create(filepath.Join(savePath, "requirements.txt"))
	utils.Check(err)

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	pyFiles, dirList := GetPaths(codesDir, ignoreDirs)
	imports := ReadImports(pyFiles, dirList)
	distPackages := GetLocalPackages(venvDir)

	localSet := make(map[string]string)

	for m := range imports {
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
	for i := range imports {
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
	pypiSet := FetchPyPIServer(pypiStore)

	importsInfo := make(map[string]string)
	maps.Copy(importsInfo, localSet)
	maps.Copy(importsInfo, pypiSet)

	for i, j := range importsInfo {
		if i != "" || j != "" {
			fullImport := i + "==" + j + "\n"
			if _, err := file.Write([]byte(fullImport)); err != nil {
				panic(err)
			}
			if print {
				fmt.Println(strings.TrimSuffix(fullImport, "\n"))
			}
		}
	}

	fmt.Println("Created successfully!")

}

// Cmd represents the create command
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Generates a requirements.txt file",
	Long:  `Generates a requirements.txt file.`,
	Run: func(cmd *cobra.Command, args []string) {

		dirPath, err1 := cmd.Flags().GetString("dirPath")
		utils.Check(err1)
		venvPath, err2 := cmd.Flags().GetString("venvPath")
		utils.Check(err2)
		ignoreDirs, err3 := cmd.Flags().GetString("ignore")
		utils.Check(err3)
		printReq, err4 := cmd.Flags().GetBool("print")
		utils.Check(err4)

		writeRequirements(venvPath, dirPath, savePath, printReq, ignoreDirs)
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
	Cmd.Flags().StringVarP(&savePath, "savePath", "s", "./", "save path for requirements.txt")
}
