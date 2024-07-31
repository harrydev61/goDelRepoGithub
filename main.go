package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	githubToken = "xxx"
	username    = "xxx"
	repoFile    = "repoToDel.txt"
)

func main() {
	var wait sync.WaitGroup
	file, err := os.Open(repoFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		repo := strings.TrimSpace(scanner.Text())
		if repo != "" {
			wait.Add(1)
			go func() {
				defer wait.Done()
				deleteRepo(repo)
			}()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	wait.Wait()
	fmt.Println("Completed deleting repositories.")
}

func deleteRepo(repo string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", username, repo)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "token "+githubToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to delete repository %s: %v", repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("Successfully deleted repository: %s\n", repo)
	} else {
		fmt.Printf("Failed to delete repository: %s (Status: %s)\n", repo, resp.Status)
	}
}
