package azureblob

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
)

const blobURLTemplate = "https://%s.blob.core.windows.net/%s"

// NewContainerURL ...
func NewContainerURL(pipeline pipeline.Pipeline, accountName, containerName string) azblob.ContainerURL {
	// https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#ContainerURL
	// Create ContainerURL object
	URL, _ := url.Parse(fmt.Sprintf(blobURLTemplate, accountName, containerName))
	containerURL := azblob.NewContainerURL(*URL, pipeline)
	return containerURL
}

// CreateContainer ...
func CreateContainer(ctx context.Context, containerURL azblob.ContainerURL, containerName string) error {
	log.Printf("CreateContainer - %s\n", containerName)
	_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	return err
}

// ListBlobsInContainer ...
func ListBlobsInContainer(ctx context.Context, containerURL azblob.ContainerURL) ([]string, error) {
	// TODO: Review with addition comments
	var blobNames []string
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			return []string{}, err
		}

		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			blobNames = append(blobNames, blobInfo.Name)
		}
	}
	return blobNames, nil
}
