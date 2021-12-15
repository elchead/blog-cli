package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/skratchdot/open-golang/open"
)

func OpenPostInBrowser(link string) *exec.Cmd {
	cmd := StartRenderBlog()
	time.Sleep(1 * time.Second)
	OpenBrowser(link)
	fmt.Println("Press Ctrl+c to stop render process")
	return cmd
}

func OpenBrowser(link string) {
	url := "http://localhost:1313/"+link
	err := open.Run(url)
	if err != nil {
		fmt.Println("Could not open browser: ", err)
	}
	fmt.Println("Open ",url)
}

func StartRenderBlog() *exec.Cmd {
	cmd := exec.Command("hugo","serve")//"--disableFastRender"
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if std:=string(output); std!= "" { 
		fmt.Println(std) 
	}
	if err != nil {
		log.Fatal("Could not serve hugo: ",err)
	}
	return cmd
}

func OpenObsidianFile(filename string) {
	err := open.Run(fmt.Sprintf("obsidian://open?file=%s",filename))
	if err != nil {
		log.Printf("Error opening obsidian: %v", err)
	}
}

// Requires https://github.com/elchead/readwise-note-extractor with inclusion in PATH
func PushToReadwise(path string){
	cmd := exec.Command("push_readwise.py")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if std:=string(output); std!= "" { 
		fmt.Println(std) 
	}
	if err != nil {
		log.Fatal("Could not push to readwise: ",err)
	}
}
