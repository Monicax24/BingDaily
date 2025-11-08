package main

import (
	"bingdaily/backend/internal/storage"
	"fmt"
	// "github.com/google/uuid"
)

func main() {
	// s := uuid.New().String()
	s := "52de8246-4c7e-479b-833b-a7225435b801"
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
