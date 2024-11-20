package gemini

import (
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"jlpt-practice-helper/jsonfileflow/contentgenerator"
	"jlpt-practice-helper/jsonfileflow/model"
	"log"
	"os"
	"time"
)

func GetShortTextWithQuiz(vocabulary []string, grammar []string) (*model.ShortTextResponse, error) {
	log.Println("Generating short text with quiz...")

	var response *model.ShortTextResponse
	var err error
	maxRetries := 3

	for retries := 0; retries < maxRetries; retries++ {
		prompt := contentgenerator.GetShortTextPrompt(vocabulary, grammar)
		schema := GetResponseSchemaForShortTextResponse()

		response, err = generateAndParseResponse[model.ShortTextResponse](prompt, schema)
		if err != nil {
			return nil, err
		}

		emptyContent := response == nil || response.TextJapanese == "" || response.TextEnglish == ""
		if !emptyContent {
			break
		}

		log.Println("Empty or lacking content detected, retrying in 60 seconds...")
		time.Sleep(60 * time.Second)
	}

	return response, nil
}

func GenerateSongLyrics(vocabulary []string, grammar []string) (*model.SongLyricsResponse, error) {
	log.Println("Generating song lyrics...")
	prompt := contentgenerator.GetSongLyricsPrompt(vocabulary, grammar)
	schema := GetResponseSchemaForSongLyricsResponse()

	return generateAndParseResponse[model.SongLyricsResponse](prompt, schema)
}

func GenerateKanjiMnemonicImages(kanji []model.Kanji) (*model.KanjiImageResponse, error) {
	log.Println("Generating kanji mnemonic images...")
	prompt := contentgenerator.GetKanjiImagePrompt(kanji)
	schema := GetResponseSchemaForKanjiImageResponse()

	return generateAndParseResponse[model.KanjiImageResponse](prompt, schema)
}

func GenerateExampleSentencesForGrammar(grammar string, vocabulary []string) (*model.GrammarSentencesResponse, error) {
	log.Println("Generating example sentences for grammar point " + grammar)

	var response *model.GrammarSentencesResponse
	var err error
	maxRetries := 3

	for retries := 0; retries < maxRetries; retries++ {
		prompt := contentgenerator.GetGrammarSentencesPrompt(grammar, vocabulary)
		schema := GetResponseSchemaForGrammarSentencesResponse()

		response, err = generateAndParseResponse[model.GrammarSentencesResponse](prompt, schema)
		if err != nil {
			return nil, err
		}

		emptyContent := response.Sentences == nil || len(response.Sentences) == 0
		for _, sentence := range response.Sentences {
			if sentence.SentenceJapanese == "" || sentence.SentenceEnglishTranslation == "" {
				emptyContent = true
				break
			}
		}

		if !emptyContent {
			break
		}

		log.Println("Empty or lacking content detected, retrying in 60 seconds...")
		time.Sleep(60 * time.Second)
	}

	return response, nil
}

func GenerateSentencesForVocabulary(vocabulary []string) (*model.VocabularySentencesResponse, error) {
	log.Println("Generating sentences for vocabulary...")

	var response *model.VocabularySentencesResponse
	var err error
	maxRetries := 3

	for retries := 0; retries < maxRetries; retries++ {
		prompt := contentgenerator.GetVocabularySentencesPrompt(vocabulary)
		schema := GetResponseSchemaForVocabularySentencesResponse()

		response, err = generateAndParseResponse[model.VocabularySentencesResponse](prompt, schema)
		if err != nil {
			return nil, err
		}

		emptyContent := false
		if len(response.VocabularySentences) != len(vocabulary) {
			emptyContent = true
		}
		for _, vocabSentence := range response.VocabularySentences {
			if vocabSentence.Sentences == nil || len(vocabSentence.Sentences) == 0 {
				emptyContent = true
				break
			}
		}

		if !emptyContent {
			break
		}

		log.Println("Empty or lacking content detected, retrying in 60 seconds...")
		time.Sleep(60 * time.Second)
	}

	return response, nil
}

func generateAndParseResponse[T any](prompt string, schema *genai.Schema) (*T, error) {
	response, err := generateText(prompt, schema)
	if err != nil {
		return nil, err
	}

	if len(response.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates found in response")
	}

	var result T
	for _, part := range response.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), &result); err != nil {
				return nil, err
			}
		}
	}

	return &result, nil
}

func generateText(prompt string, schema *genai.Schema) (*genai.GenerateContentResponse, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GENAI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	if schema != nil {
		model.ResponseMIMEType = "application/json"
		model.ResponseSchema = schema
	}
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return nil, err
	}

	return resp, nil
}
