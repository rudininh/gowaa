

package main

import (
    "encoding/csv"
    "bytes"
    "fmt"
    "os"
    "net/http"
    "io"
    "encoding/json"
    "time"
    "bufio"
    "math/rand"
    "mime/multipart"
)

type Message struct {
    Phone   string `json:"phone"`
    Message string `json:"message"`
}

func sendTextMessage(phone string, fixedMessage string, successLog, failedLog *bufio.Writer) {
    msg := Message{Phone: phone, Message: fixedMessage}
    jsonData, _ := json.Marshal(msg)

    resp, err := http.Post("http://localhost:3000/send/message",
        "application/json", bytes.NewBuffer(jsonData))

    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logMsg := ""

    if err != nil {
        logMsg = fmt.Sprintf("[%s] Gagal mengirim teks ke %s: %v\n\n", timestamp, phone, err)
        fmt.Print(logMsg)
        failedLog.WriteString(logMsg)
        failedLog.Flush()
        return
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)

    if resp.StatusCode == http.StatusOK {
        logMsg = fmt.Sprintf("[%s] Pesan teks ke %s berhasil! Response: %s\n\n", timestamp, phone, string(respBody))
        fmt.Print(logMsg)
        successLog.WriteString(logMsg)
    } else {
        logMsg = fmt.Sprintf("[%s] Gagal mengirim teks ke %s! Response: %s\n\n", timestamp, phone, string(respBody))
        fmt.Print(logMsg)
        failedLog.WriteString(logMsg)
    }

    successLog.Flush()
    failedLog.Flush()
}

func sendImageMessageFromURL(phone string, imageURL string, successLog, failedLog *bufio.Writer) {
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    _ = writer.WriteField("phone", phone)
    _ = writer.WriteField("image_url", imageURL)
    writer.Close()

    req, err := http.NewRequest("POST", "http://localhost:3000/send/image", body)
    if err != nil {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logMsg := fmt.Sprintf("[%s] Gagal membuat request: %v\n\n", timestamp, err)
        fmt.Print(logMsg)
        failedLog.WriteString(logMsg)
        failedLog.Flush()
        return
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        logMsg := fmt.Sprintf("[%s] Gagal mengirim gambar ke %s: %v\n\n", timestamp, phone, err)
        fmt.Print(logMsg)
        failedLog.WriteString(logMsg)
        failedLog.Flush()
        return
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    if resp.StatusCode == http.StatusOK {
        logMsg := fmt.Sprintf("[%s] Gambar ke %s berhasil dikirim! Response: %s\n\n", timestamp, phone, string(respBody))
        fmt.Print(logMsg)
        successLog.WriteString(logMsg)
    } else {
        logMsg := fmt.Sprintf("[%s] Gagal mengirim gambar ke %s! Response: %s\n\n", timestamp, phone, string(respBody))
        fmt.Print(logMsg)
        failedLog.WriteString(logMsg)
    }

    successLog.Flush()
    failedLog.Flush()
}

func main() {
    file, err := os.Open("messages.csv")
    if err != nil {
        fmt.Println("Error membuka file:", err)
        return
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error membaca CSV:", err)
        return
    }

    successFile, err := os.OpenFile("success.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error membuka success.log:", err)
        return
    }
    defer successFile.Close()
    successLog := bufio.NewWriter(successFile)

    failedFile, err := os.OpenFile("failed.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error membuka failed.log:", err)
        return
    }
    defer failedFile.Close()
    failedLog := bufio.NewWriter(failedFile)

    fixedMessage := `BUKU SAKU TAPERA

BADAN KEPEGAWAIAN DAERAH DIKLAT KOTA BANJARMASIN
https://bkd.banjarmasinkota.go.id

Informasi
Bagi peserta Tapera yang belum melakukan pemutakhiran data agar dapat segera melakukan updating data Tapera.

Buku Panduan Pemutakhiran Data PNS
https://drive.google.com/file/d/1HfDteZlsp14ASf5JgfwQR9QOc6LCrr9v/view?usp=sharing

Demikian disampaikan, terima kasih.

Mohon untuk tidak menghubungi/membalas WA nomor ini.`

    for _, record := range records {
        if len(record) < 1 {
            continue
        }

        phone := record[0]
        sendTextMessage(phone, fixedMessage, successLog, failedLog)

        imageURL := "https://app.banjarmasinkota.go.id/sidinketik/img/tapera.png" // Gantilah dengan URL gambar yang sesuai
        sendImageMessageFromURL(phone, imageURL, successLog, failedLog)

        delay := 3 + rand.Intn(8)
        fmt.Printf("Menunggu %d detik sebelum mengirim pesan berikutnya...\n\n", delay)
        time.Sleep(time.Duration(delay) * time.Second)
    }

    successLog.Flush()
    failedLog.Flush()
}
