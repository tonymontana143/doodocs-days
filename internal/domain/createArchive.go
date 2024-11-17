package domain

type CreateArchiveHandler struct{}

var allowedMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/xml": true,
	"image/jpeg":      true,
	"image/png":       true,
}
