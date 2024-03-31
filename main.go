package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func main() {
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

	fmt.Println("mail sent")
}
