package contentgenerator

import (
	"fmt"
	"jlpt-practice-helper/jsonfileflow/model"
	"strings"
)

func GetShortTextPrompt(kanji []string, vocabulary []string, grammar []string) string {
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

func GetSongLyricsPrompt(vocabulary []string, grammar []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")
	grammarList := strings.Join(grammar, ", ")

	return "You are an AI that generates Japanese song lyrics. Please create a song with at least 4 stanzas, separated by lines, using the following conditions:" +
		"\n1. **Level**: N4 JLPT" +
		"\n2. **Vocabulary to use**: " + vocabularyList +
		"\n3. **Grammar points to use**: " + grammarList +
		"\nThe lyrics should be natural and based on everyday themes. It must be easy to understand and appropriately use the specified vocabulary and grammar points."
}

func GetKanjiImagePrompt(kanji []model.Kanji) string {
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

func GetGrammarSentencesPrompt(grammarPoint string, vocabularyList []string) string {
	return "You are an AI that generates Japanese sentences. Please create 5 sentences in japanese, with english translation, for the given grammar points, using the specified vocabulary. The sentences should be at the N4 JLPT level.\n" +
		"Grammar Point: " + grammarPoint +
		"\nVocabulary: " + strings.Join(vocabularyList, ", ")
}

func GetVocabularySentencesPrompt(vocabulary []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")

	return "You are an AI that generates Japanese sentences. Please create 5 simple and easy-to-understand sentences in Japanese, with English translation, for each of the given vocabulary words. The sentences should be at the N4 JLPT level, using basic grammar and vocabulary.\n" +
		"Vocabulary: " + vocabularyList + "\n" +
		"Ensure that the sentences are suitable for beginners and use straightforward language."
}
