package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	parent_dir := dir[:strings.LastIndex(dir, "/")]
	fmt.Println(parent_dir)
}

func checkExt(ext string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var paths []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".txt" {
			paths = append(paths, path)
		}

		return nil
	})

	for _, path := range paths {
		//fmt.Println(path)
		//data, _ := ioutil.ReadFile(path)
		fmt.Printf("data on path %s: \n", path)
	}
	return paths
}
