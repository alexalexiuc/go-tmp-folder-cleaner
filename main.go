package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
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

func getEntriesByPrefixes(path string, prefix []string) []string {
	res := []string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if len(prefix) == 0 {
				res = append(res, entry.Name())
			} else {
				for _, p := range prefix {
					if strings.HasPrefix(entry.Name(), p) {
						res = append(res, entry.Name())
						continue
					}
				}
			}
		}
	}
	return res
}

func getUserConfirm() {
	var timeout int = 60 // seconds
	fmt.Print("Continue ? (y/n): ")
	var confirm string
	btnPressed := false
	go func() {
		for i := 0; i < timeout; i++ {
			fmt.Printf("\rContinue ? (y/n): %d", timeout-i)
			time.Sleep(1 * time.Second)
			if btnPressed {
				return
			}
		}
		fmt.Println("\rTimeout")
		os.Exit(0)
	}()
	fmt.Scanln(&confirm)
	btnPressed = true
	if confirm != "y" {
		fmt.Println("Aborted")
		os.Exit(0)
	}
}

func main() {
	fmt.Printf("Starting TEMP folder clean\n\n")
	getUserConfirm()
	fmt.Printf("Temp folder: %s\n", tmpFolderPath)

	args := parseArgs()
	prefixStr := args["prefix"]
	prefixes := strings.Split(prefixStr, ",")

	fmt.Println("Searching files by prefixes: " + prefixStr)
	entries := getEntriesByPrefixes(tmpFolderPath, prefixes)

	if len(entries) == 0 {
		fmt.Println("No entries found")
		return
	} else {
		fmt.Printf("Found %d entries\n", len(entries))
	}

	for _, f := range getEntriesByPrefixes(tmpFolderPath, prefixes) {
		fmt.Print("Cleaning " + f + "...")
		recursiveRemove(tmpFolderPath + "/" + f)
		fmt.Println("Done")
	}

	fmt.Printf("\nFinished\n")
}
