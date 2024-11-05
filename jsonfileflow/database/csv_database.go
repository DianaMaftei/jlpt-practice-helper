package database

import (
	"jlpt-practice-helper/jsonfileflow/database/dbs"
	"jlpt-practice-helper/jsonfileflow/model"
	"log"
	"math/rand"
	"strings"
)

var grammar []model.Grammar
var kanji []model.Kanji
var vocabulary []model.Vocabulary

func LoadCSVData() error {
	var err error

	log.Println("Loading grammar data...")
	grammar, err = dbs.LoadGrammarFromFile()
	if err != nil {
		log.Printf("Error loading grammar data: %v", err)
		return err
	}
	log.Println("Grammar data loaded successfully.")

	log.Println("Loading kanji data...")
	kanji, err = dbs.LoadKanjiFromFile()
	if err != nil {
		log.Printf("Error loading kanji data: %v", err)
		return err
	}
	log.Println("Kanji data loaded successfully.")

	log.Println("Loading vocabulary data...")
	vocabulary, err = dbs.LoadVocabularyFromFile()
	if err != nil {
		log.Printf("Error loading vocabulary data: %v", err)
		return err
	}
	log.Println("Vocabulary data loaded successfully.")

	return nil
}

func GetGrammar(count int) []model.Grammar {
	// Shuffle the grammar slice to ensure randomness
	rand.Shuffle(len(grammar), func(i, j int) {
		grammar[i], grammar[j] = grammar[j], grammar[i]
	})

	var selectedGrammar []model.Grammar

	// Select grammar items that haven't been seen, up to the specified count
	for _, g := range grammar {
		if !g.Seen && len(selectedGrammar) < count { // Check if 'seen' is false and we haven't selected enough
			selectedGrammar = append(selectedGrammar, g)
		}
		if len(selectedGrammar) == count { // Stop if we've selected the required count
			break
		}
	}

	return selectedGrammar
}

func GetKanji(count int) []model.Kanji {
	var selectedKanji []model.Kanji
	for _, k := range kanji {
		if !k.Seen && len(selectedKanji) < count { // Check if 'seen' is false and we haven't selected enough
			selectedKanji = append(selectedKanji, k)
		}
		if len(selectedKanji) == count { // Stop if we've selected the required count
			break
		}
	}

	return selectedKanji
}

func GetVocabulary(count int, kanjiArray []model.Kanji) []model.Vocabulary {
	// Shuffle the vocabulary slice to ensure randomness
	rand.Shuffle(len(vocabulary), func(i, j int) {
		vocabulary[i], vocabulary[j] = vocabulary[j], vocabulary[i]
	})

	var selectedVocabulary []model.Vocabulary
	var spareVocabulary []model.Vocabulary // Array to hold unmatched vocabulary

	// Search for vocabulary items containing any of the Kanji from the array
	for _, vocab := range vocabulary {
		if vocab.Seen {
			continue // Skip this item if it has been seen
		}
		for _, k := range kanjiArray {
			if strings.Contains(vocab.Kanji, k.Kanji) {
				selectedVocabulary = append(selectedVocabulary, vocab)
				if len(selectedVocabulary) == count {
					return selectedVocabulary // Return immediately if count is reached
				}
				break // Break to avoid adding the same vocab multiple times
			} else if len(spareVocabulary) < count {
				spareVocabulary = append(spareVocabulary, vocab) // Add to spare if not matched
			}
		}
		if len(selectedVocabulary) == count {
			return selectedVocabulary // Return immediately if count is reached
		}
	}

	// If not enough vocabulary items were found, fill from spareVocabulary
	for _, vocab := range spareVocabulary {
		if len(selectedVocabulary) >= count { // Stop if we've fulfilled the count
			break
		}
		selectedVocabulary = append(selectedVocabulary, vocab)
	}

	return selectedVocabulary // Return the selected vocabulary items
}

func SaveGrammar(updatedGrammar []model.Grammar) error {
	// Update the existing grammar in memory
	for _, updatedItem := range updatedGrammar {
		for i, item := range grammar {
			if item.Grammar == updatedItem.Grammar {
				grammar[i] = updatedItem
				break
			}
		}
	}

	// Write the updated grammar to file
	if err := dbs.SaveGrammarToFile(grammar); err != nil {
		log.Printf("Error saving grammar data: %v", err)
		return err
	}
	log.Println("Grammar data saved successfully.")
	return nil
}

func SaveKanji(updatedKanji []model.Kanji) error {
	// Update the existing kanji in memory
	for _, updatedItem := range updatedKanji {
		for i, item := range kanji {
			if item.Kanji == updatedItem.Kanji {
				kanji[i] = updatedItem
				break
			}
		}
	}

	// Write the updated kanji to file
	if err := dbs.SaveKanjiToFile(kanji); err != nil {
		log.Printf("Error saving kanji data: %v", err)
		return err
	}
	log.Println("Kanji data saved successfully.")
	return nil
}

func SaveVocabulary(updatedVocabulary []model.Vocabulary) error {
	// Update the existing vocabulary in memory
	for _, updatedItem := range updatedVocabulary {
		for i, item := range vocabulary {
			if item.Kanji == updatedItem.Kanji {
				vocabulary[i] = updatedItem
				break
			}
		}
	}

	// Write the updated vocabulary to file
	if err := dbs.SaveVocabularyToFile(vocabulary); err != nil {
		log.Printf("Error saving vocabulary data: %v", err)
		return err
	}
	log.Println("Vocabulary data saved successfully.")
	return nil
}
