package models

type ArchiveInfo struct {
	Filename     string       `json:"filename"`
	Archive_size float32      `json:"archive_size"`
	Total_size   float32      `json:"total_size"`
	Total_files  float32      `json:"total_files"`
	Files        []ObjectFile `json:"files"`
}
type ObjectFile struct {
	FilePath string  `json:"file_path"`
	Size     float64 `json:"size"`
	MimeType string  `json:"mime_type"`
}
type Mail struct {
	SMTPServer string
	Username   string
	Password   string
}
