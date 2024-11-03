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

func GenerateImageFromPrompt(prompt string) (string, error) {
	log.Println("Calling Hugging Face API...")

	models := []string{
		"black-forest-labs/FLUX.1-schnell",
		"stabilityai/stable-diffusion-3.5-large",
		"alvdansen/softserve_anime",
		"glif/90s-anime-art",
		"Shakker-Labs/SD3.5-LoRA-Chinese-Line-Art",
	}

	var lastErr error
	for i, model := range models {
		log.Printf("Attempt %d to call API with model %s...\n", i+1, model)
		time.Sleep(30 * time.Second)

		body, err := callHuggingFaceAPI(model, prompt)
		if err != nil {
			lastErr = err
			log.Println("Error during API call:", err)
			time.Sleep(time.Duration(10*(1<<i)) * time.Second)
			continue
		}

		imageBytes := body
		encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
		log.Println("Successfully received and encoded image from API.")
		return encodedImage, nil
	}

	log.Println("All attempts failed. Last error:", lastErr)
	return "", lastErr
}

func callHuggingFaceAPI(model string, prompt string) ([]byte, error) {
	bearerToken := os.Getenv("HUGGINGFACE_API_KEY")
	if bearerToken == "" {
		log.Println("Bearer token not found in environment variables")
		return nil, fmt.Errorf("bearer token not found in environment variables")
	}

	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model)
	data := map[string]string{"inputs": prompt}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON data: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating new HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during API call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
