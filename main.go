package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CommandLineOptions Command line options.
type commandLineOptions struct {
	configFile string
	sourceDir  string
	mysqlUser  string
	mysqlPass  string
}

// Database enrty to import.
type db struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	args := parseCommandLineArgs()

	dump, err := getDumpFiles(args.sourceDir)
	if err != nil {
		fmt.Printf("Could not find db dump files (i.e *.sql.gz) in " + args.sourceDir + ".")
		os.Exit(1)
	}

	db, err := getConfig(args.configFile)
	if err != nil {
		fmt.Printf("Config file error: %v\n", err)
		os.Exit(1)
	}

	importDB(dump, db, args)
}

// Import databases.
func importDB(dump []string, db []db, args *commandLineOptions) {
	for _, d := range db {
		f, err := findDBFile(d.ID, dump)
		if err != nil {
			fmt.Println("No file with ID " + d.ID + " found.")
		} else {
			fmt.Println("Importing " + f + " into " + d.Name)
			// TODO: probably better to use StdoutPipe().
			cmd := "zcat " + f + " | mysql -u" + args.mysqlUser + " -p" + args.mysqlPass + " " + d.Name

			_, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				fmt.Println("Failed to import DB.")
			} else {
				fmt.Println("DB imported successfully.")
			}
		}
	}
}

// Find database dump file.
func findDBFile(id string, dbDump []string) (string, error) {
	for _, s := range dbDump {
		if strings.Contains(s, id) {
			return s, nil
		}
	}

	err := errors.New("database dump file not found")

	return "", err
}

// Get database dump files - *.sql.gz.
func getDumpFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql.gz") {
			res = append(res, filepath.Join(dir, f.Name()))
		}
	}
	return res, nil
}

// Get config file.
func getConfig(fileName string) ([]db, error) {
	var db []db

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return db, err
	}

	json.Unmarshal(file, &db)

	return db, nil
}

// Parse command line arguments.
func parseCommandLineArgs() *commandLineOptions {
	args := &commandLineOptions{}
	dir, _ := os.Getwd()

	defaultConfigFile := dir + "/config.json"
	configFile := flag.String("config-file", defaultConfigFile, "Config file path.")
	sourceDir := flag.String("source-dir", dir, "DB dump source directory.")
	mysqlUser := flag.String("mysql-user", "drupal", "MySQL user name.")
	mysqlPass := flag.String("mysql-pass", "drupal", "MySQL password.")

	flag.Parse()

	args.configFile = *configFile
	args.sourceDir = *sourceDir
	args.mysqlUser = *mysqlUser
	args.mysqlPass = *mysqlPass

	// validate files.
	if _, err := os.Stat(args.configFile); os.IsNotExist(err) {
		log.Fatal("Config file " + args.configFile + " does not exist")
	}

	fi, err := os.Stat(args.sourceDir)
	if err != nil {
		fmt.Printf("Source directory " + args.sourceDir + " does not exist")
		os.Exit(1)
	}
	if !fi.IsDir() {
		fmt.Printf("Make sure " + args.sourceDir + " is a directory")
		os.Exit(1)
	}

	return args
}
