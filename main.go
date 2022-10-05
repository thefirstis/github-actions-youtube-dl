package main

import (
	"bytes"
	"fmt"
	"github.com/go-cmd/cmd"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
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

	for _, path := range ans {
		if filepath.Ext(path) == ".mp4" {
			upload(path)
		}
	}

}

func upload(path string) {
	url := "http://techocblog.qicp.vip:12760/upload"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(path)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("file", filepath.Base(path))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
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
