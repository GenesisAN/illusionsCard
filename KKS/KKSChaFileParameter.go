package KKS

//go:generate msgp
type KKSChaFileParameter struct {
	Version          string      `msg:"version"`
	Sex              int         `msg:"sex"`
	ExType           int         `msg:"exType"`
	Lastname         string      `msg:"lastname"`
	Firstname        string      `msg:"firstname"`
	Nickname         string      `msg:"nickname"`
	CallType         int         `msg:"callType"`
	Personality      int         `msg:"personality"`
	BloodType        int         `msg:"bloodType"`
	BirthMonth       int         `msg:"birthMonth"`
	BirthDay         int         `msg:"birthDay"`
	ClubActivities   int         `msg:"clubActivities"`
	VoiceRate        float64     `msg:"voiceRate"`
	WeakPoint        int         `msg:"weakPoint"`
	Awnser           Awnser      `msg:"awnser"`
	Denial           Denial      `msg:"denial"`
	Attribute        Attribute   `msg:"attribute"`
	Aggressive       int         `msg:"aggressive"`
	Diligence        int         `msg:"diligence"`
	Kindness         int         `msg:"kindness"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}

type Awnser struct {
	Animal           bool        `msg:"animal"`
	Eat              bool        `msg:"eat"`
	Cook             bool        `msg:"cook"`
	Exercise         bool        `msg:"exercise"`
	Study            bool        `msg:"study"`
	Fashionable      bool        `msg:"fashionable"`
	BlackCoffee      bool        `msg:"blackCoffee"`
	Spicy            bool        `msg:"spicy"`
	Sweet            bool        `msg:"sweet"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}
type Denial struct {
	Kiss             bool        `msg:"kiss"`
	Aibu             bool        `msg:"aibu"`
	Anal             bool        `msg:"anal"`
	Massage          bool        `msg:"massage"`
	NotCondom        bool        `msg:"notCondom"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}
type Attribute struct {
	Active           bool        `msg:"active"`
	Bitch            bool        `msg:"bitch"`
	Choroi           bool        `msg:"choroi"`
	Dokusyo          bool        `msg:"dokusyo"`
	Friendly         bool        `msg:"friendly"`
	Harapeko         bool        `msg:"harapeko"`
	Hinnyo           bool        `msg:"hinnyo"`
	Hitori           bool        `msg:"hitori"`
	Info             bool        `msg:"info"`
	Kireizuki        bool        `msg:"kireizuki"`
	LikeGirls        bool        `msg:"likeGirls"`
	Lonely           bool        `msg:"lonely"`
	Love             bool        `msg:"love"`
	Majime           bool        `msg:"majime"`
	Mutturi          bool        `msg:"mutturi"`
	Nakama           bool        `msg:"nakama"`
	Nonbiri          bool        `msg:"nonbiri"`
	Okute            bool        `msg:"okute"`
	Ongaku           bool        `msg:"ongaku"`
	Sinsyutu         bool        `msg:"sinsyutu"`
	Talk             bool        `msg:"talk"`
	Ukemi            bool        `msg:"ukemi"`
	Undo             bool        `msg:"undo"`
	ExtendedSaveData interface{} `msg:"ExtendedSaveData"`
}
