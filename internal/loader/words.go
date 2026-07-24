package loader

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/eduqr/ktype/internal/config"
)

type WordList struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

//go:embed data/words/*
var wordsContent embed.FS

func loadWords(source string, count int) ([]string, error) {
	byteValue, err := readWordSource(source)
	if err != nil {
		return nil, err
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

func readWordSource(source string) ([]byte, error) {
	ktypeDir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	customPath := filepath.Join(ktypeDir, "words", source+".json")
	byteValue, err := os.ReadFile(customPath)
	if err == nil {
		return byteValue, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read custom words source %q: %w", source, err)
	}

	filename := fmt.Sprintf("data/words/%s.json", source)
	byteValue, err = wordsContent.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("words source not found: %s", source)
	}

	return byteValue, nil
}
