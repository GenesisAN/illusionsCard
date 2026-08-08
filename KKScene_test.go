package illusionCard

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
	"github.com/GenesisAN/illusionsCard/util"
)

func TestReadKKStudioScene(t *testing.T) {
	resolveBytes, err := (&Base.ResolveInfo{
		GUID:       "example.studio.item",
		Slot:       1,
		LocalSlot:  100000001,
		Property:   "no",
		CategoryNo: 10,
	}).MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	extended := Base.MapSArrayInterface{
		Base.UARExtID: {
			int64(1),
			map[string]interface{}{"info": []interface{}{resolveBytes}},
		},
	}
	packed, err := extended.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	payload := append([]byte{7}, []byte("1.0.4.2")...)
	payload = append(payload, 0, 0, 0, 0) // empty object list
	payload = append(payload, 4, 'K', 'K', 'E', 'x')
	payload = appendInt32(payload, 3)
	payload = appendInt32(payload, len(packed))
	payload = append(payload, packed...)
	file := append(append([]byte(nil), Base.PngStartChunk...), Base.PngEndChunk...)
	file = append(file, payload...)

	png, err := util.PngBytesRead(file)
	if err != nil {
		t.Fatal(err)
	}
	png.FilePath = "synthetic_scene.png"
	card, err := ReadCard(png)
	if err != nil {
		t.Fatal(err)
	}
	scene, ok := card.(*KK.KKSceneCard)
	if !ok {
		t.Fatalf("expected *KK.KKSceneCard, got %T", card)
	}
	assertKKStudioScene(t, scene, "1.0.4.2")
	if got := scene.GetZipmodsDependencies(); len(got) != 1 || got[0] != "example.studio.item" {
		t.Fatalf("unexpected dependencies: %v", got)
	}
}

func TestReadKKStudioSceneFixture(t *testing.T) {
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
	assertKKStudioScene(t, scene, "1.0.4.2")
	if len(scene.SceneExtendedList) == 0 {
		t.Fatal("scene ExtendedSave data was not parsed")
	}
	if len(scene.GetZipmodsDependencies()) == 0 {
		t.Fatal("scene zipmod dependencies were not parsed")
	}
	if len(scene.CharaCards) == 0 {
		t.Fatal("embedded character cards were not parsed")
	}

	t.Logf("scene version=%s, scene plugins=%d, scene zipmods=%d, embedded characters=%d",
		scene.Version, len(scene.SceneExtendedList), countZipmods(scene.SceneExtendedList), len(scene.CharaCards))
	for name, plugin := range scene.SceneExtendedList {
		t.Logf("scene plugin=%s version=%d zipmods=%v", name, plugin.Version, zipmodGUIDs(plugin.RequiredZipmodGUIDs))
	}
	for index, chara := range scene.CharaCards {
		lastname, firstname, nickname := "", "", ""
		sex := 0
		if chara.CharInfo != nil {
			lastname = chara.CharInfo.Lastname
			firstname = chara.CharInfo.Firstname
			nickname = chara.CharInfo.Nickname
			sex = chara.CharInfo.Sex
		}
		t.Logf("character[%d] name=%s %s nickname=%s sex=%d plugins=%d zipmods=%v",
			index, lastname, firstname, nickname, sex, len(chara.ExtendedList), chara.GetZipmodsDependencies())
		for name, plugin := range chara.ExtendedList {
			t.Logf("character[%d] plugin=%s version=%d zipmods=%v", index, name, plugin.Version, zipmodGUIDs(plugin.RequiredZipmodGUIDs))
		}
	}
}

func assertKKStudioScene(t *testing.T, scene *KK.KKSceneCard, version string) {
	t.Helper()
	if scene.CardType != Base.CT_KKSC {
		t.Fatalf("unexpected card type %q", scene.CardType)
	}
	if scene.TypeInt() != Base.CTI_KoiKatuScene {
		t.Fatalf("unexpected numeric card type %d", scene.TypeInt())
	}
	if scene.Version != version {
		t.Fatalf("unexpected scene version %q", scene.Version)
	}
	if scene.GetPath() == "" {
		t.Fatal("scene path was not retained")
	}
	if len(scene.CompareMissingMods(nil)) == 0 {
		t.Fatal("scene ResolveInfo missing-mod comparison returned no dependencies")
	}
	if len(scene.CompareMissingZipMods(nil)) == 0 {
		t.Fatal("scene missing-mod comparison returned no dependencies")
	}
}

func appendInt32(dst []byte, value int) []byte {
	var encoded [4]byte
	binary.LittleEndian.PutUint32(encoded[:], uint32(value))
	return append(dst, encoded[:]...)
}

func countZipmods(plugins map[string]*Base.PluginDataEx) int {
	count := 0
	for _, plugin := range plugins {
		count += len(plugin.RequiredZipmodGUIDs)
	}
	return count
}

func zipmodGUIDs(infos []Base.ResolveInfo) []string {
	result := make([]string, 0, len(infos))
	for _, info := range infos {
		result = append(result, info.GUID)
	}
	return result
}
