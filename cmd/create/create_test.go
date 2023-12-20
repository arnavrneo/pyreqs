package create

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func Test_getPaths(t *testing.T) {
	dirs := "testdata"
	_, dirList := getPaths(dirs, "testdata/pyTestCodes")

	for _, j := range dirList {
		if strings.Contains(j, "venv") || strings.Contains(j, "env") || strings.Contains(j, "__pycache__") || strings.Contains(j, ".git") || strings.Contains(j, ".tox") {
			t.Error("Directory list has ignored directories.")
		}
	}
}

//func Test_readImports(t *testing.T) {
//	dirs := "testdata"
//	pyFiles, dirList := getPaths(dirs)
//
//	imports := readImports(pyFiles, dirList)
//
//	for i := range imports {
//		if strings.ContainsAny(i, "\n | \n\n | \n\n\n") { // check for "" or " " gives error; look into
//			t.Error("The imports have a new-line.")
//		}
//	}
//}

func Test_fetchPyPIServer(t *testing.T) {
	testImports := []string{"pandas", "numpy", "notapakage"}
	fetchedImports := fetchPyPIServer(testImports)

	for i := range fetchedImports {
		if strings.Contains(i, "notapakage") {
			t.Error("Unknown package crept in the imports file.")
		}
	}
}

func Test_writeRequirements(t *testing.T) {
	writeRequirements("", "testdata", "./", false)

	reqRead, err := os.Open("requirements.txt")
	if err != nil {
		panic(err)
	}
	defer reqRead.Close()

	testImports := map[string]bool{"torch": true, "pandas": true, "transformers": true, "tensorflow": true, "PyYAML": true, "ipython": true, "numpy": true}
	var imports []string
	scanner := bufio.NewScanner(reqRead)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		imports = append(imports, scanner.Text())
	}

	for _, j := range imports {
		if !testImports[strings.Split(j, "==")[0]] {
			t.Error("Missing or extra imports found.")
		}
	}
}
