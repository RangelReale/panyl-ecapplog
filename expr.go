package panylecapplog

import "github.com/RangelReale/panyl/v2"

func ExprConstants() map[string]any {
	return map[string]any{
		"MetadataExtraCategories": panyl.MetadataExtraCategories,
		"MetadataECAppLogColor":   MetadataColor,
		"MetadataECAppLogBGColor": MetadataBGColor,
	}
}
