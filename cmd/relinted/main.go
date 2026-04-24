package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"relinted/internal/formatter"
	"relinted/internal/io"
	"relinted/internal/js"
	"relinted/internal/perl"
	"relinted/internal/rust"
)

var extToLang = map[string]string{
	".c":   "c",
	".h":   "c",
	".cpp": "c",
	".cc":  "c",
	".pl":  "perl",
	".pm":  "perl",
	".js":  "js",
	".rs":  "rust",
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
	flag.StringVar(&langFlag, "l", "", "Language to use (overrides extension detection) [c, perl, rust, js]")
	flag.Usage = func() {
		fmt.Println("Relinted reformats C/C++, Perl, Rust, and JavaScript source code to visually resemble Python.")
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: relinted [-l lang] <input> [output]\n")
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
	case "js":
		output = js.Format(content)
	case "rust":
		output = rust.Format(content)
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported language %q\n", lang)
		fmt.Fprintf(os.Stderr, "Supported languages: c, perl, rust, js\n")
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
