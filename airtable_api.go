package main

import (
	"fmt"
	"github.com/mehanizm/airtable"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type AirTableApi struct {
	Client *airtable.Client
}

func (a *AirTableApi) InitializeClient() {
	apiKey := os.Getenv("AIRTABLE_API_KEY")

	client := airtable.NewClient(apiKey)
	a.Client = client
}

func (a *AirTableApi) GetKanji() []Kanji {
	var kanji []Kanji

	table := a.Client.GetTable("appLROzl1t7bORKvn", "Kanji")
	records, err := getSortedUnseenRecords(table, "record_id", 5, "kanji", "meaning", "kun_reading", "on_reading")
	if err != nil {
		fmt.Println(err)
		return kanji
	}

	for _, record := range records.Records {
		character := record.Fields["kanji"].(string)
		escapedKanji := url.QueryEscape(character)
		baseURL := "https://raw.githubusercontent.com/jcsirot/kanji.gif/master/kanji/gif/150x150/"
		gifURL := baseURL + escapedKanji + ".gif"

		kunReading := getStringFieldOrDefaultEmpty(record, "kun_reading")
		onReading := getStringFieldOrDefaultEmpty(record, "on_reading")

		kanji = append(kanji, Kanji{
			Kanji:      character,
			Meaning:    record.Fields["meaning"].(string),
			KunReading: commaSeparatedList(kunReading),
			OnReading:  commaSeparatedList(onReading),
			GifUrl:     gifURL,
		})
	}

	markRecordsAsSeen(table, records)
	return kanji
}

func (a *AirTableApi) GetVocabulary() []Vocabulary {
	var vocabulary []Vocabulary

	table := a.Client.GetTable("appLROzl1t7bORKvn", "Vocabulary")
	records, err := getSortedUnseenRecords(table, "record_id", 30, "kanji", "kana", "english", "ex1_ja_furigana", "ex1_en")
	if err != nil {
		fmt.Println(err)
		return vocabulary
	}

	for _, record := range records.Records {
		exampleJp := getStringFieldOrDefaultEmpty(record, "ex1_ja_furigana")
		exampleEn := getStringFieldOrDefaultEmpty(record, "ex1_en")
		vocabulary = append(vocabulary, Vocabulary{
			Kanji:     record.Fields["kanji"].(string),
			Kana:      record.Fields["kana"].(string),
			Meaning:   record.Fields["english"].(string),
			ExampleJp: generateHtmlWithFurigana(exampleJp),
			ExampleEn: exampleEn,
		})
	}

	markRecordsAsSeen(table, records)
	return vocabulary
}

func (a *AirTableApi) GetGrammar() []Grammar {
	var grammar []Grammar

	table := a.Client.GetTable("appLROzl1t7bORKvn", "Grammar")
	records, err := getSortedUnseenRecords(table, "record_id", 2, "grammar", "meaning", "ex1_ja_furigana", "ex1_en", "ex2_ja_furigana", "ex2_en", "bunpro")
	if err != nil {
		fmt.Println(err)
		return grammar
	}

	for _, record := range records.Records {
		exampleJp1 := getStringFieldOrDefaultEmpty(record, "ex1_ja_furigana")
		exampleEn1 := getStringFieldOrDefaultEmpty(record, "ex1_en")
		exampleJp2 := getStringFieldOrDefaultEmpty(record, "ex2_ja_furigana")
		exampleEn2 := getStringFieldOrDefaultEmpty(record, "ex2_en")
		bunpro := getStringFieldOrDefaultEmpty(record, "bunpro")

		grammar = append(grammar, Grammar{
			Grammar:     record.Fields["grammar"].(string),
			Explanation: record.Fields["meaning"].(string),
			Example1Jp:  generateHtmlWithFurigana(exampleJp1),
			Example1En:  exampleEn1,
			Example2Jp:  generateHtmlWithFurigana(exampleJp2),
			Example2En:  exampleEn2,
			Bunpro:      bunpro,
		})
	}

	markRecordsAsSeen(table, records)
	return grammar
}

func (a *AirTableApi) GetListening() string {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Listening")
	records, err := getSortedUnseenRecords(table, "title", 1, "url")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var record = records.Records[0]
	_, err = record.UpdateRecordPartial(map[string]any{"seen": true})
	if err != nil {
		log.Fatal(err)
	}

	return record.Fields["url"].(string)
}

func (a *AirTableApi) GetBook() Book {
	table := a.Client.GetTable("appLROzl1t7bORKvn", "Books")

	records, err := getSortedUnseenRecords(table, "title", 1, "url", "img")
	if err != nil {
		fmt.Println(err)
		return Book{}
	}

	markRecordsAsSeen(table, records)

	return Book{
		Url: records.Records[0].Fields["url"].(string),
		Img: records.Records[0].Fields["img"].(string),
	}
}

func getSortedUnseenRecords(table *airtable.Table, sortField string, maxRecords int, fields ...string) (*airtable.Records, error) {
	sortQuery := struct {
		FieldName string
		Direction string
	}{sortField, "desc"}

	records, err := table.GetRecords().
		FromView("Grid view").
		WithFilterFormula("AND({seen}=0)").
		WithSort(sortQuery).
		ReturnFields(fields...).
		MaxRecords(maxRecords).
		Do()

	return records, err
}

func markRecordsAsSeen(table *airtable.Table, records *airtable.Records) {
	const maxRecordsPerRequest = 10

	for i := 0; i < len(records.Records); i += maxRecordsPerRequest {
		end := i + maxRecordsPerRequest
		if end > len(records.Records) {
			end = len(records.Records)
		}

		r := &airtable.Records{}
		r.Records = []*airtable.Record{}

		for _, record := range records.Records[i:end] {
			newRecord := airtable.Record{
				ID: record.ID,
				Fields: map[string]interface{}{
					"seen": true,
				},
			}

			r.Records = append(r.Records, &newRecord)
		}

		_, err := table.UpdateRecordsPartial(r)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func removeHTMLTags(s string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, " ")
}

func commaSeparatedList(s string) string {
	text := removeHTMLTags(s)
	items := strings.Fields(text)
	return strings.Join(items, ", ")
}

func getStringFieldOrDefaultEmpty(record *airtable.Record, fieldName string) string {
	field, ok := record.Fields[fieldName].(string)
	if !ok {
		field = ""
	}

	return field
}
