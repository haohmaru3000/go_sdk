package plugin

type Uploader interface {
	UploadFile(filePath string, opt map[string]interface{}) error
}
