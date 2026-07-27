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

func ListSources() {
	ktypeDir, err := config.Dir()
	if err != nil {
		fmt.Println("getting config directory:", err)
		return
	}

	var result Sources

	wordSources, err := getEmbeddedEntries(wordsContent, "data/words", "built-in")
	if err != nil {
		fmt.Println("getting built-in word sources:", err)
		return
	}
	result.Words = wordSources

	customWordSources, err := getEntries(filepath.Join(ktypeDir, "words"), "custom")
	if err != nil {
		fmt.Println("getting custom word sources:", err)
		return
	}
	result.Words = append(result.Words, markOverrides(wordSources, customWordSources)...)

	quoteSources, err := getEmbeddedEntries(quotesContent, "data/quotes", "built-in")
	if err != nil {
		fmt.Println("getting built-in quote sources:", err)
		return
	}
	result.Quotes = quoteSources

	customQuoteSources, err := getEntries(filepath.Join(ktypeDir, "quotes"), "custom")
	if err != nil {
		fmt.Println("getting custom quote sources:", err)
		return
	}
	result.Quotes = append(result.Quotes, markOverrides(quoteSources, customQuoteSources)...)

	printSources(result)
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

func printSources(sources Sources) {
	width := maxNameLen(sources)

	fmt.Println("Words:")
	for _, source := range sources.Words {
		printSource(width, source)
	}

	fmt.Println("\nQuotes:")
	for _, source := range sources.Quotes {
		printSource(width, source)
	}
}

func printSource(width int, source Source) {
	origin := source.Origin
	if source.Overrides {
		origin += " (overrides built-in)"
	}

	fmt.Printf("  %-*s  %s\n", width, source.Name, origin)
}

func maxNameLen(sources Sources) int {
	max := 0

	for _, s := range sources.Words {
		if len(s.Name) > max {
			max = len(s.Name)
		}
	}

	for _, s := range sources.Quotes {
		if len(s.Name) > max {
			max = len(s.Name)
		}
	}

	return max + 2
}
