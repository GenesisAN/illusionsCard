// Package KK parses Koikatsu character, coordinate, and Studio scene cards.
package KK

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/GenesisAN/illusionsCard/Base"
	util "github.com/GenesisAN/illusionsCard/util"
)

var sceneExtendedMarker = append([]byte{4}, []byte(Base.KKEx)...)

// KKSceneCard represents a Koikatsu Studio scene card. The embedded Base.Card
// contains the union of scene-level and embedded-character plugin dependencies,
// so callers can use the same MOD inspection API as they do for character cards.
type KKSceneCard struct {
	*Base.Card
	PngData           []byte                        `json:"-"`
	Version           string                        `json:"version"`
	SourceFileName    string                        `json:"source_file_name"`
	CharaCards        []*KKCharaCard                `json:"chara_cards"`
	CharaCardFiles    [][]byte                      `json:"-"`
	SceneExtended     map[string]*Base.PluginData   `json:"-"`
	SceneExtendedList map[string]*Base.PluginDataEx `json:"scene_extended_list"`
	ExtendedVersion   int                           `json:"extended_version"`
}

// KKCharaCard2 is kept as a source-compatible alias for early scene prototypes.
type KKCharaCard2 = KKSceneCard

// CompareMissingZipMods returns the zipmod GUIDs referenced by this scene but
// absent from localGUIDs.
func (c *KKSceneCard) CompareMissingZipMods(localGUIDs []string) []string {
	local := make(map[string]struct{}, len(localGUIDs))
	for _, guid := range localGUIDs {
		local[guid] = struct{}{}
	}

	var missing []string
	for _, guid := range c.GetZipmodsDependencies() {
		if _, ok := local[guid]; !ok {
			missing = append(missing, guid)
		}
	}
	return dedupeStrings(missing)
}

// ExtractCharaCard returns one embedded character as a standalone KK character
// card PNG. The returned bytes are a copy and can safely be modified or written
// to disk by the caller.
func (c *KKSceneCard) ExtractCharaCard(index int) ([]byte, error) {
	if c == nil {
		return nil, errors.New("KK scene: nil card")
	}
	if index < 0 || index >= len(c.CharaCardFiles) {
		return nil, fmt.Errorf("KK scene: character index %d out of range [0,%d)", index, len(c.CharaCardFiles))
	}
	return append([]byte(nil), c.CharaCardFiles[index]...), nil
}

// ExtractCharaCards returns every embedded character as a standalone KK
// character card PNG, preserving the order in CharaCards.
func (c *KKSceneCard) ExtractCharaCards() [][]byte {
	if c == nil {
		return nil
	}
	result := make([][]byte, len(c.CharaCardFiles))
	for index, file := range c.CharaCardFiles {
		result[index] = append([]byte(nil), file...)
	}
	return result
}

// UniqueCharaCardIndexes returns the source indexes of distinct characters.
// Studio can contain multiple objects backed by the same character. Identity
// is based on character metadata plus the embedded character image, rather than
// the complete payload, because plugins may store different per-instance state.
func (c *KKSceneCard) UniqueCharaCardIndexes() []int {
	if c == nil {
		return nil
	}
	seen := make(map[[sha256.Size]byte]struct{}, len(c.CharaCards))
	result := make([]int, 0, len(c.CharaCards))
	for index, chara := range c.CharaCards {
		key, ok := charaIdentityKey(chara)
		if !ok {
			// Without enough identity data, retain the instance rather than risk
			// merging two unrelated characters.
			result = append(result, index)
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, index)
	}
	return result
}

// ExtractUniqueCharaCards returns one standalone PNG for each distinct
// character while preserving the first occurrence order from CharaCards.
func (c *KKSceneCard) ExtractUniqueCharaCards() [][]byte {
	if c == nil {
		return nil
	}
	indexes := c.UniqueCharaCardIndexes()
	result := make([][]byte, 0, len(indexes))
	for _, index := range indexes {
		if index >= len(c.CharaCardFiles) {
			continue
		}
		result = append(result, append([]byte(nil), c.CharaCardFiles[index]...))
	}
	return result
}

func charaIdentityKey(card *KKCharaCard) ([sha256.Size]byte, bool) {
	if card == nil || card.CharInfo == nil || card.Image2 == nil || len(*card.Image2) == 0 {
		return [sha256.Size]byte{}, false
	}
	hash := sha256.New()
	fmt.Fprintf(hash, "%d\x00%s\x00%s\x00%s\x00%s\x00",
		card.CharInfo.Sex,
		card.CharInfo.Version,
		card.CharInfo.Lastname,
		card.CharInfo.Firstname,
		card.CharInfo.Nickname,
	)
	hash.Write(*card.Image2)
	var result [sha256.Size]byte
	copy(result[:], hash.Sum(nil))
	return result, true
}

// KKSceneReader plugs Studio scenes into the common card factory.
type KKSceneReader struct{}

