package contentgenerator

import (
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"jlpt-practice-helper/jsonfileflow/model"
	"log"
	"os"
	"strings"
)

func GetShortTextWithQuiz(kanji []string, vocabulary []string, grammar []string) (*model.ShortTextResponse, error) {
	log.Println("Generating short text with quiz...")
	prompt := getShortTextPrompt(kanji, vocabulary, grammar)
	schema := getResponseSchemaForShortTextResponse()

	response, err := generateText(prompt, schema)
	if err != nil {
		return nil, err
	}

	if len(response.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates found in response")
	}

	shortTextResponse := &model.ShortTextResponse{}

	for _, part := range response.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), shortTextResponse); err != nil {
				return nil, err
			}
		}
	}

	return shortTextResponse, nil
}

func GenerateSongLyrics(vocabulary []string, grammar []string) (*model.SongLyricsResponse, error) {
	log.Println("Generating song lyrics...")
	prompt := getSongLyricsPrompt(vocabulary, grammar)
	schema := getResponseSchemaForSongLyricsResponse()

	response, err := generateText(prompt, schema)
	if err != nil {
		return nil, err
	}

	if len(response.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates found in response")
	}

	songLyricsResponse := &model.SongLyricsResponse{}

	for _, part := range response.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), songLyricsResponse); err != nil {
				return nil, err
			}
		}
	}

	return songLyricsResponse, nil
}

func GenerateKanjiMnemonicImages(kanji []model.Kanji) (*model.KanjiImageResponse, error) {
	log.Println("Generating kanji mnemonic images...")
	prompt := getKanjiImagePrompt(kanji)
	schema := getResponseSchemaForKanjiImageResponse()

	response, err := generateText(prompt, schema)
	if err != nil {
		return nil, err
	}

	if len(response.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates found in response")
	}

	kanjiImageResponse := &model.KanjiImageResponse{}

	for _, part := range response.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), kanjiImageResponse); err != nil {
				return nil, err
			}
		}
	}

	return kanjiImageResponse, nil
}

func GenerateExampleSentencesForGrammar(grammar string, vocabulary []string) (*model.GrammarSentencesResponse, error) {
	log.Println("Generating example sentences for grammar point " + grammar)
	prompt := getGrammarSentencesPrompt(grammar, vocabulary)
	schema := getResponseSchemaForGrammarSentencesResponse()

	response, err := generateText(prompt, schema)
	if err != nil {
		return nil, err
	}

	if len(response.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates found in response")
	}

	grammarSentencesResponse := &model.GrammarSentencesResponse{}

	for _, part := range response.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), grammarSentencesResponse); err != nil {
				return nil, err
			}
		}
	}

	return grammarSentencesResponse, nil
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

func getShortTextPrompt(kanji []string, vocabulary []string, grammar []string) string {
	kanjiList := strings.Join(kanji, ", ")
	vocabularyList := strings.Join(vocabulary, ", ")
	grammarList := strings.Join(grammar, ", ")

	return "You are an AI that generates Japanese text, their english translation and reading comprehension quizzes. You generate fun and interesting texts. Please create a text in Japanese of about 1000 words, based on the following conditions:" +
		"\n1. **Level**: N4 JLPT" +
		"\n2. **Kanji to use**: " + kanjiList +
		"\n3. **Vocabulary to use**: " + vocabularyList +
		"\n4. **Grammar points to use**: " + grammarList +
		"\nThe text should be natural and based on everyday themes. It must be easy to understand and appropriately use the specified kanji, vocabulary, and grammar points. Please include line breaks in the texts, between paragraphs, for readability." +
		"\nThe reading comprehension quiz should include 5 different questions about the text, with multiple answer choices, in easy to understand japanese, at the N4 JLPT level." +
		"\nReturn valid JSON object based on provided schema."
}

