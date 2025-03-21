package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	filePath := "http://localhost:3000/69382900-86ea-4c4e-8694-2ee008c150ea" // Pastikan file ini benar-benar JPG/JPEG/PNG
	url := "http://localhost:3000/send/image"

	// Buka file gambar
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error membuka file:", err)
		return
	}
	defer file.Close()

	// Cek ekstensi file
	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println("Error membaca informasi file:", err)
		return
	}

	// Validasi ekstensi file
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := fileStat.Name()[len(fileStat.Name())-4:] // Ambil 4 karakter terakhir
	if !validExtensions[ext] {
		fmt.Println("Error: Format file tidak didukung. Gunakan JPG, JPEG, atau PNG.")
		return
	}

	// Buat buffer dan multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Tambahkan field lainnya
	_ = writer.WriteField("phone", "6289685028129@s.whatsapp.net")
	_ = writer.WriteField("caption", "Selamat malam")
	_ = writer.WriteField("view_once", "false")
	_ = writer.WriteField("compress", "false")

	// Tambahkan file gambar
	part, err := writer.CreateFormFile("image", fileStat.Name())
	if err != nil {
		fmt.Println("Error menambahkan file:", err)
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error menyalin file:", err)
		return
	}

	// Tutup writer untuk mengakhiri multipart form-data
	writer.Close()

	// Buat request POST
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error membuat request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Kirim request ke server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error mengirim request:", err)
		return
	}
	defer resp.Body.Close()

	// Baca respon dari server
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error membaca response body:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(respBody))
}
