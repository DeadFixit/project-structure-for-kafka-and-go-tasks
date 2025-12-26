package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"task2-go-microservice/services"
	"task2-go-microservice/utils"
)

type IntegrationHandler struct {
	logger *utils.Logger
}

func NewIntegrationHandler(logger *utils.Logger) *IntegrationHandler {
	return &IntegrationHandler{logger: logger}
}

func (h *IntegrationHandler) UploadSample(w http.ResponseWriter, r *http.Request) {
	endpoint := getenv("MINIO_ENDPOINT", "localhost:9000")
	accessKey := getenv("MINIO_ACCESS_KEY", "minioadmin")
	secretKey := getenv("MINIO_SECRET_KEY", "minioadmin")
	bucket := getenv("MINIO_BUCKET", "demo-bucket")
	useSSL := getenv("MINIO_USE_SSL", "false") == "true"

	svc, err := services.NewIntegrationService(endpoint, accessKey, secretKey, bucket, useSSL)
	if err != nil {
		h.logger.Println("failed to create minio client", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectName := time.Now().UTC().Format("20060102-150405") + "-sample.txt"
	content := svc.SampleContent(r.UserAgent())
	info, err := svc.UploadSampleObject(ctx, objectName, content)
	if err != nil {
		h.logger.Println("upload failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.AsyncLog("UPLOAD object=" + objectName)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"object\": \"" + info.Key + "\"}"))
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
