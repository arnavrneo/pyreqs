/*
Copyright © 2023 ARNAV <r.arnav@icloud.com>
*/

package clean

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"pyreqs/cmd/create"
	"pyreqs/cmd/utils"
	"strings"
)

var reqPath string

func cleanReq(reqPath string, dirPath string, venvPath string, ignoreDirs string, printReq bool) {

	// TODO: better way of showing errors
	if reqPath == " " {
		panic("Path to requirements.txt not passed!")
	}

	pyFiles, dirList := create.GetPaths(dirPath, ignoreDirs)
	imports := create.ReadImports(pyFiles, dirList)
	distPackages := create.GetLocalPackages(venvPath)

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

	var pypiStore []string
	for i := range imports {
		cntr := 0
		for j := range localSet {
			if i == j {
				cntr += 1
			}
		}
		if cntr == 0 {
			pypiStore = append(pypiStore, i)
		}
	}
	pypiSet := create.FetchPyPIServer(pypiStore)

	importsInfo := make(map[string]bool)
	for i := range localSet {
		importsInfo[i] = true
	}
	for j := range pypiSet {
		importsInfo[j] = true
	}

	// read and write
	reqFile, err := os.OpenFile(reqPath, os.O_RDWR, 0644)
	utils.Check(err)
	scanner := bufio.NewScanner(reqFile)
	var bs []byte
	buf := bytes.NewBuffer(bs) // stores all the user imports

	for scanner.Scan() {
		text := scanner.Text()
		if importsInfo[strings.Split(text, "==")[0]] {
			if printReq == true {
				fmt.Println(text)
			}
			_, err = buf.WriteString(scanner.Text() + "\n")
			utils.Check(err)
		}
	}

	err = reqFile.Truncate(0)
	utils.Check(err)
	_, err = reqFile.Seek(0, 0)
	utils.Check(err)
	_, err = buf.WriteTo(reqFile)
	utils.Check(err)
}

// Cmd represents the clean command
var Cmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans up your requirements.txt",
	Long:  `It cleans up your requirements.txt by removing imports that are not imported in your project.`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, err := cmd.Flags().GetString("dirPath")
		utils.Check(err)
		venvPath, err := cmd.Flags().GetString("venvPath")
		utils.Check(err)
		ignoreDirs, err := cmd.Flags().GetString("ignore")
		utils.Check(err)
		printReq, err := cmd.Flags().GetBool("print")
		utils.Check(err)

		cleanReq(reqPath, dirPath, venvPath, ignoreDirs, printReq)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// CleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// CleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	Cmd.Flags().StringVarP(&reqPath, "reqPath", "r", " ", "path to requirements.txt")
}
