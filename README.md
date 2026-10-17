# ktype

ktype is a typing test for the terminal, built with Go and tcell.

(video placeholder)

## Features

- Real-time feedback
- Statistics summary (accuracy, WPM, raw WPM, and time)
- WPM graph
- Built-in and custom sources (compatible with Monkeytype)
- Quote support
- Stop on error modes
- Theme customization
- Simple design

## Installation

First, install [Go](https://go.dev/dl/) on your machine. Then run:

```sh
go install github.com/eduqr/ktype/cmd/ktype@latest
```

## Usage

Use `ktype --help` to view the available options:

```
Typing test for the terminal

Usage:
  ktype [flags]
  ktype [command]

Available Commands:
  list        List all ktype sources

Flags:
  -h, --help            Help for ktype
  -q, --quote           Use a quote instead of random words
  -s, --source string   Challenge source (default "english")
  -w, --words int       The number of words for the challenge (default 25)
```

### Examples

| Command                  | Description                                                         |
|--------------------------|---------------------------------------------------------------------|
| `ktype`                  | Runs with the default configuration                                 |
| `ktype -w 10`            | Generates a 10-word challenge using the default source              |
| `ktype -w 15 -s spanish` | Generates a 15-word challenge from the spanish source               |
| `ktype -q -s english`    | Generates a challenge with a random quote from the specified source |
| `ktype -w 30 -q`         | `-w 30` is ignored because `-q` takes priority                      |

### Controls

During a challenge:

- `Esc`: aborts the current challenge.
- `Ctrl+C`: exits ktype.

After a challenge:

- `r`: repeats the current challenge.
- `Tab`: starts a new challenge.
- `Esc`: exits ktype.

## Configuration

When ktype is run for the first time, a config file is created in the config directory. The exact location depends on your operating system:

| OS      | Config directory                       | Config file                                       |
|---------|----------------------------------------|---------------------------------------------------|
| Linux   | `~/.config/ktype/`                     | `~/.config/ktype/config.yaml`                     |
| MacOS   | `~/Library/Application Support/ktype/` | `~/Library/Application Support/ktype/config.yaml` |
| Windows | `%AppData%\ktype\`                     | `%AppData%\ktype\config.yaml`                     |

```yaml
challenge:
  
    # true: make the default challenge use quotes
    quote: false

    # file name (without extension) for the default source
    source: english

    # possible values: word, letter, or off
    stop_on: "word"

    # default number of words
    words: 25

# personalize and create your own theme!
theme:
  
    # accuracy color at the end of the challenge
    accuracy: '#e0af68'

    # actions color at the end of the challenge
    actions: '#565f89'

    # correct char color
    correct: '#b9f27c'

    # error messages color
    error: '#f7768e'

    # graph color at the end of the challenge
    graph: '#7aa2f7'

    # incorrect char color
    incorrect: '#f7768e'

    # missed char color
    missed: '#853F4C'

    # pending char color
    normal: '#a9b1d6'

    # raw WPM color at the end of the challenge
    raw_wpm: '#7dcfff'

    # time color at the end of the challenge
    time: '#bb9af7'

    # ktype title color
    title: '#cfc9c2'

    # WPM color at the end of the challenge
    wpm: '#7aa2f7'
```

### `stop_on`

- `off`: You can skip characters and advance to the next word even when overflow or errors are present.
- `letter`: If an error is made, the cursor will not advance to the next character until it is fixed.
- `word`: If an error is made, you can't advance to the next word until it is fixed.

## Built-in and custom sources

You can use `ktype list` to list all available word and quote sources.

```
Words:
  english        built-in
  english_10k    built-in
  english_25k    built-in
  spanish        built-in
  spanish_10k    built-in
  test           built-in
  code_c++       custom
  english        custom (overrides built-in)

Quotes:
  english        built-in
  spanish        built-in
  italian        custom
```

You can add custom sources by placing JSON files in the corresponding directory:

- Linux:
  - Words: `~/.config/ktype/words/`
  - Quotes: `~/.config/ktype/quotes/`
- macOS:
  - Words: `~/Library/Application Support/ktype/words/`
  - Quotes: `~/Library/Application Support/ktype/quotes/`
- Windows:
  - Words: `%AppData%\ktype\words\`
  - Quotes: `%AppData%\ktype\quotes\`

You can override built-in sources by using the same file name. Custom sources take priority over built-in sources.

### Finding custom sources

You can download a source from the Monkeytype repository and save it in one of the directories mentioned above.

- Words: https://github.com/monkeytypegame/monkeytype/tree/master/frontend/static/languages
- Quotes: https://github.com/monkeytypegame/monkeytype/tree/master/frontend/static/quotes

## Acknowledgments

- [tcell](https://github.com/gdamore/tcell): The main component of this project.
- [Monkeytype](https://monkeytype.com/): All built-in sources come from the Monkeytype repository. Thank you to the Monkeytype team and all its contributors for maintaining a great, curated list of sources.
- [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper): Packages for commands and configuration.
- [asciigraph](https://github.com/guptarohit/asciigraph): A package for rendering ASCII graphs.
- [Tokyo Night Theme](https://github.com/tokyo-night/tokyo-night-vscode-theme): The default theme uses Tokyo Night colors.