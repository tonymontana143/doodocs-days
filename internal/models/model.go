package models

type Response struct {
	filename     string
	archive_size float64
	total_size   float64
	total_file   float64
	files        []ObjectFile
}
type ObjectFile struct {
	filePath string
	size     float64
	mimeType string
}
