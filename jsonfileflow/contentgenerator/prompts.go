package contentgenerator

import (
	"fmt"
	"jlpt-practice-helper/jsonfileflow/model"
	"strings"
)

func GetShortTextPrompt(vocabulary []string, grammar []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")
	grammarList := strings.Join(grammar, ", ")

	return "You are an AI Japanese language tutor specializing in JLPT preparation. Generate a reading text in Japanese, with comprehension questions and English Translation, following these strict requirements:" +
		"\n\nTEXT REQUIREMENTS:" +
		"\n- Length: Strictly 450-500 Japanese characters" +
		"\n- Level: N4 JLPT base level (beginner level)" +
		"\n- The ONLY advanced elements allowed are:" +
		"\n  1. The specified vocabulary words from the input list" +
		"\n  2. The specified grammar points from the input list" +
		"\n- All other vocabulary and grammar must be only up to N4 level" +
		"\n- The English translation must accurately match the Japanese text" +
		"\n\nVOCABULARY AND GRAMMAR CONTROL:" +
		"\n- You MUST use the exact Japanese words provided (no synonyms or variations)" +
		"\n- You MUST use all specified grammar points naturally in the text" +
		"\n\nTEXT STRUCTURE:" +
		"\n- Genre: Can be story, blog, letter, advertisement, etc." +
		"\n- Theme: Must be everyday situations that N4 learners can relate to" +
		"\n- Format: Include clear paragraph breaks for readability" +
		"\n- Context: Should be practical and interesting for language learners" +
		"\n\nCOMPREHENSION QUIZ REQUIREMENTS:" +
		"\n- Exactly 5 questions" +
		"\n- Each question must:" +
		"\n  * Focus on a different aspect/paragraph of the text" +
		"\n  * Have exactly 4 answer options" +
		"\n  * Have only one correct answer" +
		"\n  * Randomize the position of the correct answer" +
		"\n- All questions and answers must be use easy Japanese, at the beginner level" +
		"\n\nVERIFICATION STEPS:" +
		"\n1. Text Length: Count characters to ensure 450-500 limit" +
		"\n2. Vocabulary Check: Verify all non-input words are N4 level" +
		"\n3. Grammar Check: Verify all non-input grammar is N4 level" +
		"\n4. Question Distribution:" +
		"\n   - Confirm each question covers different content" +
		"\n   - Verify correct answers are randomly distributed" +
		"\n   - Ensure all questions are at N4 (easy) level" +
		"\n5. Comprehensiveness: Verify text includes all required vocabulary and grammar" +
		"\n\nINPUT:" +
		"\nVocabulary List: " + vocabularyList +
		"\nGrammar List: " + grammarList
}

func GetSongLyricsPrompt(vocabulary []string, grammar []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")
	grammarList := strings.Join(grammar, ", ")

	return "You are an AI that generates Japanese song lyrics from given vocabulary and grammar. Please create a song with at least 5 stanzas, separated by lines, using the following conditions:" +
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
	return "You are an AI Japanese language tutor specializing in JLPT preparation. Your task is to generate example sentences in Japanese (and their English translations), demonstrating specific grammar points, following these strict requirements:" +
		"\nBASE LEVEL REQUIREMENTS:" +
		"\n* All supporting vocabulary and structures must be JLPT N4 level or below" +
		"\n* The ONLY advanced elements allowed are:" +
		"\n    1. The specified grammar point being practiced" +
		"\n    2. The vocabulary words provided in the input list" +
		"\n* Do not use any other vocabulary or grammar from N3, N2, or N1 levels" +
		"\n* Do not use words that aren't in official JLPT vocabulary lists" +
		"\nSENTENCE CONSTRUCTION:" +
		"\n* Create 5 unique sentences that naturally demonstrate the target grammar point" +
		"\n* Each sentence must use the grammar point exactly as specified, with the given meaning" +
		"\n* Each sentence must incorporate at least one word from the provided vocabulary list" +
		"\n* Use clear, practical contexts that a N4-level learner would encounter" +
		"\n* Provide an accurate English translation for each sentence" +
		"\nINPUT GRAMMAR POINT: " + grammarPoint +
		"\nINPUT VOCABULARY LIST: " + strings.Join(vocabularyList, ", ")
}

func GetVocabularySentencesPrompt(vocabulary []string) string {
	vocabularyList := strings.Join(vocabulary, ", ")

	return "You are an AI Japanese language tutor specializing in JLPT preparation. Your task is to generate Japanese example sentences, and their English translation, following these strict requirements:" +
		"\nBASE LEVEL REQUIREMENTS:" +
		"\n* All sentences must use N5-N4 JLPT level vocabulary and grammar as their foundation" +
		"\n* The ONLY exceptions allowed are the N3 vocabulary words provided in the input list" +
		"\n* Do not use any vocabulary or grammar from N2 or N1 levels, or above" +
		"\n* Do not use any words that aren't included in official JLPT vocabulary lists" +
		"\nSENTENCE STRUCTURE:" +
		"\n* Create 5 unique sentences for each provided N3 vocabulary word" +
		"\n* Use simple, clear sentence structures common in everyday Japanese" +
		"\n* Focus on practical, real-life situations" +
		"\nCRITICAL WORD USAGE RULE:" +
		"\n* In each sentence you MUST use the exact Japanese word provided, and it must be used with the meaning that is given" +
		"\n* Do NOT substitute the given Japanese word with any synonyms or alternative words with the same meaning" +
		"\n* Do NOT change the word form unless specifically instructed" +
		"\n* The English meaning is provided only for disambiguation" +
		"\n* Example: For input {\"続く\": \"to continue\"}, you must use 続く, not 継続する or any other synonym" +
		"\nINPUT VOCABULARY LIST: " + vocabularyList +
		"\nVERIFICATION STEPS: Before outputting each sentence:" +
		"\n1. Verify every word against N5-N4 vocabulary list (except for provided N3 terms)" +
		"\n2. Confirm grammar structures are N5-N4 level" +
		"\n3. Check that sentence length and complexity are appropriate for N4 learners" +
		"\n4. Ensure the context is clear and practical"

}
