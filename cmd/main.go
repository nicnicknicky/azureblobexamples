package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/joho/godotenv"
	"github.com/nicnicknicky/azureblob/pkg/azureblob"
)

// Based off Azure Blob storage documentation - Blob Storage > Quickstarts > Develop with Blobs
// https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-go?tabs=linux

func main() {
	// Azure Request Pipeline
	azPipeline, azCreds := NewAzureRequestPipeline()

	// Created containers
	tempContainer, permContainer := os.Getenv("EXISTING_TEMP_CONTAINER_NAME"), os.Getenv("EXISTING_PERM_CONTAINER_NAME")
	if len(tempContainer) == 0 || len(permContainer) == 0 {
		log.Fatalf("Created container name has not been set in .env")
	}
	tempContainerURL := azureblob.NewContainerURL(azPipeline, azCreds.accName, tempContainer)
	//
	// permContainerURL := azureblob.NewContainerURL(azPipeline, azCreds.accName, permContainer)

	ctx := context.Background()

	// Upload file to tempContainer
	testFilePath := "/Users/nicholaslim/Desktop/LoremIpsum.pdf"
	err := azureblob.UploadFileFromLocalToBlockBlob(ctx, testFilePath, tempContainerURL)
	azureblob.HandleErrors(err)

	// Create newContainer
	newContainerURL := azureblob.NewContainerURL(azPipeline, azCreds.accName, "newbucket")
	err = azureblob.CreateContainer(ctx, newContainerURL, "newbucket")
	azureblob.HandleErrors(err)

	// Upload temporary file to newContainer
	tempFile, err := ioutil.TempFile("", "temp")
	azureblob.HandleErrors(err)
	defer os.Remove(tempFile.Name())

	data := []byte("Hello temporary blob here!\n")
	_, err = tempFile.Write(data)
	azureblob.HandleErrors(err)
	err = tempFile.Close()
	azureblob.HandleErrors(err)

	err = azureblob.UploadFileFromLocalToBlockBlob(ctx, tempFile.Name(), newContainerURL)
	azureblob.HandleErrors(err)

	// TODO: Copy Blob from Container to Container

	// List blobs in containers
	tempContainerBlobs, err := azureblob.ListBlobsInContainer(ctx, tempContainerURL)
	newContainerBlobs, err := azureblob.ListBlobsInContainer(ctx, newContainerURL)
	azureblob.HandleErrors(err)
	fmt.Println("tempContainerBlobs:", tempContainerBlobs)
	fmt.Println("newContainerBlobs:", newContainerBlobs)

	// PAUSE
	fmt.Println("Press <ENTER> to proceed with cleanup: Delete newContainer and testContainer testBlob.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Delete newContainer
	newContainerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	// Delete random testFile from tempContainer
	azureblob.DeleteBlockBlob(ctx, tempContainerURL, filepath.Base(testFilePath))
}

// AzCredentials ...
type AzCredentials struct {
	accName string
	accKey  string
}

// NewAzureRequestPipeline creates a pipeline using the defined storage account name and account key
func NewAzureRequestPipeline() (pipeline.Pipeline, AzCredentials) {
	// Loading env from .env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error - Unable to load .env file")
	}

	// Setup credentials
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	// Create default request pipeline
	return azblob.NewPipeline(credential, azblob.PipelineOptions{}), AzCredentials{accountName, accountKey}
}

// fileExistsOnLocal checks the existence of a file and returns a bool
func fileExistsOnLocal(filePath string) bool {
	var fileExists bool
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		fileExists = true
	}
	return fileExists
}

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}
