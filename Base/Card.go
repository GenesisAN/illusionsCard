package Base

// Card is illusionCard Bset struct
type Card struct {
	Extended     map[string]*PluginData //
	ExtendedList map[string]*PluginDataEx
	CharInfo     *ChaFileParameterEx
	Image1       *[]byte
	Image2       *[]byte
	CardType     string
	LoadVersion  string
	Path         string
}

type ChaFileParameterEx struct {
	Version   string
	Lastname  string
	Firstname string
	Nickname  string
}