func getResponseSchemaForShortTextResponse() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"text_japanese": {
				Type: genai.TypeString,
			},
			"text_english": {
				Type: genai.TypeString,
			},
			"reading_comprehension_quiz_in_japanese": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"question": {
							Type: genai.TypeString,
						},
						"options": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"text": {
										Type: genai.TypeString,
									},
									"isCorrect": {
										Type: genai.TypeBoolean,
									},
								},
							},
						},
					},
					Required: []string{"question", "options"},
				},
			},
		},
		Required: []string{"text_japanese", "text_english", "reading_comprehension_quiz_in_japanese"},
	}
}

func getSongLyricsPrompt(vocabulary []string, grammar []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")
	grammarList := strings.Join(grammar, ", ")

	return "You are an AI that generates Japanese song lyrics. Please create a song with at least 4 stanzas, separated by lines, using the following conditions:" +
		"\n1. **Level**: N4 JLPT" +
		"\n2. **Vocabulary to use**: " + vocabularyList +
		"\n3. **Grammar points to use**: " + grammarList +
		"\nThe lyrics should be natural and based on everyday themes. It must be easy to understand and appropriately use the specified vocabulary and grammar points."
}

func getResponseSchemaForSongLyricsResponse() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"lyrics_japanese": {
				Type: genai.TypeString,
			},
			"lyrics_english_translation": {
				Type: genai.TypeString,
			},
			"song_title": {
				Type: genai.TypeString,
			},
		},
		Required: []string{"lyrics_japanese", "lyrics_english_translation", "song_title"},
	}
}

func getKanjiImagePrompt(kanji []model.Kanji) string {
	basePrompt := "Create visually interesting images to help memorize the meaning of kanji based on their components. Each image should help in remembering the meaning of the kanji by illustrating its components in a wildly creative and playful way. Do not mention the actual japanese radicals in the prompt, only use visual imagery. These will be used as prompts for an AI image generation model." +
		"\n Here is an example for the kanji candy 菓 (which is made up of the components grass, rice paddy and tree):\"A whimsical illustration of candy growing like fruit from a magical tree in a rice field, surrounded by tall grass. The candy is bright and colorful, hanging from green branches. The rice field is structured in traditional paddy style with water reflecting the sky. The grass is swaying in the breeze. Storybook style, bright colors, cheerful mood.\"" +
		"\n Your descriptions can be in any style or mood. The goal is to create memorable images that will help learners remember the meaning of the kanji."
	var prompts []string

	for _, k := range kanji {
		prompt := fmt.Sprintf("Kanji: '%s' - Meaning: '%s' - Components: (%s)", k.Kanji, k.Meaning, k.Components)
		prompts = append(prompts, prompt)
	}

	return basePrompt + strings.Join(prompts, "\n")
}

func getResponseSchemaForKanjiImageResponse() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"kanji_list": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"kanji": {
							Type: genai.TypeString,
						},
						"image_description": {
							Type: genai.TypeString,
						},
					},
					Required: []string{"kanji", "image_description"},
				},
			},
		},
		Required: []string{"kanji_list"},
	}
}

func getGrammarSentencesPrompt(grammarPoint string, vocabularyList []string) string {

	return "You are an AI that generates Japanese sentences. Please create 5 sentences in japanese, with english translation, for the given grammar points, using the specified vocabulary. The sentences should be at the N4 JLPT level.\n" +
		"Grammar Point: " + grammarPoint +
		"\nVocabulary: " + strings.Join(vocabularyList, ", ")
}

func getResponseSchemaForGrammarSentencesResponse() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"grammar_point": {
				Type: genai.TypeString,
			},
			"sentences": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"sentence_japanese": {
							Type: genai.TypeString,
						},
						"sentence_english_translation": {
							Type: genai.TypeString,
						},
					},
					Required: []string{"sentence_japanese", "sentence_english_translation"},
				},
			},
		},
		Required: []string{"grammar_point", "sentences"},
	}
}
