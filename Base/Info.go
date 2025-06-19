package Base

var (
	PngEndChunk   = []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}
	PngStartChunk = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D}
)

// GUID
const (
	PluginGUID  = "com.bepis.bepinex.sideloader"
	UARExtID    = "com.bepis.sideloader.universalautoresolver"
	UARExtIDOld = "EC.Core.Sideloader.UniversalAutoResolver"
)

const KKEx = "KKEx"

// CardType
const (
	CT_KK    = "【KoiKatuChara】"    // Koikatu
	CT_AI    = "【AIS_Chara】"       // AiSyoujyo
	CT_EC    = "【EroMakeChara】"    // EmotionCreators
	CT_KKS   = "【KoiKatuCharaSun】" // KoikatsuSunshine
	CT_KKP   = "【KoiKatuCharaS】"   // KoikatsuParty
	CT_KKCSP = "【KoiKatuCharaSP】"  // KoikatsuPartySpecialPatch
	CT_KKC   = "【KoiKatuClothes】"  // KoiKatuClothes
)

const (
	CTI_Unknown = iota
	CTI_Koikatu
	CTI_KoikatsuParty
	CTI_KoikatsuPartySpecialPatch
	CTI_EmotionCreators
	CTI_AiSyoujyo
	CTI_KoikatsuSunshine
	CTI_KoiKatuClothes
	CTI_KoiKatuScene
)
