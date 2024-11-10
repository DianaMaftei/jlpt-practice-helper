package contentgenerator

import (
	"fmt"
	"jlpt-practice-helper/jsonfileflow/model"
	"strings"
)

func GetShortTextPrompt(vocabulary []string, grammar []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")
	grammarList := strings.Join(grammar, ", ")

	return "You are an AI that generates a Japanese text, its English translation and a reading comprehension quiz, about the text content, with 5 questions and answers. Please create a cohesive and interesting text of any kind (short story, news, diary, blog, advertisement, letter, etc), in Japanese, based on the following conditions:" +
		"\n1. The language difficulty MUST be lower intermediate, at around the N4 JLPT exam level. Do not use vocabulary, grammar or kanji above the N3 level." +
		"\n2. You must include these words in the text: " + vocabularyList +
		"\n3. You must use these grammar structures in the text" + grammarList +
		"\n4. You must also provide the English translation of the Japanese text." +
		"\n5. The text should be natural and based on everyday themes. It must be easy to understand and appropriately use the specified vocabulary and grammar points. It must be around 500 Japanese characters long. You must include line breaks in the Japanese text, and the English translation, between paragraphs, for readability." +
		"\n6. The reading comprehension quiz should include 5 different questions about the text that was generated, with 4 answer choices per question, in easy to understand Japanese, at the N4 JLPT level." +
		"\n7. The reading comprehension quiz questions should cover aspects of the whole text, each quiz questions should be about a different thing mentioned in the text." +
		"\n8. Return a valid JSON object based on provided schema."
}

func GetSongLyricsPrompt(vocabulary []string, grammar []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")
	grammarList := strings.Join(grammar, ", ")

	return "You are an AI that generates Japanese song lyrics from givem vocabulary and grammar. Please create a song with at least 5 stanzas, separated by lines, using the following conditions:" +
		"\n1. The language difficulty should be at around the N4 JLPT exam level." +
		"\n2. The lyrics MUST include these exact Japanese words in them: " + vocabularyList +
		"\n3. You must use these grammar structures" + grammarList +
		"\n4. The lyrics should be natural and based on everyday themes. The song should make sense. It must be easy to understand and MUST appropriately use the given vocabulary and grammar points."
}

func GetKanjiImagePrompt(kanji []model.Kanji) string {
	basePrompt := "Create visually interesting images to help memorize the meaning of kanji based on their components. Each image should help in remembering the meaning of the kanji by illustrating its components in a wildly creative and playful way. Do not mention the actual japanese radicals in the prompt, only use visual imagery. These will be used as prompts for an AI image generation model." +
		"\n Here is an example for the kanji candy 菓 (which is made up of the components grass, rice paddy and tree):\"A whimsical illustration of candy growing like fruit from a magical tree in a rice field, surrounded by tall grass. The candy is bright and colorful, hanging from green branches. The rice field is structured in traditional paddy style with water reflecting the sky.\"" +
		"\n Your descriptions can be in any style or mood. The goal is to create memorable images that will help learners remember the meaning of the kanji."
	var prompts []string

	for _, k := range kanji {
		prompt := fmt.Sprintf("Kanji: '%s' - Meaning: '%s' - Components: (%s)", k.Kanji, k.Meaning, k.Components)
		prompts = append(prompts, prompt)
	}

	return basePrompt + strings.Join(prompts, "\n")
}

func GetGrammarSentencesPrompt(grammarPoint string, vocabularyList []string) string {
	return "You are an AI that generates Japanese sentences. Please create 5 sentences in Japanese, with English translation, for the given grammar point, using the specified vocabulary. The sentences should be at the N4 JLPT level.\n" +
		"\n Each sentence must include the exact Japanese grammar point: " + grammarPoint +
		"\n You must incorporate this vocabulary into the sentences: " + strings.Join(vocabularyList, ", ") +
		"\n The sentences must be grammatically correct."
}

func GetVocabularySentencesPrompt(vocabulary []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")

	return "You are an AI that generates Japanese sentences. Please create 5 simple and easy-to-understand sentences in Japanese, with English translation, for each of the given Japanese vocabulary words. " +
		"\n You must use simple Japanese, at the beginner level, around N4 JLPT level, using basic grammar and vocabulary." +
		"\n Vocabulary list: " + vocabularyList + "\n" +
		"\n Ensure each sentence is unique and uses the vocabulary word in a different context." +
		"\n The sentences MUST be suitable for beginners and use straightforward language. The sentences must be grammatically correct." +
		"\n You MUST use the exact given Japanese word in each sentence. Do not change the word form or use synonyms."

}
