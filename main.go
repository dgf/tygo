// Tygo is a terminal-based typing test.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/dgf/tygo/internal/config"
	"github.com/dgf/tygo/internal/dict"
	"github.com/dgf/tygo/internal/display"
	"github.com/dgf/tygo/internal/game"
	"github.com/dgf/tygo/internal/input"
	"golang.org/x/term"
)

// Exit codes.
const (
	ExitSuccess          = 0
	ExitUserError        = 1
	ExitEnvironmentError = 2
	ExitInternalError    = 3
)

func Dictionary(name string) dict.Dictionary {
	switch name {
	case "german":
		return dict.German10K
	case "english":
		return dict.English10K
	default:
		return dict.English10K
	}
}

func MustLoadWords(cfg config.Config, file string) []string {
	if len(file) == 0 {
		return dict.LoadDict(Dictionary(cfg.Dictionary), cfg.TopWords)
	}

	words, err := dict.LoadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dictionary load failed: %v\n", err)
		os.Exit(ExitUserError)
	}

	return words
}

func MustMakeRaw(in *os.File) (int, *term.State) {
	fd := int(in.Fd())

	if !term.IsTerminal(fd) {
		fmt.Fprintln(os.Stderr, "Use a terminal (requires a TTY)")
		os.Exit(ExitUserError)
	}

	termState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Raw mode activation failed: %v\n", err)
		os.Exit(ExitEnvironmentError)
	}

	return fd, termState
}

func main() {
	cfg, loadErr := config.LoadUserConfig()
	if loadErr != nil {
		cfg = config.Default()
		_ = config.WriteUserConfig(cfg)
	}

	var file string

	flag.StringVar(&cfg.Dictionary, "dict", cfg.Dictionary, "dictionary to use, available: german, english")

	flag.IntVar(&cfg.TopWords, "top", cfg.TopWords, "top count of words to load from source (dict or file)")
	flag.IntVar(&cfg.WordCount, "count", cfg.WordCount, "number of words to include in the typing test")
	flag.IntVar(&cfg.Width, "width", cfg.Width, "display width for the typing text")

	flag.BoolVar(&cfg.Numbers, "nums", cfg.Numbers, "enable number mode")
	flag.BoolVar(&cfg.Punctuation, "punct", cfg.Punctuation, "enable punctuation marks")
	flag.BoolVar(&cfg.StrictMode, "strict", cfg.StrictMode, "enable strict mode, restarts on every error")

	flag.StringVar(&file, "file", "", "vocabulary JSON file with 'words' list")

	flag.Parse()

	in := os.Stdin
	out := os.Stdout
	words := MustLoadWords(cfg, file)
	fd, oldState := MustMakeRaw(in)

	defer func() {
		_ = term.Restore(fd, oldState)

		if r := recover(); r != nil {
			fmt.Fprintf(out, "%v - %s", r, debug.Stack())
			os.Exit(ExitInternalError)
		}
	}()

	input.Loop(in, game.NewGame(cfg, words, display.NewRenderer(out)))
}
