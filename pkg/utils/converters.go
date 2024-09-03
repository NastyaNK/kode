package utils

import (
	"encoding/json"
	. "myproject/pkg/models"
)

func ReMarshal(m, result interface{}) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return err
	}
	return nil
}

func ConvertSpellResult(result SpellCheckResult) SpellResultTO {
	return SpellResultTO{
		Word:        result.Word,
		Suggestions: result.S,
		Position: Position{
			Pos: result.Pos,
			Row: result.Row,
			Col: result.Col,
		},
	}
}

func ConvertAllSpellResults(results []SpellCheckResult) []SpellResultTO {
	var resultsTO []SpellResultTO
	for _, result := range results {
		resultsTO = append(resultsTO, ConvertSpellResult(result))
	}
	return resultsTO
}
