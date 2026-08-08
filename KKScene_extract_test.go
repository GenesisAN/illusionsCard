package illusionCard

import (
	"bytes"
	"os"
	"testing"

	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
	"github.com/GenesisAN/illusionsCard/util"
)

func TestExtractKKStudioCharacters(t *testing.T) {
	const path = "./KKTest/133439749_p29.png"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("optional KK Studio fixture is not present")
	}

	card, err := ReadCardFromPath(path)
	if err != nil {
		t.Fatalf("read KK Studio scene: %v", err)
	}
	scene, ok := card.(*KK.KKSceneCard)
	if !ok {
		t.Fatalf("expected *KK.KKSceneCard, got %T", card)
	}

	files := scene.ExtractCharaCards()
	if len(files) != len(scene.CharaCards) {
		t.Fatalf("extracted %d cards for %d characters", len(files), len(scene.CharaCards))
	}
	for index, file := range files {
		png, err := util.PngBytesRead(file)
		if err != nil {
			t.Fatalf("parse extracted character %d PNG: %v", index, err)
		}
		if png.Type != Base.CT_KK && png.Type != Base.CT_KKCSP {
			t.Fatalf("extracted character %d has type %q", index, png.Type)
		}
		extracted, err := ReadCard(png)
		if err != nil {
			t.Fatalf("read extracted character %d: %v", index, err)
		}
		chara, ok := extracted.(*KK.KKCharaCard)
		if !ok {
			t.Fatalf("expected extracted *KK.KKCharaCard, got %T", extracted)
		}
		original := scene.CharaCards[index]
		if chara.CharInfo == nil || original.CharInfo == nil ||
			chara.CharInfo.Lastname != original.CharInfo.Lastname ||
			chara.CharInfo.Firstname != original.CharInfo.Firstname ||
			chara.CharInfo.Nickname != original.CharInfo.Nickname {
			t.Fatalf("extracted character %d identity does not match source", index)
		}
		if len(chara.ExtendedList) != len(original.ExtendedList) {
			t.Fatalf("extracted character %d plugin count=%d, want %d", index, len(chara.ExtendedList), len(original.ExtendedList))
		}
	}

	first, err := scene.ExtractCharaCard(0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, files[0]) {
		t.Fatal("single-card extraction differs from bulk extraction")
	}
	if _, err := scene.ExtractCharaCard(len(files)); err == nil {
		t.Fatal("out-of-range extraction did not return an error")
	}

	indexes := scene.UniqueCharaCardIndexes()
	if len(indexes) != 1 || indexes[0] != 0 {
		t.Fatalf("unique character indexes=%v, want [0]", indexes)
	}
	unique := scene.ExtractUniqueCharaCards()
	if len(unique) != 1 {
		t.Fatalf("extracted %d unique characters, want 1", len(unique))
	}
	if !bytes.Equal(unique[0], files[0]) {
		t.Fatal("unique extraction did not retain the first character instance")
	}
}
