package dbs

import (
	"encoding/csv"
	"jlpt-practice-helper/jsonfileflow/model"
	"os"
	"strconv"
)

var vocabularyFilePath = "static/data/vocabulary.csv"

func LoadVocabularyFromFile() ([]model.Vocabulary, error) {
	var vocabularies []model.Vocabulary

	f, err := os.Open(vocabularyFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records[1:] { // Skip header
		var seen bool
		if record[3] == "" {
			seen = false
		} else {
			seen, err = strconv.ParseBool(record[3])
			if err != nil {
				return nil, err
			}
		}
		vocabularies = append(vocabularies, model.Vocabulary{
			Kanji:   record[0],
			Reading: record[1],
			Meaning: record[2],
			Metadata: model.Metadata{
				Seen:       seen,
				UpdateDate: record[4],
			},
		})
	}

	return vocabularies, nil
}

func SaveVocabularyToFile(vocabularies []model.Vocabulary) error {
	file, err := os.Create(vocabularyFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Kanji", "Reading", "Meaning", "Seen", "Update_Date"})

	for _, v := range vocabularies {
		err := writer.Write([]string{v.Kanji, v.Reading, v.Meaning, strconv.FormatBool(v.Metadata.Seen), v.Metadata.UpdateDate})
		if err != nil {
			return err
		}
	}

	return nil
}
