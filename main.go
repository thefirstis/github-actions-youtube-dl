package main

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/imroc/req/v3"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	url := "https://frp.acgh.top/upload"
	//DownloadYouTube(url)
	//UploadMusic("gof")

	//获取文件夹中的待上传的文件名
	ans := []string{""}
	filepath.Walk("./downloads", func(path string, info fs.FileInfo, err error) error {
		fmt.Println("path ====" + path)
		if !info.IsDir() {
			if filepath.Ext(path) == ".mp4" {
				ans = append(ans, path)
			}
		}
		return nil
	})

	client := req.C()
	for _, path := range ans {
		fmt.Println(path)
		response, _ := client.R().
			SetFile("file", "./"+path).
			SetUploadCallback(func(infos req.UploadInfo) {
				fmt.Printf("%q uploaded %.2f%%\n", infos.FileName, float64(infos.UploadedSize)/float64(infos.FileSize)*100.0)
			}).Post(url)
		fmt.Println(response.String())
	}

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
