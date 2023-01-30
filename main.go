package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var tmpFolderPath = os.TempDir()

func parseArgs() map[string]string {
	args := make(map[string]string)
	flag.Parse()
	for i := 0; i < flag.NArg(); i++ {
		arg := flag.Arg(i)
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			args[parts[0]] = parts[1]
		}
	}
	return args
}

func isEmptyDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	return len(entries) == 0
}

func safeRemove(path string) {
	err := os.Remove(path)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Println("Removal permission denied for path: " + path)
		} else {
			fmt.Println("Error removing file: " + path)
		}
	}
}

func recursiveRemove(path string) {
	if isEmptyDir(path) {
		safeRemove(path)
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Println("Permission denied for path: " + path)
			return
		}
		panic(err)
	}
	for _, entry := range entries {
		p := path + "/" + entry.Name()
		if entry.IsDir() && !isEmptyDir(p) {
			recursiveRemove(path + "/" + entry.Name())
		} else {
			safeRemove(path + "/" + entry.Name())
		}
	}
	// after all entries were removed, remove folder
	safeRemove(path)
}

func getEntriesByPrefix(path string, prefix string) []string {
	res := []string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if prefix == "" {
				res = append(res, entry.Name())
			} else if strings.HasPrefix(entry.Name(), prefix) {
				res = append(res, entry.Name())
			}
		}
	}
	return res
}

func main() {
	fmt.Printf("Starting...\n\n")

	args := parseArgs()
	prefix := args["prefix"]

	entries := getEntriesByPrefix(tmpFolderPath, prefix)

	if len(entries) == 0 {
		fmt.Println("No entries found")
		return
	}

	for _, f := range getEntriesByPrefix(tmpFolderPath, prefix) {
		fmt.Print("Cleaning " + f + "...")
		recursiveRemove(tmpFolderPath + "/" + f)
		fmt.Println("Done")
	}

	fmt.Printf("\nFinished\n")
}
