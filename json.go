package jacl

import (
	"encoding/json"
)

// ------------------------------------------------------------
// TO-FROM-JSON

// toFromJson() takes a list of pairs, where each pair is an
// input and output, and converts the input to json and then
// to the output.
func toFromJson(pairs ...interface{}) error {
	var input interface{} = nil
	for _, p := range pairs {
		if input != nil {
			b, err := json.Marshal(input)
			if err != nil {
				return err
			}
			err = json.Unmarshal(b, p)
			if err != nil {
				return err
			}
			input = nil
		} else {
			input = p
		}
	}
	return nil
}
