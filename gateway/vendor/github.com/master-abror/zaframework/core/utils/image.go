package utils

import (
	"fmt"
	"image"
	"image/png" // Tambahkan import ini agar tipe png.CompressionLevel valid
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ==========================================
// KONFIGURASI PRESET (Untuk Controller)
// ==========================================

// ImagePreset mendefinisikan target ukuran & mode
type ImagePreset struct {
	Width   int
	Height  int
	Quality int    // 1-100 (Default 80)
	Mode    string // "fill" (Crop/Square), "fit" (Keep Ratio), "resize" (Stretch)
}

// Preset Ukuran Populer
var (
	ImgSquareSmall  = ImagePreset{200, 200, 80, "fill"}   // Avatar/Thumb
	ImgSquareMedium = ImagePreset{500, 500, 85, "fill"}   // Katalog Product
	ImgIGFeed       = ImagePreset{1080, 1080, 90, "fill"} // Instagram Post
	ImgIGStory      = ImagePreset{1080, 1920, 90, "fill"} // Instagram Story
	ImgHD           = ImagePreset{1280, 720, 85, "fit"}   // HD Landscape
	ImgFullHD       = ImagePreset{1920, 1080, 85, "fit"}  // Full HD
)

// ==========================================
// LOGIC UTAMA
// ==========================================

// ProcessImage memproses gambar yang SUDAH diupload.
func ProcessImage(filePath string, preset ImagePreset) error {

	// 1. Buka File Gambar
	src, err := imaging.Open(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuka gambar: %v", err)
	}

	var dst *image.NRGBA

	// 2. Lakukan Manipulasi sesuai Mode
	switch preset.Mode {
	case "fill":
		// FILL: Crop tengah (Center Crop)
		dst = imaging.Fill(src, preset.Width, preset.Height, imaging.Center, imaging.Lanczos)

	case "fit":
		// FIT: Resize agar muat (Keep Ratio)
		dst = imaging.Fit(src, preset.Width, preset.Height, imaging.Lanczos)

	case "resize":
		// RESIZE: Paksa ukuran (Stretch)
		dst = imaging.Resize(src, preset.Width, preset.Height, imaging.Lanczos)

	default:
		dst = imaging.Fit(src, preset.Width, preset.Height, imaging.Lanczos)
	}

	// 3. Simpan Kembali
	err = saveImage(dst, filePath, preset.Quality)
	if err != nil {
		return fmt.Errorf("gagal menyimpan gambar: %v", err)
	}

	return nil
}

// Internal Helper untuk Save dengan Kompresi
func saveImage(img *image.NRGBA, path string, quality int) error {
	ext := strings.ToLower(filepath.Ext(path))

	if quality <= 0 {
		quality = 80
	}

	switch ext {
	case ".jpg", ".jpeg":
		return imaging.Save(img, path, imaging.JPEGQuality(quality))

	case ".png":
		// PERBAIKAN DISINI:
		// pngCompression(quality) sudah me-return EncodeOption.
		// Jadi tidak perlu dibungkus imaging.PNGCompressionLevel lagi.
		return imaging.Save(img, path, pngCompression(quality))

	default:
		return imaging.Save(img, path)
	}
}

// Konversi Quality (0-100) ke PNG Compression Level
// Kita gunakan casting ke png.CompressionLevel agar tipe datanya sesuai
func pngCompression(q int) imaging.EncodeOption {
	if q > 90 {
		// No Compression (Cepat, File Besar)
		return imaging.PNGCompressionLevel(png.NoCompression)
	} else if q > 70 {
		// Default Compression
		return imaging.PNGCompressionLevel(png.DefaultCompression)
	} else {
		// Best Compression (Lambat, File Kecil)
		return imaging.PNGCompressionLevel(png.BestCompression)
	}
}
