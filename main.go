package main

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"os"
	"runtime"
)

func main() {
	//url := "https://www.youtube.com/watch?v=0vgQSJgb_eY"
	//DownloadYouTube(url)
	UploadMusic("gof")
}

func DownloadYouTube(url string) {
	var envCmd *cmd.Cmd
	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	sysType := runtime.GOOS
	if sysType == "windows" {
		// Create Cmd with options
		envCmd = cmd.NewCmdOptions(cmdOptions, "pkg/yt-dlp/yt-dlp.exe", "-f 140", "-P", "./YouTube/", url)
	} else if sysType == "linux" {
		envCmd = cmd.NewCmdOptions(cmdOptions, "pkg/yt-dlp/yt-dlp", "-f 140", "-o ~/YouTube/%(title)s.%(ext)s", url)
	}

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for envCmd.Stdout != nil || envCmd.Stderr != nil {
			select {
			case line, open := <-envCmd.Stdout:
				if !open {
					envCmd.Stdout = nil
					continue
				}
				fmt.Println(line)
			case line, open := <-envCmd.Stderr:
				if !open {
					envCmd.Stderr = nil
					continue
				}
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	<-envCmd.Start()

	// Wait for goroutine to print everything
	<-doneChan
}

func UploadMusic(bucketName string) {
	var envCmd *cmd.Cmd
	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	sysType := runtime.GOOS
	if sysType == "windows" {
		// Create Cmd with options
		envCmd = cmd.NewCmdOptions(cmdOptions, "pkg/transfer/transfer.exe", bucketName, "./YouTube/")
	} else if sysType == "linux" {
		envCmd = cmd.NewCmdOptions(cmdOptions, "pkg/transfer/transfer", bucketName, "/YouTube/")
	}

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for envCmd.Stdout != nil || envCmd.Stderr != nil {
			select {
			case line, open := <-envCmd.Stdout:
				if !open {
					envCmd.Stdout = nil
					continue
				}
				fmt.Println(line)
			case line, open := <-envCmd.Stderr:
				if !open {
					envCmd.Stderr = nil
					continue
				}
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	<-envCmd.Start()

	// Wait for goroutine to print everything
	<-doneChan
}
