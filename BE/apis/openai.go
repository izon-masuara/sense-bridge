package apis

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ChatCompletion struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
}

func encodeImage(imagePath string) (string, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	imageData, err := io.ReadAll(imageFile)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}

func ImageAnlyze(imagePath string) string {
	base64Image, err := encodeImage(imagePath)
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", "OPEN_API_KEY"),
	}

	payload := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "Analisa gambar berikut dan kembalikan hanya konteks judul dari gambar ini",
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
						},
					},
				},
			},
		},
		"max_tokens": 300,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from API: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response to JSON: %v", err)
	}

	var completion ChatCompletion

	err = json.Unmarshal([]byte(responseJSON), &completion)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	if len(completion.Choices) > 0 {
		fmt.Println(completion.Choices[0].Message.Content)
	} else {
		fmt.Println("No choices available")
	}

	return completion.Choices[0].Message.Content
}

func GenerateSummary(text string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{"role": "system", "content": "Kamu adalah seorang guru yang mengajar di kelas tunarungu. Kemudian kamu mendapatkan sebuah teks yang sangat panjang sehingga sulit untuk di konversi ke dalam bahasa isyarat. Tugas kamu adalah mengubah teks tersebut agar mudah di pahami sehingga mudah untuk di ubah ke dalam bahasa isyarat. Dan bagi topik menjadi beberapa paragraf sesuai dengan lingkup pembahasannya. Buat teks hanya mengandung huruf kecil dan datar dan hanya memanfaatkan baris baru jika pembahasannya berbeda."},
			{"role": "user", "content": text},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+"OPEN_API_KEY")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok {
		if len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if message, ok := choice["message"].(map[string]interface{}); ok {
					if content, ok := message["content"].(string); ok {
						return content, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("invalid response format")
}
