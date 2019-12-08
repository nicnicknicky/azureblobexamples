package azureblob

import (
	"log"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// HandleErrors ...
// TODO: Add more serviceCode processing
func HandleErrors(err error) {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				log.Println("Received 409. Container already exists")
				return
			}
		}
		log.Fatal(err)
	}
}