func (KKSceneReader) Read(pgb *util.PngBuff) (Base.CardInterface, error) {
	if pgb.Type != Base.CT_KKSC {
		return nil, errors.New("KKSceneReader: invalid type " + pgb.Type)
	}
	return ParseKKSceneCard(pgb)
}

// ParseKKSceneCard parses the scene header and ExtendedSave trailer. It also
// extracts valid character-card payloads embedded in the scene and aggregates
// their dependencies into the returned scene card.
func ParseKKSceneCard(pgb *util.PngBuff) (*KKSceneCard, error) {
	if pgb == nil || pgb.B == nil {
		return nil, errors.New("KK scene: nil buffer")
	}
	if pgb.Type != Base.CT_KKSC {
		return nil, errors.New("KK scene: invalid type " + pgb.Type)
	}

	payload := append([]byte(nil), pgb.B.Bytes()...)
	versionReader := &util.PngBuff{B: bytes.NewBuffer(payload)}
	version, err := versionReader.StringRead()
	if err != nil {
		return nil, fmt.Errorf("KK scene: read version: %w", err)
	}
	if version == "" {
		return nil, errors.New("KK scene: empty version")
	}

	base := &Base.Card{
		CardType:     Base.CT_KKSC,
		LoadVersion:  version,
		Path:         pgb.FilePath,
		Image1:       pgb.Png1,
		Extended:     make(map[string]*Base.PluginData),
		ExtendedList: make(map[string]*Base.PluginDataEx),
	}
	card := &KKSceneCard{
		Card:              base,
		Version:           version,
		SourceFileName:    pgb.FilePath,
		CharaCards:        make([]*KKCharaCard, 0),
		CharaCardFiles:    make([][]byte, 0),
		SceneExtended:     make(map[string]*Base.PluginData),
		SceneExtendedList: make(map[string]*Base.PluginDataEx),
	}
	if pgb.Png1 != nil {
		card.PngData = append([]byte(nil), (*pgb.Png1)...)
	}

	extVersion, raw, parsed, err := parseSceneExtendedData(payload)
	if err != nil {
		return nil, err
	}
	card.ExtendedVersion = extVersion
	card.SceneExtended = raw
	card.SceneExtendedList = parsed
	mergePluginMaps(card.Extended, card.ExtendedList, raw, parsed)

	card.CharaCards, card.CharaCardFiles = parseEmbeddedSceneCharacters(payload, pgb)
	for _, chara := range card.CharaCards {
		mergePluginMaps(card.Extended, card.ExtendedList, chara.Extended, chara.ExtendedList)
	}

	return card, nil
}

func parseSceneExtendedData(payload []byte) (int, map[string]*Base.PluginData, map[string]*Base.PluginDataEx, error) {
	emptyRaw := make(map[string]*Base.PluginData)
	emptyParsed := make(map[string]*Base.PluginDataEx)
	searchEnd := len(payload)

	for searchEnd > 0 {
		idx := bytes.LastIndex(payload[:searchEnd], sceneExtendedMarker)
		if idx < 0 {
			break
		}
		headerEnd := idx + len(sceneExtendedMarker) + 8
		if headerEnd <= len(payload) {
			version := int(int32(binary.LittleEndian.Uint32(payload[idx+len(sceneExtendedMarker):])))
			length := int(int32(binary.LittleEndian.Uint32(payload[idx+len(sceneExtendedMarker)+4:])))
			dataEnd := headerEnd + length
			if version >= 0 && length > 0 && dataEnd == len(payload) {
				raw, parsed, err := decodePluginData(payload[headerEnd:dataEnd])
				if err == nil {
					return version, raw, parsed, nil
				}
			}
		}
		searchEnd = idx
	}

	return 0, emptyRaw, emptyParsed, nil
}

func decodePluginData(data []byte) (map[string]*Base.PluginData, map[string]*Base.PluginDataEx, error) {
	var packed Base.MapSArrayInterface
	rest, err := packed.UnmarshalMsg(data)
	if err != nil {
		return nil, nil, err
	}
	if len(rest) != 0 {
		return nil, nil, errors.New("unexpected bytes after MessagePack plugin map")
	}

	raw := make(map[string]*Base.PluginData, len(packed))
	parsed := make(map[string]*Base.PluginDataEx, len(packed))
	for name, fields := range packed {
		if len(fields) < 2 {
			continue
		}
		version, ok := integerValue(fields[0])
		if !ok {
			continue
		}
		plugin := &Base.PluginData{Version: version, Data: fields[1]}
		raw[name] = plugin
		parsed[name] = &Base.PluginDataEx{
			Version:             version,
			Name:                name,
			RequiredZipmodGUIDs: decodeResolveInfos(fields[1]),
		}
	}
	return raw, parsed, nil
}

