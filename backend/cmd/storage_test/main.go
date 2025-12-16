package main

import (
	"bingdaily/backend/internal/storage"
	"fmt"
	// "github.com/google/uuid"
)

func main() {
	// s := uuid.New().String()
	s := "a6505678-0692-4df6-9334-299606a281a4"

	bucket := "bingdaily-pictures"

	fmt.Printf("%s\n", s)

	strg := storage.InitializeStorage()

	// url, err := strg.GenerateUploadURL(bucket, s)
	// if err != nil {
	// 	panic("Error occured " + err.Error())
	// }

	url, err := strg.GenerateDownloadURL(bucket, s)
	if err != nil {
		panic("Error occured " + err.Error())
	}

	fmt.Printf("Output: %v\n", url)
}
