package veTos

import (
	"context"
	"fmt"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io/ioutil"
	"negative-screen/config"
	"os"
	"path/filepath"
)

var tosClient *tos.ClientV2

func init() {
	if err := config.Init(""); err != nil {
		panic(err)
	}
	var (
		ak = config.Config.VeTos.Ak
		sk = config.Config.VeTos.Sk
		// endpoint 若没有指定 HTTP 协议（HTTP/HTTPS），默认使用 HTTPS
		endpoint = config.Config.VeTos.Endpoint
		region   = config.Config.VeTos.Region
	)
	var err error
	credential := tos.NewStaticCredentials(ak, sk)
	tosClient, err = tos.NewClientV2(endpoint, tos.WithCredentials(credential), tos.WithRegion(region))
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}

func GetTosClient() *tos.ClientV2 {
	return tosClient
}

// UploadTos Demo
func UploadTos(filePath string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件名字
	fileName := filepath.Base(filePath)
	// 上传对象
	_, err = tosClient.PutObjectV2(context.Background(), &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: "ve-test",
			Key:    fmt.Sprintf("test001/%s", fileName),
		},
		Content: file,
	})
	if err != nil {
		checkErr(err)
		return err
	}
	return nil
}

// DownloadTos Demo
func DownloadTos() error {
	fileName := "性能1.jpeg"
	localDir := "./obsData"
	// 下载对象
	output, err := tosClient.GetObjectV2(context.Background(), &tos.GetObjectV2Input{
		Bucket: "ve-test",
		Key:    fmt.Sprintf("test001/%s", fileName),
	})
	if err != nil {
		checkErr(err)
		return err
	}
	defer output.Content.Close()

	body, err := ioutil.ReadAll(output.Content)
	if err != nil {
		checkErr(err)
		return err
	}
	localFilePath := filepath.Join(localDir, fileName)
	err = ioutil.WriteFile(localFilePath, body, 0644)
	if err != nil {
		checkErr(err)
		return fmt.Errorf("failed to save file %s: %v", localFilePath, err)
	}

	fmt.Println("File saved to:", localFilePath)
	return nil
}

// DelTos Demo
func DelTos() error {
	fileName := "性能1.jpeg"
	// 删除对象
	_, err := tosClient.DeleteObjectV2(context.Background(), &tos.DeleteObjectV2Input{
		Bucket: "ve-test",
		Key:    fmt.Sprintf("test001/%s", fileName),
	})
	if err != nil {
		checkErr(err)
		return err
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		if serverErr, ok := err.(*tos.TosServerError); ok {
			fmt.Println("Error:", serverErr.Error())
			fmt.Println("Request ID:", serverErr.RequestID)
			fmt.Println("Response Status Code:", serverErr.StatusCode)
			fmt.Println("Response Header:", serverErr.Header)
			fmt.Println("Response Err Code:", serverErr.Code)
			fmt.Println("Response Err Msg:", serverErr.Message)
		} else if clientErr, ok := err.(*tos.TosClientError); ok {
			fmt.Println("Error:", clientErr.Error())
			fmt.Println("Client Cause Err:", clientErr.Cause.Error())
		} else {
			fmt.Println("Error:", err)
		}
	}
}
