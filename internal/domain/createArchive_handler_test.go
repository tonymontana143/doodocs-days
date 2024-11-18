package domain_test

import (
	"bytes"
	"doodocs-days/internal/domain"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockCreateArchiveService struct {
	mock.Mock
}

func (m *MockCreateArchiveService) ValidateAndZipFiles(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	args := m.Called(files)
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}

func TestCreateArchiveHandler_CreateArchive(t *testing.T) {
	mockService := new(MockCreateArchiveService)
	handler := domain.NewCreateArchiveHandler(mockService)

	// Test POST request with valid files
	files := []*multipart.FileHeader{
		{Filename: "test1.txt", Header: map[string][]string{"Content-Type": {"text/plain"}}},
	}

	mockService.On("ValidateAndZipFiles", files).Return(bytes.NewBuffer([]byte("dummy zip content")), nil)

	req := httptest.NewRequest(http.MethodPost, "/create-archive", nil)
	req.MultipartForm = &multipart.Form{
		File: map[string][]*multipart.FileHeader{
			"files[]": files,
		},
	}

	rr := httptest.NewRecorder()
	handler.CreateArchive(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, status)
	}

	mockService.AssertExpectations(t)
}

func TestCreateArchiveHandler_InvalidMethod(t *testing.T) {
	mockService := new(MockCreateArchiveService)
	handler := domain.NewCreateArchiveHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/create-archive", nil)
	rr := httptest.NewRecorder()

	handler.CreateArchive(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("expected status code %v, got %v", http.StatusMethodNotAllowed, status)
	}
}

func TestCreateArchiveHandler_NoFilesProvided(t *testing.T) {
	mockService := new(MockCreateArchiveService)
	handler := domain.NewCreateArchiveHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/create-archive", nil)
	rr := httptest.NewRecorder()

	handler.CreateArchive(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %v, got %v", http.StatusBadRequest, status)
	}
}

func TestCreateArchiveHandler_InvalidFileType(t *testing.T) {
	mockService := new(MockCreateArchiveService)
	handler := domain.NewCreateArchiveHandler(mockService)

	// Simulate invalid file type
	files := []*multipart.FileHeader{
		{Filename: "invalid.exe", Header: map[string][]string{"Content-Type": {"application/octet-stream"}}},
	}

	mockService.On("ValidateAndZipFiles", files).Return(nil, fmt.Errorf("file %s has unsupported MIME type", files[0].Filename))

	req := httptest.NewRequest(http.MethodPost, "/create-archive", nil)
	req.MultipartForm = &multipart.Form{
		File: map[string][]*multipart.FileHeader{
			"files[]": files,
		},
	}

	rr := httptest.NewRecorder()
	handler.CreateArchive(rr, req)

	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("expected status code %v, got %v", http.StatusUnsupportedMediaType, status)
	}

	mockService.AssertExpectations(t)
}

func TestCreateArchiveHandler_ValidationError(t *testing.T) {
	mockService := new(MockCreateArchiveService)
	handler := domain.NewCreateArchiveHandler(mockService)

	// Simulate validation error
	files := []*multipart.FileHeader{
		{Filename: "test1.txt", Header: map[string][]string{"Content-Type": {"text/plain"}}},
	}

	mockService.On("ValidateAndZipFiles", files).Return(nil, fmt.Errorf("validation error"))

	req := httptest.NewRequest(http.MethodPost, "/create-archive", nil)
	req.MultipartForm = &multipart.Form{
		File: map[string][]*multipart.FileHeader{
			"files[]": files,
		},
	}

	rr := httptest.NewRecorder()
	handler.CreateArchive(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %v, got %v", http.StatusBadRequest, status)
	}

	mockService.AssertExpectations(t)
}
