package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"jlpt-practice-helper/deprecatedemailtemplateflow"
	"jlpt-practice-helper/jsonfileflow/contentgenerator"
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

	err := jsonFileFlow()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Email has been successfully sent!")
	endTime := time.Now()
	fmt.Printf("Application ended at: %s\n", endTime.Format("2006-01-02 15:04:05"))

	duration := endTime.Sub(startTime)
	fmt.Printf("The application took %v to run.\n", duration)
}

func jsonFileFlow() error {
	err := database.LoadCSVData()
	if err != nil {
		return err
	}

	kanji := database.GetKanji(5)
	vocabulary := database.GetVocabulary(15, kanji)
	grammar := database.GetGrammar(3)

	kanjiArray := make([]string, len(kanji))
	for i, k := range kanji {
		kanjiArray[i] = k.Kanji
	}

	vocabularyArray := make([]string, len(vocabulary))
	for i, v := range vocabulary {
		vocabularyArray[i] = fmt.Sprintf("%s: %s", v.Kanji, v.Meaning)
	}

	grammarArray := make([]string, len(grammar))
	for i, g := range grammar {
		grammarArray[i] = fmt.Sprintf("%s: %s", g.Grammar, g.Meaning)
	}

	//get a short text with quiz based on the kanji, vocabulary, and grammar of today's lesson
	shortText, err := contentgenerator.GetShortTextWithQuiz(kanjiArray, vocabularyArray, grammarArray)
	if err != nil {
		log.Println(err)
	}

	//get song lyrics based on the vocabulary and grammar of today's lesson
	lyrics, err := contentgenerator.GenerateSongLyrics(vocabularyArray, grammarArray)
	if err != nil {
		log.Println(err)
	}

	//generate 2 songs from the lyrics
	songsFromLyrics, style, err := contentgenerator.GenerateSongFromLyrics(lyrics.LyricsJapanese, lyrics.SongTitle)
	var songIds []string
	for _, song := range songsFromLyrics {
		songIds = append(songIds, song.Id)
	}

	song := model.Song{
		Ids:        songIds,
		Style:      style,
		SongLyrics: *lyrics,
	}

	// get kanji mnemonic images based on the kanji of today's lesson
	images, err := contentgenerator.GenerateKanjiMnemonicImages(kanji)
	if err != nil {
		log.Println(err)
		return err
	}

	// generate image from mnemonic description
	for i, image := range images.KanjiList {
		encodedImage, err := contentgenerator.GenerateImageForKanji(image.ImageDescription)
		if err != nil {
			log.Println(err)
		}
		images.KanjiList[i].EncodedImage = encodedImage
	}

	// get example sentences for the grammar of today's lesson, using today's vocabulary
	var grammarInfo []model.GrammarInfo
	for _, pieceOfGrammar := range grammar {
		sentencesForGrammar, err := contentgenerator.GenerateExampleSentencesForGrammar(pieceOfGrammar.Grammar, vocabularyArray)
		if err != nil {
			log.Println(err)
			continue
		}
		grammarInfo = append(grammarInfo, model.GrammarInfo{
			Grammar:          pieceOfGrammar,
			GrammarSentences: *sentencesForGrammar,
		})
	}

	file := model.File{
		Kanji:             kanji,
		Vocabulary:        vocabulary,
		GrammarInfo:       grammarInfo,
		ShortTextWithQuiz: *shortText,
		Song:              song,
		KanjiImages:       *images,
	}

	// Call the new private method to write the file object to a JSON file
	if err := writeFileToJSON(file); err != nil {
		return err
	}

	return nil
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

func deprecatedEmailTemplateFlow() error {
	var AirTableApi = deprecatedemailtemplateflow.AirTableApi{}

	AirTableApi.InitializeClient()
	kanji := AirTableApi.GetKanji()
	vocabulary := AirTableApi.GetVocabulary()
	grammar := AirTableApi.GetGrammar()
	videoUrl := AirTableApi.GetListening()
	book := AirTableApi.GetBook()

	data := struct {
		Kanji      []deprecatedemailtemplateflow.Kanji
		Vocabulary []deprecatedemailtemplateflow.Vocabulary
		Grammar    []deprecatedemailtemplateflow.Grammar
		VideoUrl   string
		Book       deprecatedemailtemplateflow.Book
	}{
		Kanji:      kanji,
		Vocabulary: vocabulary,
		Grammar:    grammar,
		VideoUrl:   videoUrl,
		Book:       book,
	}

	return sendEmail(data, "deprecatedemailtemplateflow/mail_template_full_lesson.html")
}
