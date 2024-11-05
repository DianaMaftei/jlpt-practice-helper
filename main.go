package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"jlpt-practice-helper/jsonfileflow/contentgenerator"
	"jlpt-practice-helper/jsonfileflow/contentgenerator/gemini"
	"jlpt-practice-helper/jsonfileflow/database"
	"jlpt-practice-helper/jsonfileflow/model"
	"log"
	"os"
	"time"
)

func main() {
	startTime := time.Now()
	fmt.Printf("Application started at: %s\n", startTime.Format("2006-01-02 15:04:05"))

	godotenv.Load()

	if err := jsonFileFlow(); err != nil {
		log.Fatal(err)
	}

	endTime := time.Now()
	fmt.Printf("Application ended at: %s\n", endTime.Format("2006-01-02 15:04:05"))

	duration := endTime.Sub(startTime)
	fmt.Printf("The application took %v to run.\n", duration)
}

func jsonFileFlow() error {
	if err := database.LoadCSVData(); err != nil {
		return err
	}

	// fetch data from DB
	kanji := database.GetKanji(5)
	vocabulary := database.GetVocabulary(21, kanji)
	grammar := database.GetGrammar(2)

	vocabularyArray := extractVocabularyArray(vocabulary)
	grammarArray := extractGrammarArray(grammar)

	// generate content
	shortText, err := gemini.GetShortTextWithQuiz(vocabularyArray, grammarArray)
	if err != nil {
		logError(err)
	}

	lyrics, err := gemini.GenerateSongLyrics(vocabularyArray, grammarArray)
	if err != nil {
		logError(err)
	}

	songsFromLyrics, style, err := contentgenerator.GenerateSongFromLyrics(lyrics.LyricsJapanese, lyrics.SongTitle)
	if err != nil {
		logError(err)
	}

	songIds := extractSongIds(songsFromLyrics)
	song := model.Song{
		Ids:        songIds,
		Style:      style,
		SongLyrics: *lyrics,
	}

	images, err := gemini.GenerateKanjiMnemonicImages(kanji)
	if err != nil {
		logError(err)
	}

	kanjiInfoList := generateKanjiInfoList(kanji, images)

	grammarInfoList, err := generateGrammarInfoList(grammar, vocabularyArray)
	if err != nil {
		logError(err)
	}

	vocabSentences, err := gemini.GenerateSentencesForVocabulary(vocabularyArray)
	if err != nil {
		logError(err)
	}

	vocabularyInfoList := generateVocabularyInfoList(vocabulary, vocabSentences)

	// write to file
	file := model.File{
		KanjiInfo:         kanjiInfoList,
		VocabularyInfo:    vocabularyInfoList,
		GrammarInfo:       grammarInfoList,
		ShortTextWithQuiz: *shortText,
		Song:              song,
	}

	err = writeFileToJSON(file)
	if err != nil {
		logError(err)
		return err
	}

	// update database
	err = updateDatabase(kanji, vocabulary, grammar)
	if err != nil {
		logError(err)
		return err
	}

	// send email
	today := time.Now().Format("2006-01-02")
	err = sendEmail(struct {
		CurrentDate string
		PageURL     string
	}{
		CurrentDate: today,
		PageURL:     fmt.Sprintf("https://dianamaftei.github.io/japanese-study-pages/page/%s", today),
	}, "static/email_template.html")
	if err != nil {
		logError(err)
		return err
	}
	return nil
}

func extractKanjiArray(kanji []model.Kanji) []string {
	kanjiArray := make([]string, len(kanji))
	for i, k := range kanji {
		kanjiArray[i] = k.Kanji
	}
	return kanjiArray
}

func extractVocabularyArray(vocabulary []model.Vocabulary) []string {
	vocabularyArray := make([]string, len(vocabulary))
	for i, v := range vocabulary {
		vocabularyArray[i] = fmt.Sprintf("%s: %s", v.Kanji, v.Meaning)
	}
	return vocabularyArray
}

