package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/file"
	"github.com/chanler/prosel/backend/internal/interfaces/http/middleware"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/file"
	"github.com/gin-gonic/gin"
)

type FileService interface {
	Upload(ctx context.Context, req usecase.UploadRequest) (*domain.FileAsset, error)
	AttachToRef(ctx context.Context, fileID string, refType string, refID string) error
	ListFiles(ctx context.Context, filter domain.FileFilter) ([]domain.FileAsset, domain.Pagination, error)
	DeleteFile(ctx context.Context, id string) error
}

type FileHandler struct{ service FileService }

func NewFileHandler(service FileService) *FileHandler { return &FileHandler{service: service} }

func (h *FileHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.POST("/files/upload", h.upload)
	admin.GET("/files", h.list)
	admin.PATCH("/files/:id/ref", h.attachRef)
	admin.DELETE("/files/:id", h.delete)
}

type fileResponse struct {
	ID           string    `json:"id"`
	UploaderID   string    `json:"uploaderId,omitempty"`
	OriginalName string    `json:"originalName"`
	FileName     string    `json:"fileName"`
	StorageType  string    `json:"storageType"`
	ObjectKey    string    `json:"objectKey"`
	PublicURL    string    `json:"publicUrl"`
	MimeType     string    `json:"mimeType"`
	ByteSize     int64     `json:"byteSize"`
	Width        *int      `json:"width,omitempty"`
	Height       *int      `json:"height,omitempty"`
	RefType      string    `json:"refType,omitempty"`
	RefID        string    `json:"refId,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type attachRefRequest struct {
	RefType string `json:"refType" binding:"required"`
	RefID   string `json:"refId" binding:"required"`
}

func (h *FileHandler) upload(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "File is required", nil)
		return
	}
	opened, err := header.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Unable to read file", nil)
		return
	}
	defer opened.Close()
	mimeType := header.Header.Get("Content-Type")
	asset, err := h.service.Upload(c.Request.Context(), usecase.UploadRequest{UploaderID: middleware.CurrentUserID(c), OriginalName: header.Filename, MimeType: mimeType, Size: header.Size, Data: opened, RefType: c.PostForm("refType"), RefID: c.PostForm("refId")})
	if err != nil {
		h.handleFileError(c, err)
		return
	}
	response.OK(c, toFileResponse(asset))
}

func (h *FileHandler) list(c *gin.Context) {
	files, pagination, err := h.service.ListFiles(c.Request.Context(), fileFilter(c))
	if err != nil {
		h.handleFileError(c, err)
		return
	}
	response.OKWithMeta(c, toFileResponses(files), pagination)
}

func (h *FileHandler) attachRef(c *gin.Context) {
	var req attachRefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid file reference request", nil)
		return
	}
	if err := h.service.AttachToRef(c.Request.Context(), c.Param("id"), req.RefType, req.RefID); err != nil {
		h.handleFileError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *FileHandler) delete(c *gin.Context) {
	if err := h.service.DeleteFile(c.Request.Context(), c.Param("id")); err != nil {
		h.handleFileError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *FileHandler) handleFileError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrFileNotFound):
		response.Error(c, http.StatusNotFound, "FILE_NOT_FOUND", "File not found", nil)
	case errors.Is(err, domain.ErrFileTooLarge):
		response.Error(c, http.StatusBadRequest, "FILE_TOO_LARGE", "File is too large", nil)
	case errors.Is(err, domain.ErrInvalidFile):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid file request", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "File request failed", nil)
	}
}

func fileFilter(c *gin.Context) domain.FileFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	filter := domain.FileFilter{Page: page, PerPage: perPage, Search: c.Query("search"), MimeType: c.Query("type")}
	if status := c.Query("status"); status != "" {
		fileStatus := domain.FileStatus(status)
		if fileStatus.Valid() {
			filter.Status = &fileStatus
		}
	}
	return filter
}

func toFileResponses(files []domain.FileAsset) []fileResponse {
	result := make([]fileResponse, 0, len(files))
	for _, file := range files {
		fileCopy := file
		result = append(result, *toFileResponse(&fileCopy))
	}
	return result
}

func toFileResponse(file *domain.FileAsset) *fileResponse {
	return &fileResponse{ID: file.ID, UploaderID: file.UploaderID, OriginalName: file.OriginalName, FileName: file.FileName, StorageType: string(file.StorageType), ObjectKey: file.ObjectKey, PublicURL: file.PublicURL, MimeType: file.MimeType, ByteSize: file.ByteSize, Width: file.Width, Height: file.Height, RefType: file.RefType, RefID: file.RefID, Status: string(file.Status), CreatedAt: file.CreatedAt, UpdatedAt: file.UpdatedAt}
}
