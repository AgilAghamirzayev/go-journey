package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var count int

func searchFiles(rootDir string, filename string, wg *sync.WaitGroup) {

	defer wg.Done()

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		count++
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), filename) {
			fmt.Printf("Found: %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	rootDir := "/" // Replace this with your directory path

	searchFileName := "tasks.json" // Replace this with the file name to search for

	var wg sync.WaitGroup

	files, err := ioutil.ReadDir(rootDir)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			go searchFiles(filepath.Join(rootDir, file.Name()), searchFileName, &wg)
		}

	}

	wg.Wait()
	fmt.Println("All searches completed")
	fmt.Println("Total files searched:", count)
}
