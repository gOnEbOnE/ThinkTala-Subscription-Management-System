package orders

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// generateInvoicePDF renders a dark-themed invoice PDF matching the ThinkNalyze email template.
func generateInvoicePDF(rec *OrderRecord) ([]byte, string, error) {
	clientName := strings.TrimSpace(rec.ClientName)
	if clientName == "" || clientName == "Unknown" {
		clientName = "Pelanggan"
	}
	packageName := strings.TrimSpace(rec.PackageName)
	if packageName == "" {
		packageName = "-"
	}
	invoiceNumber := strings.TrimSpace(rec.InvoiceNumber)
	if invoiceNumber == "" {
		invoiceNumber = "-"
	}
	paymentStatus := formatPaymentStatus(rec.Status)
	orderDate := formatInvoiceDate(rec.CreatedAt)
	totalPrice := formatRupiahPDF(rec.TotalPrice)
	paymentMethod := strings.TrimSpace(rec.PaymentMethod)
	if paymentMethod == "" {
		paymentMethod = "-"
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()

	pageW, pageH := pdf.GetPageSize()

	// Page background #0F1923
	pdf.SetFillColor(15, 25, 35)
	pdf.Rect(0, 0, pageW, pageH, "F")

	const (
		cardW    = 148.0 // ~560px
		pad      = 8.5   // ~32px
		cardY    = 15.0
		innerPad = 6.0
	)
	cardXPos := (pageW - cardW) / 2
	innerX := cardXPos + pad
	innerW := cardW - pad*2

	// --- Measure content height first (cursor simulation) ---
	headerH := 32.0
	bodyTop := cardY + headerH
	y := bodyTop + pad

	y += 6  // greeting line
	y += 14 // intro paragraph
	y += 5  // section title
	y += 24 // package card
	y += 6  // section title
	invoiceBoxH := 7*9.0 + 6
	y += invoiceBoxH
	y += 6  // gap
	noteH := 20.0
	y += noteH
	y += pad
	footerH := 18.0
	y += footerH
	cardH := y - cardY + pad*0.5

	// --- Card shell ---
	pdf.SetFillColor(26, 37, 53)
	pdf.SetDrawColor(42, 58, 80)
	pdf.SetLineWidth(0.3)
	pdf.RoundedRect(cardXPos, cardY, cardW, cardH, 4.2, "FD", "1234")

	// --- Header gradient (#0D2137 → #0A3D62) ---
	gradientSteps := 24
	stepH := headerH / float64(gradientSteps)
	for i := 0; i < gradientSteps; i++ {
		t := float64(i) / float64(gradientSteps-1)
		r := lerp(13, 10, t)
		g := lerp(33, 61, t)
		b := lerp(55, 98, t)
		pdf.SetFillColor(r, g, b)
		pdf.Rect(cardXPos, cardY+float64(i)*stepH, cardW, stepH+0.2, "F")
	}
	pdf.SetDrawColor(30, 64, 96)
	pdf.Line(cardXPos, cardY+headerH, cardXPos+cardW, cardY+headerH)

	// Logo: Think (cyan) + Nalyze (white)
	pdf.SetFont("Helvetica", "B", 15)
	thinkW := pdf.GetStringWidth("Think")
	nalyzeW := pdf.GetStringWidth("Nalyze")
	logoX := cardXPos + (cardW-thinkW-nalyzeW)/2
	logoY := cardY + 9
	pdf.SetXY(logoX, logoY)
	pdf.SetTextColor(0, 212, 255)
	pdf.Cell(thinkW, 7, "Think")
	pdf.SetTextColor(255, 255, 255)
	pdf.Cell(nalyzeW, 7, "Nalyze")

	// Badge: ✓ Pembayaran Berhasil
	badgeLabel := "Pembayaran Berhasil"
	pdf.SetFont("Helvetica", "B", 7.5)
	badgeTextW := pdf.GetStringWidth(badgeLabel)
	checkSize := 3.0
	badgeW := badgeTextW + checkSize + 12
	badgeH := 6.5
	badgeX := cardXPos + (cardW-badgeW)/2
	badgeY := cardY + 20
	pdf.SetFillColor(10, 61, 42)
	pdf.SetDrawColor(26, 102, 68)
	pdf.RoundedRect(badgeX, badgeY, badgeW, badgeH, 3.2, "FD", "1234")
	drawCheckMark(pdf, badgeX+5, badgeY+1.5, checkSize)
	pdf.SetTextColor(46, 232, 154)
	pdf.SetXY(badgeX+5+checkSize+2, badgeY+1.6)
	pdf.Cell(badgeTextW, 4, badgeLabel)

	// --- Body ---
	y = bodyTop + pad

	// Greeting
	pdf.SetXY(innerX, y)
	pdf.SetFont("Helvetica", "", 9.5)
	pdf.SetTextColor(160, 180, 200)
	pdf.Write(5, "Halo, ")
	pdf.SetFont("Helvetica", "B", 9.5)
	pdf.SetTextColor(232, 240, 248)
	pdf.Write(5, clientName)
	y += 7

	// Intro
	pdf.SetXY(innerX, y)
	pdf.SetFont("Helvetica", "", 8.5)
	pdf.SetTextColor(122, 154, 184)
	pdf.MultiCell(innerW, 4.2,
		"Terima kasih! Subscription Trader kamu di ThinkNalyze telah aktif. Berikut detail transaksi kamu.",
		"", "L", false)
	y = pdf.GetY() + 5

	// Section: Paket Berlangganan
	drawSectionTitle(pdf, innerX, y, innerW, "Paket Berlangganan")
	y = pdf.GetY() + 3

	packageCardH := 22.0
	drawPackageCard(pdf, innerX, y, innerW, packageCardH, packageName)
	y += packageCardH + 6

	// Section: Detail Invoice
	drawSectionTitle(pdf, innerX, y, innerW, "Detail Invoice")
	y = pdf.GetY() + 3

	boxY := y
	pdf.SetFillColor(15, 25, 35)
	pdf.SetDrawColor(42, 58, 80)
	pdf.RoundedRect(innerX, boxY, innerW, invoiceBoxH, 2.6, "FD", "1234")

	rows := []struct {
		label string
		value string
		color [3]int
	}{
		{"No. Invoice", invoiceNumber, [3]int{0, 212, 255}},
		{"Tanggal", orderDate, [3]int{200, 216, 232}},
		{"Nama Klien", clientName, [3]int{200, 216, 232}},
		{"Paket", packageName, [3]int{200, 216, 232}},
		{"Harga", totalPrice, [3]int{200, 216, 232}},
		{"Metode Pembayaran", paymentMethod, [3]int{200, 216, 232}},
		{"Status Pembayaran", paymentStatus, [3]int{46, 232, 154}},
	}
	rowY := boxY + 3
	for i, row := range rows {
		drawInvoiceRow(pdf, innerX+innerPad, rowY, innerW-innerPad*2, row.label, row.value, row.color)
		rowY += 9
		if i < len(rows)-1 {
			pdf.SetDrawColor(30, 46, 64)
			pdf.Line(innerX+4, rowY-1.5, innerX+innerW-4, rowY-1.5)
		}
	}
	y = boxY + invoiceBoxH + 6

	// Note box
	pdf.SetFillColor(13, 30, 48)
	pdf.RoundedRect(innerX, y, innerW, noteH, 2, "F", "12")
	pdf.SetFillColor(0, 212, 255)
	pdf.Rect(innerX, y, 1.0, noteH, "F")
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(122, 154, 184)
	pdf.SetXY(innerX+5, y+4)
	pdf.MultiCell(innerW-10, 3.8,
		"Akses subscription kamu sudah aktif. Mulai analisis pasar dan tingkatkan strategi trading kamu sekarang melalui aplikasi ThinkNalyze.",
		"", "L", false)
	y += noteH + pad

	// Footer
	pdf.SetDrawColor(30, 46, 64)
	pdf.Line(innerX, y, innerX+innerW, y)
	y += 5
	pdf.SetFont("Helvetica", "", 7)
	pdf.SetTextColor(58, 85, 112)
	pdf.SetXY(innerX, y)
	pdf.CellFormat(innerW, 3.5, "Email ini dikirim otomatis oleh sistem ThinkNalyze.", "", 1, "C", false, 0, "")
	pdf.SetX(innerX)
	pdf.CellFormat(innerW, 3.5, "Butuh bantuan? Hubungi support kami di support@thinknalyze.com", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("gagal generate PDF: %w", err)
	}
	return buf.Bytes(), rec.InvoiceNumber, nil
}

func formatInvoiceDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	loc := time.FixedZone("WIB", 7*3600)
	t = t.In(loc)
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()-1], t.Year())
}

