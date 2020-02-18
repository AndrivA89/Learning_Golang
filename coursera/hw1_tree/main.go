package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	printTreeRecursive("", out, path, printFiles)
	return nil
}

func printTreeRecursive(tabString string, out io.Writer, path string, printFiles bool) {
	fileInfo, err := os.Open(path)
	defer fileInfo.Close()

	if err != nil {
		log.Fatalf("Could not open %s: %s", path, err.Error())
	}

	fileName := fileInfo.Name()
	files, err := ioutil.ReadDir(fileName)

	if err != nil {
		log.Fatalf("Could not read dir names in %s: %s", path, err.Error())
	}

	var filesMap map[string]os.FileInfo = map[string]os.FileInfo{}
	var notSortFile []string = []string{}

	for _, file := range files {
		notSortFile = append(notSortFile, file.Name())
		filesMap[file.Name()] = file
	}

	sort.Strings(notSortFile)
	var sortFile []os.FileInfo = []os.FileInfo{}

	for _, stringName := range notSortFile {
		sortFile = append(sortFile, filesMap[stringName])
	}

	files = sortFile
	var tempFileList []os.FileInfo = []os.FileInfo{}

	if !printFiles {
		for _, file := range files {
			if file.IsDir() {
				tempFileList = append(tempFileList, file)
			}
		}
		files = tempFileList
	}

	for i, file := range files {
		if file.IsDir() {
			var newTabString string

			if len(files) > i+1 {
				fmt.Fprintf(out, tabString+"├───"+"%s\n", file.Name())
				newTabString = tabString + "│\t"
			} else {
				fmt.Fprintf(out, tabString+"└───"+"%s\n", file.Name())
				newTabString = tabString + "\t"
			}

			newDir := path + "/" + file.Name()
			printTreeRecursive(newTabString, out, newDir, printFiles)

		} else if printFiles {
			if file.Size() > 0 {
				if len(files) > i+1 {
					fmt.Fprintf(out, tabString+"├───%s (%vb)\n", file.Name(), file.Size())
				} else {
					fmt.Fprintf(out, tabString+"└───%s (%vb)\n", file.Name(), file.Size())
				}
			} else {
				if len(files) > i+1 {
					fmt.Fprintf(out, tabString+"├───%s (empty)\n", file.Name())
				} else {
					fmt.Fprintf(out, tabString+"└───%s (empty)\n", file.Name())
				}
			}
		}
	}
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	/*
		path := os.Args[1]
		printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	*/
	err := dirTree(out, "testdata", true)
	if err != nil {
		panic(err.Error())
	}
}
