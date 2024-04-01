package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	startTime := time.Now()
	fmt.Printf("Application started at: %s\n", startTime.Format("2006-01-02 15:04:05"))

	godotenv.Load()

	var AirTableApi = AirTableApi{}

	AirTableApi.InitializeClient()
	kanji := AirTableApi.GetKanji()
	vocabulary := AirTableApi.GetVocabulary()
	grammar := AirTableApi.GetGrammar()
	videoUrl := AirTableApi.GetListening()
	book := AirTableApi.GetBook()

	err := sendEmail(kanji, vocabulary, grammar, videoUrl, book, "mail_template_fun.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email has been successfully sent!")
	endTime := time.Now()
	fmt.Printf("Application ended at: %s\n", endTime.Format("2006-01-02 15:04:05"))

	duration := endTime.Sub(startTime)
	fmt.Printf("The application took %v to run.\n", duration)
}