func formatRupiahPDF(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	n := len(s)
	if n <= 3 {
		return "Rp " + s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return "Rp " + string(result)
}

func formatPaymentStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "PAID":
		return "Pembayaran Berhasil"
	case "PENDING_PAYMENT":
		return "Menunggu Pembayaran"
	case "CANCELLED":
		return "Dibatalkan"
	case "EXPIRED":
		return "Kedaluwarsa"
	default:
		return status
	}
}

func lerp(a, b int, t float64) int {
	return int(float64(a) + t*float64(b-a))
}

func drawSectionTitle(pdf *gofpdf.Fpdf, x, y, w float64, title string) {
	pdf.SetXY(x, y)
	pdf.SetFont("Helvetica", "B", 6.5)
	pdf.SetTextColor(74, 106, 138)
	pdf.CellFormat(w, 4, strings.ToUpper(title), "", 1, "L", false, 0, "")
}

func drawPackageCard(pdf *gofpdf.Fpdf, x, y, w, h float64, packageName string) {
	// Gradient-like package card background
	gradientSteps := 8
	stepH := h / float64(gradientSteps)
	for i := 0; i < gradientSteps; i++ {
		t := float64(i) / float64(gradientSteps-1)
		r := lerp(10, 13, t)
		g := lerp(42, 30, t)
		b := lerp(64, 48, t)
		pdf.SetFillColor(r, g, b)
		pdf.Rect(x, y+float64(i)*stepH, w, stepH+0.2, "F")
	}
	pdf.SetDrawColor(30, 64, 96)
	pdf.RoundedRect(x, y, w, h, 2.6, "D", "1234")

	iconSize := 11.5
	iconX := x + 5
	iconY := y + (h-iconSize)/2
	drawChartIcon(pdf, iconX, iconY, iconSize, iconSize)

	textX := iconX + iconSize + 5
	pdf.SetXY(textX, y+5)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(232, 244, 255)
	pdf.CellFormat(w-iconSize-14, 5, packageName, "", 1, "L", false, 0, "")
	pdf.SetXY(textX, y+11)
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetTextColor(90, 133, 168)
	pdf.CellFormat(w-iconSize-14, 4, "Trader Subscription \u2014 ThinkNalyze", "", 1, "L", false, 0, "")
}

