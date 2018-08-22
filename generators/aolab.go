/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-08-2018
 * |
 * | File Name:     aolab.go
 * +===============================================
 */

package generators

import (
	"encoding/json"
	"fmt"

	"github.com/I1820/lanserver/models"
)

// AolabGenerator generates data based on
// lanserver protocol
// and aolab model.
type AolabGenerator struct {
	DevEUI string
}

// Topic returns lanserver mqtt topic
func (g AolabGenerator) Topic() []byte {
	return []byte(fmt.Sprintf("device/%s/rx", g.DevEUI))
}

// Generate generates lanserver message by converting input into cbor and generator
// parameters.
func (g AolabGenerator) Generate(input interface{}) ([]byte, error) {
	// input into json
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// lora message
	message, err := json.Marshal(models.RxMessage{
		DevEUI: g.DevEUI,
		Data:   b,
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}
