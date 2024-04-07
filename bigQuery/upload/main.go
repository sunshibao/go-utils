package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	uploadFile(ioutil.Discard, "overseas-manage-test", "aaaa_20231005.json")
}

// uploadFile uploads an object.
func uploadFile(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	relativePath := "./bigQuery/config/market-381602-54a84ede4ff6.json"
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", absPath)
	if err != nil {
		fmt.Printf("storage.NewClient: %w", err)
		return err
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("storage.NewClient: %w", err)
		return err
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open("/Users/sunshibao/Desktop/aaaa_20231005.json")
	if err != nil {
		fmt.Printf("os.Open: %w", err)
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		fmt.Printf("io.Copy: %w", err)
		return err
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("Writer.Close: %v", errors.New("访问权限不足或对象已存在"))
		return err
	}
	fmt.Printf("Blob uploaded job success", object)
	return nil
}
