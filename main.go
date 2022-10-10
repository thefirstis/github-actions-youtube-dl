package main

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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
			uploadInfo(path)
		}
	}

}

func uploadInfo(path string) {
	//req.DevMode()
	id := uuid.NewString()
	info, _ := os.Stat(path)
	fmt.Println(info.Size())
	var limit int64 = 1024 * 1024 * 10
	blocked := (info.Size() / limit) + 1
	fmt.Println(blocked)
	req.R().SetBody(map[string]interface{}{ // Set form data while uploading
		"id":    id,
		"name":  filepath.Base(path),
		"block": blocked,
	}).Post("http://techocblog.qicp.vip:12760/uploadInfo")

	//开始分片上传
	upload(path, id)
}

func upload(path string, id string) {
	url := "http://techocblog.qicp.vip:12760/uploads"
	file, _ := os.Open(path)
	//关闭文件
	defer file.Close()
	count := 0
	for {
		//读取文件内容
		buf := make([]byte, 1024*1024*50) //50M大小
		count += 1
		name := id + "." + strconv.Itoa(count)
		//n表示从文件读取内容的长度，buf为切片类型
		n, err := file.Read(buf)
		//文件出错同时没有到结尾
		if err != nil && err != io.EOF {
			fmt.Println("err = ", err)
			return
		}
		if n == 0 {
			break
		}
		_, err = req.R().
			SetFileBytes("file", name, buf[:n]).
			SetUploadCallback(func(info req.UploadInfo) {
				fmt.Printf("%q uploaded\n", info.FileName)
			}).
			Post(url)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("file %s uploaded\n", path)
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
