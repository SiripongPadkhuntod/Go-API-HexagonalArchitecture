package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hexagonalarchitecture/internal/core/port"
)

type FileHandler struct {
	storage port.StoragePort
}

func NewFileHandler(storage port.StoragePort) *FileHandler {
	return &FileHandler{storage: storage}
}

// Upload godoc
// @Summary Upload an image
// @Description Upload an image to the object storage.
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image to upload"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/files/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, newErrorResponse(port.ErrCodeBadRequest, "image is required"))
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse(port.ErrCodeInternalServer, "failed to open image"))
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse(port.ErrCodeInternalServer, "failed to read image"))
		return
	}

	// Generate a unique filename to prevent collisions
	ext := ""
	if len(file.Filename) > 4 && file.Filename[len(file.Filename)-4] == '.' {
		ext = file.Filename[len(file.Filename)-4:]
	} else if len(file.Filename) > 5 && file.Filename[len(file.Filename)-5] == '.' {
		ext = file.Filename[len(file.Filename)-5:]
	}
	objectName := uuid.New().String() + ext

	url, err := h.storage.UploadImage(c.Request.Context(), "", objectName, data, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse(port.ErrCodeInternalServer, "failed to upload image"))
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{"url": url})
}
