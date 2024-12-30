package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

var FromEmail = "youremail@gmail.com"
var FromPass = "yourapppassword"

func SaveUploadedFile(r *http.Request, formFileName string) (string, error) {
	file, header, err := r.FormFile(formFileName)
	if err != nil {
		return "", nil // no file uploaded
	}
	defer file.Close()

	os.MkdirAll("uploads", 0755)
	filePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("error saving file: %v", err)
	}
	return filePath, nil
}

func SaveUploadedFiles(r *http.Request, formFileName string) ([]string, error) {
	var filePaths []string
	files := r.MultipartForm.File[formFileName]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		os.MkdirAll("uploads", 0755)
		filePath := filepath.Join("uploads", fileHeader.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %v", err)
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			return nil, fmt.Errorf("error saving file: %v", err)
		}
		filePaths = append(filePaths, filePath)
	}
	return filePaths, nil
}

func AttachFile(w *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	partHeader.Set("Content-Type", "application/octet-stream")
	part, err := w.CreatePart(partHeader)
	if err != nil {
		return fmt.Errorf("error creating part: %v", err)
	}
	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(contentBytes)))
	base64.StdEncoding.Encode(encoded, contentBytes)
	if _, err := part.Write(encoded); err != nil {
		return fmt.Errorf("error writing part: %v", err)
	}
	return nil
}

func SplitEmails(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	var res []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
