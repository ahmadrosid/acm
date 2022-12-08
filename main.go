package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/ahmadrosid/acm/command"
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
	apiKey := os.Getenv("OPENAI_API_KEY")
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	resp, err := client.CompletionWithEngine(ctx, "code-davinci-002", gpt3.CompletionRequest{
		Prompt: []string{
			"git diff HEAD\\^!",
			diff,
			"",
			"# Write a commit message describing the changes and the reasoning behind them",
			"git commit -F- <<EOF",
		},
		MaxTokens:   gpt3.IntPtr(2000),
		Temperature: gpt3.Float32Ptr(0.0),
		Stop:        []string{"EOF"},
	})
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatalln(err)
	}

	println("You commit message!")
	println(strings.Repeat("-", command.GetCmdWidth()), "\n")
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
	_, err = command.ExecCmd("git", "commit", "-m", commitMessage)
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
