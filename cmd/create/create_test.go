package create

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func Test_getPaths(t *testing.T) {
	dirs := "ultralytics/"
	_, dirList := getPaths(dirs, "")

	for _, j := range dirList {
		if strings.Contains(j, "venv") || strings.Contains(j, "env") || strings.Contains(j, "__pycache__") || strings.Contains(j, ".git") || strings.Contains(j, ".tox") {
			t.Error("Directory list has ignored directories.")
		}
	}
}

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
	writeRequirements("", "ultralytics/", "./", false)

	reqRead, err := os.Open("requirements.txt")
	if err != nil {
		panic(err)
	}
	defer reqRead.Close()

	testImports := map[string]bool{"tensorflow": true, "coremltools": true, "clip": true,
		"pytest": true, "comet-ml": true, "pycocotools": true, "shapely": true, "nncf": true,
		"thop": true, "onnxsim": true, "onnxruntime": true, "albumentations": true, "ipython": true,
		"x2paddle": true, "dvclive": true, "opencv-python": true, "matplotlib": true, "wandb": true,
		"ncnn": true, "setuptools": true, "SciPy": true, "yt-dlp": true, "psutil": true, "super-gradients": true,
		"torchvision": true, "Pillow": true, "tflite-runtime": true, "tflite-support": true, "seaborn": true,
		"tqdm": true, "lap": true, "requests": true, "numpy": true, "tritonclient": true, "pandas": true}

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