func drawChartIcon(pdf *gofpdf.Fpdf, x, y, w, h float64) {
	pdf.SetFillColor(0, 58, 92)
	pdf.RoundedRect(x, y, w, h, 2.5, "F", "1234")
	pdf.SetDrawColor(0, 212, 255)
	pdf.SetLineWidth(0.45)
	pdf.Line(x+2, y+h-2.5, x+w*0.35, y+h*0.45)
	pdf.Line(x+w*0.35, y+h*0.45, x+w*0.55, y+h*0.58)
	pdf.Line(x+w*0.55, y+h*0.58, x+w-2, y+2.5)
}

func drawCheckMark(pdf *gofpdf.Fpdf, x, y, size float64) {
	pdf.SetDrawColor(46, 232, 154)
	pdf.SetLineWidth(0.45)
	pdf.Line(x, y+size*0.55, x+size*0.32, y+size*0.9)
	pdf.Line(x+size*0.32, y+size*0.9, x+size, y+size*0.1)
}

func drawInvoiceRow(pdf *gofpdf.Fpdf, x, y, w float64, label, value string, valueColor [3]int) {
	pdf.SetXY(x, y)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(106, 133, 160)
	labelW := w * 0.45
	pdf.CellFormat(labelW, 7, label, "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(valueColor[0], valueColor[1], valueColor[2])
	pdf.CellFormat(w-labelW, 7, value, "", 0, "R", false, 0, "")
}
