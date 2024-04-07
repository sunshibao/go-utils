package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cloud.google.com/go/bigquery"
)

func main() {
	importJSONTruncate("market-381602", "manage_report", "DDU-H1-2023-10-13_13")
}

func importJSONTruncate(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	relativePath := "./bigQuery/config/market-381602-54a84ede4ff6.json"
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", absPath)

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("bigquery.NewClient: %v", err)
		return err
	}
	defer client.Close()

	gcsRef := bigquery.NewGCSReference("gs://overseas-manage-test/DDU-H1/2023-10-13/DDU-H1-2023-10-13_13.json")
	gcsRef.SourceFormat = bigquery.JSON
	gcsRef.AutoDetect = true
	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteTruncate

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		fmt.Printf("job completed with error: %v", status.Err())
		return err
	}
	fmt.Printf("tableID:%s job success \n", tableID)
	return nil
}
