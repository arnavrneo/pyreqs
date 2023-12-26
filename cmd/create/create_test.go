package create

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func Test_fetchPyPIServer(t *testing.T) {
	testImports := []string{"pandas", "numpy", "notapakage"}
	fetchedImports := FetchPyPIServer(testImports)

	for i := range fetchedImports {
		if strings.Contains(i, "notapakage") {
			t.Error("Unknown package crept in the imports file.")
		}
	}
}

func Test_writeRequirements(t *testing.T) {
	writeRequirements("", "testdata/", "./", false, " ")

	reqRead, err := os.Open("requirements.txt")
	if err != nil {
		panic(err)
	}
	defer reqRead.Close()

	testImports := map[string]bool{
		"torch":        true,
		"ipython":      true,
		"tensorflow":   true,
		"pandas":       true,
		"PyYAML":       true,
		"transformers": true,
		"numpy":        true}

	var imports []string
	scanner := bufio.NewScanner(reqRead)
	for scanner.Scan() {
		imports = append(imports, scanner.Text())
	}

	for _, j := range imports {
		if !testImports[strings.Split(j, "==")[0]] {
			t.Error("Missing or extra imports found.")
		}
	}
}
