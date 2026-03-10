package utils

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// R2Config holds Cloudflare R2 credentials loaded from env vars.
// Set these in Railway env vars:
//
//	R2_ACCOUNT_ID      — Cloudflare Account ID
//	R2_ACCESS_KEY_ID   — R2 Access Key ID
//	R2_SECRET_KEY      — R2 Secret Access Key
//	R2_BUCKET          — Bucket name (e.g. "thinktala-uploads")
//	R2_PUBLIC_URL      — Public URL of the bucket (e.g. "https://pub-xxxx.r2.dev")
func r2Client() (*minio.Client, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_KEY")

	if accountID == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("R2 env vars not set")
	}

	endpoint := fmt.Sprintf("%s.r2.cloudflarestorage.com", accountID)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create R2 client: %v", err)
	}
	return client, nil
}

// IsR2Enabled returns true when all required R2 env vars are set.
func IsR2Enabled() bool {
	return os.Getenv("R2_ACCOUNT_ID") != "" &&
		os.Getenv("R2_ACCESS_KEY_ID") != "" &&
		os.Getenv("R2_SECRET_KEY") != "" &&
		os.Getenv("R2_BUCKET") != ""
}

// UploadFileToR2 uploads a multipart file from an HTTP request to Cloudflare R2.
// Returns the public URL of the uploaded file.
//
// Falls back to local disk upload (UploadFile) if R2 env vars are not set.
func UploadFileToR2(r *http.Request, fieldName string, folder string, allowedExts []string, maxSize int64) (string, error) {
	// Fall back to local upload if R2 is not configured
	if !IsR2Enabled() {
		return UploadFile(r, fieldName, folder, allowedExts, maxSize)
	}

	// 1. Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return "", fmt.Errorf("gagal memproses form: %v", err)
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("file tidak ditemukan: %v", err)
	}
	defer file.Close()

	// 2. Validate size
	if header.Size > maxSize {
		return "", fmt.Errorf("ukuran file terlalu besar (Max: %d MB)", maxSize/(1024*1024))
	}

	// 3. Validate extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedExtension(ext, allowedExts) {
		return "", fmt.Errorf("tipe file tidak diizinkan. Hanya boleh: %v", strings.Join(allowedExts, ", "))
	}

	// 4. Validate MIME type
	if err := validateMimeType(file, ext); err != nil {
		return "", err
	}
	file.Seek(0, 0)

	// 5. Build object key: folder/timestamp-rand.ext
	objectName := fmt.Sprintf("%s/%d-%s%s", folder, time.Now().UnixNano(), RandomString(6), ext)

	// 6. Upload to R2
	client, err := r2Client()
	if err != nil {
		return "", err
	}

	bucket := os.Getenv("R2_BUCKET")
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = client.PutObject(ctx, bucket, objectName, file, header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("gagal upload ke R2: %v", err)
	}

	// 7. Return public URL
	publicBase := strings.TrimRight(os.Getenv("R2_PUBLIC_URL"), "/")
	if publicBase == "" {
		// Fallback: construct from account ID (works if bucket has public access enabled)
		publicBase = fmt.Sprintf("https://%s.r2.dev", os.Getenv("R2_BUCKET"))
	}

	return fmt.Sprintf("%s/%s", publicBase, objectName), nil
}
