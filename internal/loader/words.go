package loader

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
)

type WordList struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

//go:embed data/words/*
var wordsContent embed.FS

func loadWords(language string, count int) ([]string, error) {
	filename := fmt.Sprintf("data/words/%s.json", language)

	byteValue, err := wordsContent.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("language not found: %s", language)
	}

	var list WordList
	if err := json.Unmarshal(byteValue, &list); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(list.Words) < count {
		count = len(list.Words)
	}

	rand.Shuffle(len(list.Words), func(i, j int) {
		list.Words[i], list.Words[j] = list.Words[j], list.Words[i]
	})

	return list.Words[:count], nil
}
