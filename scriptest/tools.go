package scriptest

import (
	"errors"
	"fmt"
	"strconv"
)

// GetTLV to split string message into Tag Length Value
func GetTLV(s string) (res map[string]string, err error) {
	var tag string
	var leng int
	res = make(map[string]string)

	for {
		if len(s) == 0 {
			break
		} else {
			if len(s) <= 4 {
				return nil, errors.New("TLV should have first 4 alphanumeric")
			}
			tag = s[0:2]
			leng, err = strconv.Atoi(s[2:4])
			if err != nil {
				return nil, err
			}
			if len(s) < 4+leng {
				return nil, errors.New("TLV Not enough length for tag [" + tag + "]" + " need length " + strconv.Itoa(leng) + ".")
			}
			res[tag] = s[4 : 4+leng]
			s = s[4+leng:]
		}
	}
	return res, nil
}

// SetTLV to pack Tag Length Value into one string
func SetTLV(res map[string]string) (s string, err error) {
	for k, v := range res {
		s += k
		p := len(v)
		if p >= 100 {
			return s, errors.New("Data Tag should under 100 length")
		}
		s += fmt.Sprintf("%02d", p)
		s += v
	}
	return s, err
}
