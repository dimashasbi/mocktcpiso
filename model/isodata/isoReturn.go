package isodata

import (
	"errors"
	"fmt"
)

// CreateReturnISO to create Simple Response Message Error
func CreateReturnISO(i Message, RC string) (Message, error) {
	// create msg
	i.MTI = "0210"
	i.IsoMessageMap[39] = RC
	i.IsoMessageMap[123] = "00000000000000000000"
	err := i.PackISO()
	if err != nil {
		errinf := fmt.Sprintf("Error Pack ISO %+v", err)
		return i, (errors.New(errinf))
	}
	return i, nil
}
