package domain_test

import (
	"bytes"
	"doodocs-days/internal/domain"
	"doodocs-days/internal/models"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockArchiveInfoService struct {
	mock.Mock
}

func (m *MockArchiveInfoService) GetZipFileInfo(file multipart.File, handler *multipart.FileHeader) (*models.ArchiveInfo, error) {
	args := m.Called(file, handler)
	return args.Get(0).(*models.ArchiveInfo), args.Error(1)
}

func (m *MockArchiveInfoService) IsZipFile(file multipart.File) (bool, error) {
	args := m.Called(file)
	return args.Bool(0), args.Error(1)
}

func TestArchiveInfoHandle(t *testing.T) {
	mockService := new(MockArchiveInfoService)
	mockService.On("IsZipFile", mock.Anything).Return(true, nil)
	mockService.On("GetZipFileInfo", mock.Anything, mock.Anything).Return(&models.ArchiveInfo{
		Filename:     "example.zip",
		Archive_size: 1024,
		Total_size:   2048,
		Total_files:  5,
		Files: []models.ObjectFile{
			{FilePath: "file1.txt", Size: 100, MimeType: "text/plain"},
			{FilePath: "file2.txt", Size: 200, MimeType: "text/plain"},
		},
	}, nil)

	handler := domain.NewFileHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/archive", bytes.NewReader([]byte("dummy data")))
	req.Header.Set("Content-Type", "multipart/form-data")
	rr := httptest.NewRecorder()

	handler.ArchiveInfoHandle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("got %v, want %v", status, http.StatusOK)
	}

	expected := `{"Filename":"example.zip","Archive_size":1024,"Total_size":2048,"Total_files":5,"Files":[{"FilePath":"file1.txt","Size":100,"MimeType":"text/plain"},{"FilePath":"file2.txt","Size":200,"MimeType":"text/plain"}]}`
	if rr.Body.String() != expected {
		t.Errorf("got %v, want %v", rr.Body.String(), expected)
	}

	mockService.AssertExpectations(t)
}

func TestArchiveInfoHandleInvalidMethod(t *testing.T) {
	mockService := new(MockArchiveInfoService)
	handler := domain.NewFileHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/archive", nil)
	rr := httptest.NewRecorder()

	handler.ArchiveInfoHandle(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestArchiveInfoHandleMissingFile(t *testing.T) {
	mockService := new(MockArchiveInfoService)
	handler := domain.NewFileHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/archive", nil)
	rr := httptest.NewRecorder()

	handler.ArchiveInfoHandle(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestArchiveInfoHandleInvalidZip(t *testing.T) {
	mockService := new(MockArchiveInfoService)
	mockService.On("IsZipFile", mock.Anything).Return(false, nil)

	handler := domain.NewFileHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/archive", bytes.NewReader([]byte("dummy data")))
	req.Header.Set("Content-Type", "multipart/form-data")
	rr := httptest.NewRecorder()

	handler.ArchiveInfoHandle(rr, req)

	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("got %v, want %v", status, http.StatusUnsupportedMediaType)
	}

	mockService.AssertExpectations(t)
}
