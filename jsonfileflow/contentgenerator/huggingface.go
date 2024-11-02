package contentgenerator

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func GenerateImageForKanji(prompt string) (string, error) {
	log.Println("Calling Hugging Face API...")

	models := []string{
		"black-forest-labs/FLUX.1-schnell",
		"stabilityai/stable-diffusion-3.5-large",
		"alvdansen/softserve_anime",
		"glif/90s-anime-art",
		"Shakker-Labs/SD3.5-LoRA-Chinese-Line-Art",
	}

	bearerToken := os.Getenv("HUGGINGFACE_API_KEY") // Ensure you set this environment variable
	if bearerToken == "" {
		log.Println("Bearer token not found in environment variables")
		return "", fmt.Errorf("bearer token not found in environment variables")
	}

	data := map[string]string{"inputs": prompt}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON data:", err)
		return "", err
	}

	var lastErr error
	client := &http.Client{}

	for i, model := range models {
		url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model)
		log.Printf("Attempt %d to call API with model %s...\n", i+1, model)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error creating new HTTP request:", err)
			return "", err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+bearerToken)

		time.Sleep(1 * time.Minute)

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if err != nil {
				lastErr = err
				log.Println("Error during API call:", err)
			} else {
				lastErr = fmt.Errorf("received non-200 response: %d", resp.StatusCode)
				log.Println(lastErr)
			}
			time.Sleep(time.Duration(10*(1<<i)) * time.Second)
			continue
		}

		defer resp.Body.Close()

		imageBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return "", err
		}

		// Encode the image bytes to base64
		encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
		log.Println("Successfully received and encoded image from API.")
		return encodedImage, nil
	}

	log.Println("All attempts failed. Last error:", lastErr)
	return "", lastErr
}
