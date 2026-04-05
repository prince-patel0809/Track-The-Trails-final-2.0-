package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"track-the-trails/config"
	"track-the-trails/utils"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadProfileImage(w http.ResponseWriter, r *http.Request) {

	// ===== ALLOW ONLY POST =====
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ===== CHECK CONTENT TYPE =====
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		http.Error(w, "Send image using form-data (multipart/form-data)", http.StatusBadRequest)
		return
	}

	// ===== GET USER FROM TOKEN =====
	userID, err := utils.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ===== PARSE FORM =====
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// ===== GET IMAGE FILE =====
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image field required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ===== UPLOAD TO CLOUDINARY =====
	result, err := config.CLD.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			Folder: "track-the-trails/profile",
		},
	)

	if err != nil {
		http.Error(w, "Cloudinary upload failed", http.StatusInternalServerError)
		return
	}

	imageURL := result.SecureURL

	// ===== UPDATE DATABASE =====
	_, err = config.DB.Exec(
		"UPDATE Users SET profile_image=$1 WHERE user_id=$2",
		imageURL,
		userID,
	)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// ===== SUCCESS RESPONSE =====
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile image uploaded successfully",
		"image":   imageURL,
	})
}
