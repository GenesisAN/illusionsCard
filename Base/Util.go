package Base

func ParsePluginData(raw MapSArrayInterface) (map[string]*PluginData, map[string]*PluginDataEx) {
	exData := make(map[string]*PluginData)
	exDataEx := make(map[string]*PluginDataEx)

	for S, v := range raw {
		if v != nil && len(v) >= 2 {
			var pd PluginData
			if ver, ok := v[0].(int64); ok {
				pd.Version = int(ver)
			}
			pd.Data = v[1]
			exData[S] = &pd

			dex := pd.DeserializeObjects()
			dex.Name = S
			dex.Version = pd.Version
			exDataEx[S] = &dex
		}
	}

	return exData, exDataEx
}