func decodeResolveInfos(data interface{}) []Base.ResolveInfo {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	info, ok := dataMap["info"]
	if !ok {
		return nil
	}
	entries, ok := info.([]interface{})
	if !ok {
		return nil
	}

	result := make([]Base.ResolveInfo, 0, len(entries))
	for _, entry := range entries {
		encoded, ok := entry.([]byte)
		if !ok {
			continue
		}
		var resolve Base.ResolveInfo
		if rest, err := resolve.UnmarshalMsg(encoded); err == nil && len(rest) == 0 {
			result = append(result, resolve)
		}
	}
	return result
}

func integerValue(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		if uint64(int(v)) != v {
			return 0, false
		}
		return int(v), true
	default:
		return 0, false
	}
}

func parseEmbeddedSceneCharacters(payload []byte, scene *util.PngBuff) ([]*KKCharaCard, [][]byte) {
	type marker struct {
		cardType string
		bytes    []byte
	}
	markers := make([]marker, 0, 2)
	for _, cardType := range []string{Base.CT_KK, Base.CT_KKCSP} {
		encoded := []byte(cardType)
		if len(encoded) < 128 {
			markers = append(markers, marker{cardType: cardType, bytes: append([]byte{byte(len(encoded))}, encoded...)})
		}
	}

	var cards []*KKCharaCard
	var files [][]byte
	seenOffsets := make(map[int]struct{})
	for _, mark := range markers {
		searchFrom := 0
		for searchFrom < len(payload) {
			relative := bytes.Index(payload[searchFrom:], mark.bytes)
			if relative < 0 {
				break
			}
			offset := searchFrom + relative + len(mark.bytes)
			searchFrom = offset
			if _, exists := seenOffsets[offset]; exists {
				continue
			}

			reader := &util.PngBuff{
				B:        bytes.NewBuffer(payload[offset:]),
				FilePath: scene.FilePath,
				Png1:     scene.Png1,
				Type:     mark.cardType,
			}
			parsed, err := ParseKKChara(reader)
			if err != nil || parsed.CharParmeter == nil {
				continue
			}
			consumed := len(payload[offset:]) - reader.B.Len()
			if consumed <= 0 {
				continue
			}
			file, err := buildStandaloneCharaCard(&parsed, mark.bytes, payload[offset:offset+consumed], scene.Png1)
			if err != nil {
				continue
			}
			parsed.Path = scene.FilePath
			cards = append(cards, &parsed)
			files = append(files, file)
			seenOffsets[offset] = struct{}{}
		}
	}
	return cards, files
}

func buildStandaloneCharaCard(card *KKCharaCard, encodedType, payload []byte, fallbackImage *[]byte) ([]byte, error) {
	var image []byte
	if card != nil && card.Image2 != nil && isCompletePNG(*card.Image2) {
		image = *card.Image2
	} else if fallbackImage != nil && isCompletePNG(*fallbackImage) {
		image = *fallbackImage
	} else {
		return nil, errors.New("KK scene: embedded character has no usable PNG image")
	}

	result := make([]byte, 0, len(image)+4+len(encodedType)+len(payload))
	result = append(result, image...)
	var productNo [4]byte
	binary.LittleEndian.PutUint32(productNo[:], 100)
	result = append(result, productNo[:]...)
	result = append(result, encodedType...)
	result = append(result, payload...)
	return result, nil
}

func isCompletePNG(data []byte) bool {
	return bytes.HasPrefix(data, Base.PngStartChunk) && bytes.Contains(data, Base.PngEndChunk)
}

func mergePluginMaps(dstRaw map[string]*Base.PluginData, dstParsed map[string]*Base.PluginDataEx, raw map[string]*Base.PluginData, parsed map[string]*Base.PluginDataEx) {
	for name, plugin := range raw {
		if _, exists := dstRaw[name]; !exists {
			dstRaw[name] = plugin
		}
	}
	for name, plugin := range parsed {
		if existing, exists := dstParsed[name]; exists {
			existing.RequiredPluginGUIDs = dedupeStrings(append(existing.RequiredPluginGUIDs, plugin.RequiredPluginGUIDs...))
			existing.RequiredZipmodGUIDs = mergeResolveInfos(existing.RequiredZipmodGUIDs, plugin.RequiredZipmodGUIDs)
			continue
		}
		copyPlugin := *plugin
		copyPlugin.RequiredPluginGUIDs = append([]string(nil), plugin.RequiredPluginGUIDs...)
		copyPlugin.RequiredZipmodGUIDs = append([]Base.ResolveInfo(nil), plugin.RequiredZipmodGUIDs...)
		dstParsed[name] = &copyPlugin
	}
}

func mergeResolveInfos(left, right []Base.ResolveInfo) []Base.ResolveInfo {
	result := append([]Base.ResolveInfo(nil), left...)
	seen := make(map[string]struct{}, len(result))
	for _, info := range result {
		seen[info.GUID+"\x00"+info.Property] = struct{}{}
	}
	for _, info := range right {
		key := info.GUID + "\x00" + info.Property
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, info)
	}
	return result
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
