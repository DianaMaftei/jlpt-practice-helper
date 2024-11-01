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

func CallHuggingFaceAPI(prompt string) (string, error) {
	log.Println("Calling Hugging Face API...")

	url := fmt.Sprintf("https://api-inference.huggingface.co/models/black-forest-labs/FLUX.1-schnell")

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

	const maxRetries = 5
	var lastErr error

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating new HTTP request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}

	for i := 0; i < maxRetries; i++ {
		time.Sleep(1 * time.Minute)
		log.Printf("Attempt %d to call API...\n", i+1)

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if err != nil {
				lastErr = err
				log.Println("Error during API call:", err)
			} else {
				lastErr = fmt.Errorf("received non-200 response: %d", resp.StatusCode)
				fmt.Println(resp)
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
	return "", lastErr // Return the last error after all retries
}
