package imgprocessing

import (
	"mime/multipart"

	"github.com/haohmaru3000/go_sdk/sdkcm"
)

type ImgProcessing interface {
	// call img processing service to resize img and upload to s3
	Resize(file *multipart.FileHeader, folder string, longEdge int, quality int) (*sdkcm.Image, error)

	ResizeFile(filePath string, folder string, longEdge int, quality int) (*sdkcm.Image, error)
}
