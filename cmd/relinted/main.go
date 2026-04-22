package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"relinted/internal/formatter"
	"relinted/internal/io"
	"relinted/internal/perl"
)

var extToLang = map[string]string{
	".c":   "c",
	".h":   "c",
	".cpp": "c",
	".cc":  "c",
	".pl":  "perl",
	".pm":  "perl",
}

func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if lang, ok := extToLang[ext]; ok {
		return lang
	}
	return "c"
}

func main() {
	var langFlag string
	flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection)")
	flag.StringVar(&langFlag, "lang", "", "Language to use (overrides extension detection)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: relinted [-l|--lang lang] <input> [output]\n")
		os.Exit(1)
	}

	inputPath := args[0]

	content, err := io.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	lang := langFlag
	if lang == "" {
		lang = detectLanguage(inputPath)
	}

	var output string
	switch lang {
	case "c":
		output = formatter.Format(content)
	case "perl":
		output = perl.Format(content)
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported language %q\n", lang)
		fmt.Fprintf(os.Stderr, "Supported languages: c, perl\n")
		os.Exit(1)
	}

	if len(args) >= 2 {
		outputPath := args[1]
		if err := io.WriteFile(outputPath, output); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
			os.Exit(1)
		}
	} else {
		fmt.Print(output)
	}
}
