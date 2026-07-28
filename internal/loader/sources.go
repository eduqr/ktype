package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/eduqr/ktype/internal/config"
)

type Source struct {
	Name      string
	Origin    string
	Overrides bool
}

type Sources struct {
	Words  []Source
	Quotes []Source
}

func ListSources() (Sources, error) {
	ktypeDir, err := config.Dir()
	if err != nil {
		return Sources{}, fmt.Errorf("get config directory: %w", err)
	}

	var result Sources

	wordSources, err := getEmbeddedEntries(wordsContent, "data/words", "built-in")
	if err != nil {
		return Sources{}, fmt.Errorf("get built-in word sources: %w", err)
	}
	result.Words = wordSources

	customWordSources, err := getEntries(filepath.Join(ktypeDir, "words"), "custom")
	if err != nil {
		return Sources{}, fmt.Errorf("get custom word sources: %w", err)
	}
	result.Words = append(result.Words, markOverrides(wordSources, customWordSources)...)

	quoteSources, err := getEmbeddedEntries(quotesContent, "data/quotes", "built-in")
	if err != nil {
		return Sources{}, fmt.Errorf("get built-in quote sources: %w", err)
	}
	result.Quotes = quoteSources

	customQuoteSources, err := getEntries(filepath.Join(ktypeDir, "quotes"), "custom")
	if err != nil {
		return Sources{}, fmt.Errorf("get custom quote sources: %w", err)
	}
	result.Quotes = append(result.Quotes, markOverrides(quoteSources, customQuoteSources)...)

	return result, nil
}

func getEntries(path, origin string) ([]Source, error) {
	var sources []Source

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		basename := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		sources = append(sources, Source{Name: basename, Origin: origin})
	}

	return sources, nil
}

func getEmbeddedEntries(content fs.FS, path, origin string) ([]Source, error) {
	var sources []Source

	entries, err := fs.ReadDir(content, path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		basename := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		sources = append(sources, Source{Name: basename, Origin: origin})
	}

	return sources, nil
}

func markOverrides(builtIn, custom []Source) []Source {
	builtInNames := make(map[string]bool, len(builtIn))
	for _, source := range builtIn {
		builtInNames[source.Name] = true
	}

	for i := range custom {
		custom[i].Overrides = builtInNames[custom[i].Name]
	}

	return custom
}
