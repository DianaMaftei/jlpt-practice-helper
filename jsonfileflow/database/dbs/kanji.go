package dbs

import (
	"encoding/csv"
	"jlpt-practice-helper/jsonfileflow/model"
	"os"
	"strconv"
)

var kanjiFilePath = "static/data/kanji.csv"

func LoadKanjiFromFile() ([]model.Kanji, error) {
	var kanjis []model.Kanji

	f, err := os.Open(kanjiFilePath)
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
		if record[8] == "" {
			seen = false
		} else {
			seen, err = strconv.ParseBool(record[8])
			if err != nil {
				return nil, err
			}
		}

		kanjis = append(kanjis, model.Kanji{
			Kanji:      record[0],
			OnReading:  record[1],
			KunReading: record[2],
			Meaning:    record[3],
			Examples:   record[4],
			Components: record[5],
			Koohii1:    record[6],
			Koohii2:    record[7],
			Metadata: model.Metadata{
				Seen:       seen,
				UpdateDate: record[9],
			},
		})
	}

	return kanjis, nil
}

func SaveKanjiToFile(kanjis []model.Kanji) error {
	file, err := os.Create(kanjiFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Kanji", "On_Reading", "Kun_Reading", "Meaning", "Examples", "Components", "Koohii_Story_1", "Koohii_Story_2", "Seen", "Update_Date"})

	for _, k := range kanjis {
		err := writer.Write([]string{k.Kanji, k.OnReading, k.KunReading, k.Meaning, k.Examples, k.Components, k.Koohii1, k.Koohii2, strconv.FormatBool(k.Metadata.Seen), k.Metadata.UpdateDate})
		if err != nil {
			return err
		}
	}

	return nil
}
