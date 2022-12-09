package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ahmadrosid/acm/command"
	"github.com/ahmadrosid/acm/gpt3"
	"github.com/joho/godotenv"
)

func printHelp() {
	fmt.Print(`acm: Auto commit your git changes with AI!

Usage: acm [OPTIONS]

Options:
  -v --version    Print version 
  -h --help       Print help information
`)
}

func printVersion() {
	fmt.Println("acm version v0.0.1")
}

func commitChanges(diff string) {
	prompt := strings.Join([]string{
		"git diff HEAD\\^!",
		diff,
		"",
		"# Write a commit message describing the changes and the reasoning behind them",
		"git commit -F- <<EOF",
	}, "\n")
	resp, err := gpt3.RequestCompletion(prompt)
	if err != nil {
		log.Fatalln(err)
	}

	println("The proposed commit message!")
	println(strings.Repeat("-", command.GetCmdWidth()))
	commitMessage := strings.TrimSpace(resp.Choices[0].Text)
	fmt.Println(commitMessage)
	println(strings.Repeat("-", command.GetCmdWidth()))

	var input string
	fmt.Print("Do you want to continue ? (y/n): ")
	fmt.Scanln(&input)
	if input != "y" {
		log.Fatalln("Aborted!")
		return
	}
	_, err = command.CommitChanges(commitMessage)
	if err != nil {
		log.Fatalln(err)
	}
}

func executeAutoCommit() {
	res := command.GetStagedFiles()
	if res == "" {
		status := command.GetGitStatus()
		if status == "" {
			log.Println("No file changes!!!")
		} else {
			var input string
			fmt.Println("File changes exist in your git repo, but you have not added them to git.")
			fmt.Print("Do you want to execute `git add .` ? (y/n): ")
			fmt.Scanln(&input)
			if input != "y" {
				return
			}
			_, err := command.GitAddAll()
			if err != nil {
				log.Fatal(err)
			}
			commitChanges(command.GetStagedFiles())
		}
	} else {
		commitChanges(res)
	}
}

func main() {
	godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY, please add OPENAI_API_KEY to your env.")
	}

	if len(os.Args) == 1 {
		executeAutoCommit()
		return
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		return
	}

	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		printVersion()
	}

}
