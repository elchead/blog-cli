package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/elchead/blog-cli/blog"
)

func getInput(read io.Reader,inputQuestion string) bool {
	fmt.Print(inputQuestion)
	reader := bufio.NewReader(read)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSuffix(answer, "!\n")
	desiredInput := "y" 
	if answer != desiredInput {
		fmt.Println("Answer was:", answer)
	}
	return answer == desiredInput		
}

func readMetadata(title string) blog.Metadata {
	fmt.Printf("Create new post %s\n",title)
	fmt.Print("Enter category: ")
	reader := bufio.NewReader(os.Stdin)
	category, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occured while reading input. Please try again", err)
	}
	category = strings.TrimSuffix(category, "\n")
	return blog.Metadata{Title: title, Categories : []string{category}, Date: time.Now().Format("2006-01-02")}
}

func createBlog() blog.BlogWriter {
	templateFile, err := os.Open(bookTemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	return blog.BlogWriter{RepoPath:repoDir,WritingDir: writingDir,FS:filesystem,BookDir:bookDir,BookTemplate:templateFile}
}

// instantiate and draft Post
func createAndWritePost(title string,isBook bool) blog.Post {
	var post blog.Post
	var err error
	if !isBook {
		meta := readMetadata(title)
		post,err = blogWriter.DraftArticle(meta)
	} else {
		bookmeta := blog.Metadata{Title: title}
		post, err = blogWriter.DraftBook(bookmeta)
		fmt.Printf("Created new book note %s\n",title)
	}
	if err != nil {
		log.Fatal(err)
	}
	return post
}

// only instantiate Post
func newPost(title string,isBook bool) blog.Post {
	var post blog.Post
	meta := blog.Metadata{Title:title}
	if isBook {
		post = blog.NewBook(meta, blog.GetFilepath(title,bookDir))
	} else {
		post = blog.NewArticle(meta,blog.GetFilepath(title,writingDir))
	}
	return post
}


// not needed at the moment
func startGoRoutine(exitChan chan os.Signal, done chan bool) {
	go func() {
		cmd := StartRenderBlog()		
		for s := range exitChan {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,os.Interrupt:
				done <- true
				return
			}
		}
		cmd.Wait()
	}()
}
