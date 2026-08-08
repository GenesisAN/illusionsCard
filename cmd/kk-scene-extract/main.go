package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	illusionCard "github.com/GenesisAN/illusionsCard"
	"github.com/GenesisAN/illusionsCard/KK"
)

func main() {
	input := flag.String("input", "", "path to a KK Studio scene PNG")
	output := flag.String("output", "", "directory for extracted character-card PNGs")
	mode := flag.String("mode", "unique", "extraction mode: unique or all")
	flag.Parse()

	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := extract(*input, *output, *mode); err != nil {
		fmt.Fprintln(os.Stderr, "extract:", err)
		os.Exit(1)
	}
}

func extract(input, output, mode string) error {
	card, err := illusionCard.ReadCardFromPath(input)
	if err != nil {
		return err
	}
	scene, ok := card.(*KK.KKSceneCard)
	if !ok {
		return fmt.Errorf("%s is not a KK Studio scene card", input)
	}

	var indexes []int
	switch mode {
	case "all":
		indexes = make([]int, len(scene.CharaCards))
		for index := range indexes {
			indexes[index] = index
		}
	case "unique":
		indexes = scene.UniqueCharaCardIndexes()
	default:
		return fmt.Errorf("unknown mode %q: expected unique or all", mode)
	}

	files := make([][]byte, 0, len(indexes))
	for _, index := range indexes {
		file, err := scene.ExtractCharaCard(index)
		if err != nil {
			return err
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return errors.New("scene contains no extractable character cards")
	}
	if err := os.MkdirAll(output, 0755); err != nil {
		return err
	}

	for index, data := range files {
		name := "character"
		sourceIndex := indexes[index]
		if sourceIndex < len(scene.CharaCards) && scene.CharaCards[sourceIndex].CharInfo != nil {
			info := scene.CharaCards[sourceIndex].CharInfo
			name = safeFilePart(strings.TrimSpace(info.Lastname + "_" + info.Firstname))
		}
		path := filepath.Join(output, fmt.Sprintf("%02d_%s.png", index+1, name))
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return err
		}
		if _, err := file.Write(data); err != nil {
			file.Close()
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
		fmt.Println(path)
	}
	return nil
}

func safeFilePart(value string) string {
	replacer := strings.NewReplacer(
		"<", "_", ">", "_", ":", "_", "\"", "_", "/", "_",
		"\\", "_", "|", "_", "?", "_", "*", "_",
	)
	value = strings.Trim(replacer.Replace(value), " .")
	if value == "" {
		return "character"
	}
	return value
}
