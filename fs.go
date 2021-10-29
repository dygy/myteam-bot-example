package main

import (
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var jiraLink = GoDotEnvVariable("jiralink")

func readFIle(fileName string) string {
	// Open file for reading
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func replaceTasks(text string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			text, "task_", jiraLink,
		), "hotfix_", jiraLink,
	)
}

func WriteFile(filename string, message string) {
	err := os.WriteFile("task/"+strings.ReplaceAll(filename, "/", "_"), []byte(replaceTasks(message)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteAllFiles() {
	var files, _ = os.ReadDir("task")
	for _, f := range files {
		os.Remove("task/" + f.Name())
	}
}
func allFilesData() string {
	var files, _ = os.ReadDir("task")
	var res = ""
	for _, f := range files {
		res += replaceTasks(f.Name()) + " : " + readFIle("task/"+f.Name()) + "\n"
	}
	return res
}

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	println(os.Getenv(key), key)
	return os.Getenv(key)
}