func extractGrammarArray(grammar []model.Grammar) []string {
	grammarArray := make([]string, len(grammar))
	for i, g := range grammar {
		grammarArray[i] = fmt.Sprintf("%s: %s", g.Grammar, g.Meaning)
	}
	return grammarArray
}

func extractSongIds(songs []model.SongResponse) []string {
	var songIds []string
	for _, song := range songs {
		songIds = append(songIds, song.Id)
	}
	return songIds
}

func generateKanjiInfoList(kanji []model.Kanji, images *model.KanjiImageResponse) []model.KanjiInfo {
	var kanjiInfoList []model.KanjiInfo
	for _, image := range images.KanjiList {
		encodedImage, err := contentgenerator.GenerateImageFromPrompt(image.ImageDescription)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, k := range kanji {
			if k.Kanji == image.Kanji {
				kanjiInfoList = append(kanjiInfoList, model.KanjiInfo{
					Kanji:            k,
					ImageDescription: image.ImageDescription,
					EncodedImage:     encodedImage,
				})
				break
			}
		}
	}
	return kanjiInfoList
}

func generateGrammarInfoList(grammar []model.Grammar, vocabularyArray []string) ([]model.GrammarInfo, error) {
	var grammarInfoList []model.GrammarInfo
	for _, pieceOfGrammar := range grammar {
		sentencesForGrammar, err := gemini.GenerateExampleSentencesForGrammar(pieceOfGrammar.Grammar, vocabularyArray)
		if err != nil {
			return nil, err
		}
		grammarInfoList = append(grammarInfoList, model.GrammarInfo{
			Grammar:          pieceOfGrammar,
			GrammarSentences: *sentencesForGrammar,
		})
	}
	return grammarInfoList, nil
}

func generateVocabularyInfoList(vocabulary []model.Vocabulary, vocabSentences *model.VocabularySentencesResponse) []model.VocabularyInfo {
	var vocabularyInfoList []model.VocabularyInfo
	for _, vocab := range vocabulary {
		var sentences []model.Sentence
		for _, vocabSentence := range vocabSentences.VocabularySentences {
			if vocabSentence.Vocabulary == vocab.Kanji {
				sentences = vocabSentence.Sentences
				break
			}
		}
		vocabularyInfoList = append(vocabularyInfoList, model.VocabularyInfo{
			Vocabulary: vocab,
			Sentences:  sentences,
		})
	}
	return vocabularyInfoList
}

func writeFileToJSON(file model.File) error {
	today := time.Now().Format("2006-01-02")
	jsonFileName := fmt.Sprintf("out/%s.json", today)

	log.Printf("Writing file to JSON: %s\n", jsonFileName)

	jsonData, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return err
	}

	err = os.WriteFile(jsonFileName, jsonData, 0644)
	if err != nil {
		log.Println("Error writing JSON to file:", err)
		return err
	}

	log.Println("Successfully wrote JSON to file:", jsonFileName)
	return nil
}

func updateDatabase(kanji []model.Kanji, vocabulary []model.Vocabulary, grammar []model.Grammar) error {
	log.Println("Updating database...")

	today := time.Now().Format("2006-01-02")

	for i := range kanji {
		kanji[i].Seen = true
		kanji[i].UpdateDate = today
	}
	for i := range vocabulary {
		vocabulary[i].Seen = true
		vocabulary[i].UpdateDate = today
	}
	for i := range grammar {
		grammar[i].Seen = true
		grammar[i].UpdateDate = today
	}

	if err := database.SaveKanji(kanji); err != nil {
		return err
	}
	if err := database.SaveVocabulary(vocabulary); err != nil {
		return err
	}
	if err := database.SaveGrammar(grammar); err != nil {
		return err
	}
	log.Println("Successfully updated the database")
	return nil
}

func logError(err error) {
	log.Println(err)
}
