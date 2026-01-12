package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
	"se_practice/backend/internal/util"
	"strconv"
	"strings"
	"time"
)

var carouselRepo = repo.NewCarouselRepo()

func GetPublicCarousel(w http.ResponseWriter, r *http.Request) {
	images, err := carouselRepo.GetAll(true)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.OK(w, images)
}

func GetAllAdminCarousel(w http.ResponseWriter, r *http.Request) {
	images, err := carouselRepo.GetAll(false)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.OK(w, images)
}

func CreateCarousel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var img model.CarouselImage
	if err := json.NewDecoder(r.Body).Decode(&img); err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if img.ImageURL == "" {
		util.Error(w, http.StatusBadRequest, "Image URL is required")
		return
	}
	if err := carouselRepo.Create(&img); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.OK(w, img)
}

func UpdateCarousel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var img model.CarouselImage
	if err := json.NewDecoder(r.Body).Decode(&img); err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if img.ID == 0 {
		util.Error(w, http.StatusBadRequest, "ID is required")
		return
	}
	if err := carouselRepo.Update(&img); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.OK(w, nil)
}

func UploadCarouselImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 10MB limit
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		util.Error(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	// Validate extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		util.Error(w, http.StatusBadRequest, "Invalid file type. Only jpg, png, gif, webp allowed")
		return
	}

	// Create directory if not exists
	uploadDir := "./uploads/carousel"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		util.Error(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// Create file
	dst, err := os.Create(savePath)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		util.Error(w, http.StatusInternalServerError, "Failed to copy file content")
		return
	}

	// Return URL
	url := fmt.Sprintf("/uploads/carousel/%s", filename)

	util.OK(w, map[string]string{"url": url})
}

func ReorderCarouselImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var ids []int64
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for i, id := range ids {
		if err := carouselRepo.UpdateSortOrder(id, i); err != nil {
			util.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update sort order for id %d: %v", id, err))
			return
		}
	}

	util.OK(w, nil)
}

func DeleteCarousel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := carouselRepo.Delete(id); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.OK(w, nil)
}
