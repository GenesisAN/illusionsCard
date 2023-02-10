package Base

var PngEndChunk = []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

var PngStartChunk = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D}

const PluginGUID = "com.bepis.bepinex.sideloader"

const UARExtID = "com.bepis.sideloader.universalautoresolver"

const UARExtIDOld = "EC.Core.Sideloader.UniversalAutoResolver"
const Parameter = "Parameter"
const KKEx = "KKEx"

const CT_KK = "【KoiKatuChara】"     // Koikatu
const CT_AI = "【AIS_Chara】"        // AiSyoujyo
const CT_EC = "【EroMakeChara】"     // EmotionCreators
const CT_KKS = "【KoiKatuCharaSun】" // KoikatsuSunshine

const CT_KKP = "【KoiKatuCharaS】"    // KoikatsuParty
const CT_KKCSP = "【KoiKatuCharaSP】" // KoikatsuPartySpecialPatch
const CT_KKC = "【KoiKatuClothes】"   // KoiKatuClothes
