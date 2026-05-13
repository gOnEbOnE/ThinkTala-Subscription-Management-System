package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ==========================================================
// KONSTANTA TIPE FILE (PRESET)
// ==========================================================
// Gunakan variabel ini di Controller agar tidak perlu ngetik manual

// FileImages : JPG, JPEG, PNG, WEBP, GIF
var FileImages = []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}

// FileDocs : PDF, Word, Excel, PowerPoint
var FileDocs = []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}

// FilePDF : Khusus PDF
var FilePDF = []string{".pdf"}

// FileExcel : Khusus Excel
var FileExcel = []string{".xls", ".xlsx", ".csv"}

// ==========================================================
// LOGIC UTAMA
// ==========================================================

// UploadFile menangani proses upload, validasi, dan penyimpanan file.
//
// Parameters:
//   - r: Request HTTP
//   - fieldName: Nama key di form-data (misal: "avatar", "document")
//   - destFolder: Folder tujuan penyimpanan (relatif dari root, misal: "public/uploads/avatars")
//   - allowedExts: Slice string ekstensi yang diperbolehkan (gunakan preset di atas)
//   - maxSize: Ukuran maksimal dalam bytes (misal: 2 * 1024 * 1024 untuk 2MB)
//
// Return:
//   - string: Nama file baru yang disimpan (untuk disimpan di database)
//   - error: Error jika gagal
func UploadFile(r *http.Request, fieldName string, destFolder string, allowedExts []string, maxSize int64) (string, error) {

	// 1. Parse Multipart Form (Max memory 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return "", fmt.Errorf("gagal memproses form: %v", err)
	}

	// 2. Ambil file dari request
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", errors.New("file tidak ditemukan atau kosong")
		}
		return "", err
	}
	defer file.Close()

	// 3. Validasi Ukuran File
	if header.Size > maxSize {
		return "", fmt.Errorf("ukuran file terlalu besar (Max: %d MB)", maxSize/(1024*1024))
	}

	// 4. Validasi Ekstensi (Case Insensitive)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedExtension(ext, allowedExts) {
		return "", fmt.Errorf("tipe file tidak diizinkan. Hanya boleh: %v", strings.Join(allowedExts, ", "))
	}

	// 5. Validasi Content-Type (MIME Sniffing - Optional tapi Recommended)
	// Kita baca 512 byte awal untuk memastikan isi file sesuai ekstensinya (Mencegah rename .exe jadi .jpg)
	if err := validateMimeType(file, ext); err != nil {
		return "", err
	}

	// 6. Buat Folder Tujuan jika belum ada
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat folder upload: %v", err)
	}

	// 7. Generate Nama File Unik (UUID like)
	// Format: unix_nano-random.ext (Contoh: 1709221122-a8b2c.jpg)
	newFilename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), RandomString(4), ext)
	dstPath := filepath.Join(destFolder, newFilename)

	// 8. Simpan File
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}
	defer dst.Close()

	// Reset pointer file setelah sniffing di step 5
	file.Seek(0, 0)

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("gagal menyalin file: %v", err)
	}

	// Return nama file saja (bukan full path) agar fleksibel di DB
	return newFilename, nil
}

// Helper: Cek ekstensi di whitelist
func isAllowedExtension(ext string, allowed []string) bool {
	for _, a := range allowed {
		if ext == a {
			return true
		}
	}
	return false
}

// Helper: Cek MIME type asli (Header Sniffing)
func validateMimeType(file multipart.File, ext string) error {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return errors.New("gagal membaca header file")
	}

	contentType := http.DetectContentType(buffer)

	// Validasi sederhana untuk gambar
	// Jika ekstensi gambar tapi isinya bukan image/*, tolak.
	if (ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif") && !strings.HasPrefix(contentType, "image/") {
		return errors.New("file rusak atau bukan gambar valid")
	}

	// Untuk PDF, Excel, dll, validasi MIME agak kompleks karena variasi office,
	// jadi untuk tutorial ini kita skip deep check selain gambar.

	return nil
}
