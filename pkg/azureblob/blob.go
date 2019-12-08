package azureblob

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// Operations for Block Blobs
// https://docs.microsoft.com/en-us/rest/api/storageservices/understanding-block-blobs--append-blobs--and-page-blobs#about-block-blobs
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#BlockBlobURL

// @@@@@ UPLOAD @@@@@

// UploadFileFromLocalToBlockBlob checks file existence before executing azblob.UploadFileToBlockBlob
func UploadFileFromLocalToBlockBlob(ctx context.Context, localFilePath string, containerURL azblob.ContainerURL) error {
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}

	// BlockBlobURL
	fileName := filepath.Base(localFilePath)
	blockBlobURL := containerURL.NewBlockBlobURL(fileName)

	// UploadFileToBlockBlob - High Level API
	// Uses StageBlock (Putblock) operations to concurrently upload a file in chunks to optimize the throughput
	// For files less than 256 MB, it uses Upload (PutBlob) instead to complete the transfer in a single transaction
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blockBlobURL, azblob.UploadToBlockBlobOptions{
		// MAX: BlockBlobMaxUploadBlobBytes
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16,
	})
	return err
}

// <FILE>
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadFileToBlockBlob
// Calls <BUFFER>

// <BUFFER>
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadBufferToBlockBlob
// Maximum size of block blob: 4.75 TB (100 MB X 50,000 blocks)
// Calls - <UPLOAD> / StageBlock-CommitBlockList

// <UPLOAD> ( Simple ) File < 256MB
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#BlockBlobURL.Upload
// Low Level API

// <STREAM>
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadStreamToBlockBlob

// @@@@@ DOWNLOAD @@@@@

// DownloadFullBlockBlobtoLocalFile ...
func DownloadFullBlockBlobtoLocalFile(ctx context.Context, containerURL azblob.ContainerURL, file *os.File) error {
	//  offset, count ( optional ) - set 0 to download entire blob
	options := azblob.DownloadFromBlobOptions{
		RetryReaderOptionsPerBlock: azblob.RetryReaderOptions{MaxRetryRequests: 20},
	}

	// BlockBlobURL
	fileName := file.Name()
	blockBlobURL := containerURL.NewBlockBlobURL(fileName)

	err := azblob.DownloadBlobToFile(ctx, blockBlobURL.BlobURL, 0, 0, file, options)
	return err
}

// <FILE>
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#DownloadBlobToFile

// <BUFFER>
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#DownloadBlobToBuffer

// <DOWNLOAD>
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#BlobURL.Download
// Low Level API

// @@@@@ DELETE @@@@

// DeleteBlockBlob ...
func DeleteBlockBlob(ctx context.Context, containerURL azblob.ContainerURL, fileName string) error {
	// BlockBlobURL
	blockBlobURL := containerURL.NewBlockBlobURL(fileName)
	_, err := blockBlobURL.BlobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	return err
}

// <DELETE>
// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#BlobURL.Delete
// Low Level API

// @@@@@ MOVE @@@@
// TODO
