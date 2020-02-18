package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	fileObj, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not open %s: %s", path, err.Error())
	}
	defer fileObj.Close()
	fileName := fileObj.Name()
	files, err := ioutil.ReadDir(fileName)
	if err != nil {
		log.Fatalf("Could not read dir names in %s: %s", path, err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Fprintf(out, "%s\n", file.Name())
			newDir := filepath.Join(path, file.Name())
			dirTree(out, newDir, printFiles)
		} else if printFiles {
			if file.Size() > 0 {
				fmt.Fprintf(out, "%s (%vb)\n", file.Name(), file.Size())
			} else {
				fmt.Fprintf(out, "%s (empty)\n", file.Name())
			}
		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := "testdata" //os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
