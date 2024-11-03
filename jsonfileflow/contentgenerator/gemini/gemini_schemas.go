package gemini

import "github.com/google/generative-ai-go/genai"

func GetResponseSchemaForShortTextResponse() *genai.Schema {
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

func GetResponseSchemaForSongLyricsResponse() *genai.Schema {
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

func GetResponseSchemaForKanjiImageResponse() *genai.Schema {
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

func GetResponseSchemaForGrammarSentencesResponse() *genai.Schema {
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

func GetResponseSchemaForVocabularySentencesResponse() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"vocabulary_sentences": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"vocabulary": {
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
					Required: []string{"vocabulary", "sentences"},
				},
			},
		},
		Required: []string{"vocabulary_sentences"},
	}
}
