package isodata

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type (
	// IsoPackage for PackagerISO message
	IsoPackage struct {
		IsoHeader     bool
		IsoPackHeader int
		IsoPackager   []int
		IsoKey        []int
	}

	// Message for Component Message
	Message struct {
		IsoHeader       string         // iso header
		MTI             string         // MTI
		firstBitmap     string         // First bitmap
		actSecondBitmap bool           // Second Bitmap activation
		secondBitmap    string         // Second bitmap
		AsString        string         // ASCII full data
		ErrorIso        error          // info Error to take
		IsoMessageMap   map[int]string // Get Element of Iso Message in Map
		BitActive       []int
	}
)

var (
	isoHeader     = false
	isoPackHeader = 12
	isoKey        = []int{2, 7, 37}

	// IsoPack for
	IsoPack = []int{ // Default Package ISO if no Package
		//1  2  3   4   5   6   7  8  9  10
		16, -2, 6, 10, 12, 12, 10, 8, 8, 8, // 01..10
		6, 4, 4, 4, 4, 4, 4, 4, 4, 4, // 11..20
		3, 3, 4, 4, 2, 2, 3, 9, 3, 3, // 21..30
		-3, -3, -2, -2, -2, -2, 12, 6, 2, -2, // 31..40
		12, 15, 30, -2, -2, -3, -3, -2, 3, 3, // 41..50
		3, 16, 16, -3, -3, -3, -3, -3, -3, -3, // 51..60
		-2, -3, -3, 16, 1, 1, 2, 3, 3, 3, // 61..70
		4, 4, 6, 10, 10, 10, 10, 10, 10, 10, // 71..80
		12, 12, 12, 12, 12, 16, 16, 16, 16, 42, // 81..90
		1, 2, 5, 7, 42, 16, 16, 25, -2, -2, // 91..100
		-2, -2, -2, -3, -3, -3, -3, -3, -3, -3, // 101..110
		-3, -3, -3, -3, -3, -3, -3, -3, -3, -3, // 111..120
		-3, -3, -3, -3, -3, -3, -3, 16, // 121..128
	}
)

// SetPackageIso for set package Iso Message length
func (i *IsoPackage) SetPackageIso() {
	isoHeader = i.IsoHeader
	isoPackHeader = i.IsoPackHeader
	IsoPack = i.IsoPackager
	isoKey = i.IsoKey
}

// CreateTraceISO is to Trace ISO Message
func CreateTraceISO(isoData Message) string {
	var result string
	// combination ISO => PAN , STAN, RRN, TransmissionDateTime
	isomap := isoData.IsoMessageMap
	for _, v := range isoKey {
		result += isomap[v]
	}
	return result
}

// ParseISO to Single Data Element
func (m *Message) ParseISO() error {
	// Check the message is ISO or not
	isomsg, err := m.checkMessageISO()
	if err != nil {
		err = errors.New("Message is not ISO : " + err.Error())
		return err
	}

	leng := 0
	IsoMessage := make(map[int]string)
	data := []rune(isomsg)

	// get MTI
	m.MTI = string(data[0:4])
	copy(data, data[4:])
	// get First Bitmap
	m.firstBitmap = string(data[0:16])
	copy(data, data[16:])
	// get Second Bitmap
	FirstBitmap, err := strconv.ParseUint(m.firstBitmap, 16, 64)
	if err != nil {
		err = errors.New("Error First Bitmap : " + err.Error())
		return err
	}
	firstBitmap := fmt.Sprintf("%016b", FirstBitmap)
	firstBitmap = fmt.Sprintf("%064s", firstBitmap)
	firstbmp := []rune(string(firstBitmap))

	for i, v := range firstbmp {
		if i == 0 && string(v) == "1" {
			m.secondBitmap = string(data[0:16])
			copy(data, data[16:])
			m.actSecondBitmap = true
		} else if string(v) == "1" {
			if IsoPack[i] > 0 {
				if leng = IsoPack[i]; leng == 0 {
					IsoMessage[i+1] = ""
					continue
				}
				IsoMessage[i+1] = string(data[0:leng])
				copy(data, data[leng:])
			} else if IsoPack[i] == -2 {
				leng, err = strconv.Atoi(string(data[0:2]))
				if err != nil {
					errInfo := fmt.Sprintf("Error in Parsing bit %v, error : %+v \n", i+1, err)
					err = errors.New(errInfo)
					return err
				} else if leng == 0 {
					IsoMessage[i+1] = ""
					copy(data, data[2:])
					continue
				}
				IsoMessage[i+1] = string(data[2 : leng+2])
				copy(data, data[leng+2:])
			} else if IsoPack[i] == -3 {
				leng, err = strconv.Atoi(string(data[0:3]))
				if err != nil {
					errInfo := fmt.Sprintf("Error in Parsing bit %v, error : %+v \n", i+1, err)
					err = errors.New(errInfo)
					return err
				} else if leng == 0 {
					IsoMessage[i+1] = ""
					copy(data, data[3:])
					continue
				}
				IsoMessage[i+1] = string(data[3 : leng+3])
				copy(data, data[leng+3:])
			}
		}
	}

	if m.actSecondBitmap {
		SecondBitmap, err := strconv.ParseUint(m.secondBitmap, 16, 64)
		if err != nil {
			errInfo := fmt.Sprintf("Error Second Bitmap : %+v \n", err)
			err = errors.New(errInfo)
			return err
		}
		secondBitmap := fmt.Sprintf("%016b", SecondBitmap)
		secondBitmap = fmt.Sprintf("%064s", secondBitmap)
		secondbmp := []rune(secondBitmap)
		if m.secondBitmap != "" {
			for j, val := range secondbmp {
				if string(val) == "1" {
					if IsoPack[j+64] > 0 {
						if leng = IsoPack[j+64]; leng == 0 {
							IsoMessage[j+65] = ""
							continue
						}
						IsoMessage[j+65] = string(data[0:leng])
						copy(data, data[leng:])
					} else if IsoPack[j+64] == -2 {
						leng, err = strconv.Atoi(string(data[0:2]))
						if err != nil {
							errInfo := fmt.Sprintf("Error in Parsing bit %v, error : %+v \n", j+65, err)
							err = errors.New(errInfo)
							return err
						} else if leng == 0 {
							IsoMessage[j+65] = ""
							copy(data, data[2:])
							continue
						}
						IsoMessage[j+65] = string(data[2 : leng+2])
						copy(data, data[leng+2:])
					} else if IsoPack[j+64] == -3 {
						leng, err = strconv.Atoi(string(data[0:3]))
						if err != nil {
							errInfo := fmt.Sprintf("Error in Parsing bit %v, error : %+v\n", j+65, err)
							err = errors.New(errInfo)
							return err
						} else if leng == 0 {
							IsoMessage[j+65] = ""
							copy(data, data[3:])
							continue
						}
						IsoMessage[j+65] = string(data[3 : leng+3])
						copy(data, data[leng+3:])
					}
				}
			}
		}
	}
	m.IsoMessageMap = IsoMessage
	// check active bit
	m.GetKeyActive()
	return nil
}

// GetKeyActive to get Bit Active in Message ISO
func (m *Message) GetKeyActive() error {
	keys := make([]int, 0, len(m.IsoMessageMap))
	for k := range m.IsoMessageMap {
		keys = append(keys, k)
	}
	m.BitActive = keys
	sort.Ints(m.BitActive)
	return nil
}

// PackISO to get Create Message ISO in string
func (m *Message) PackISO() error {
	var result string

	// check active bit
	m.GetKeyActive()

	// check all variable needed is not empty (MTI)
	if m.MTI == "" {
		return errors.New("No MTI to Pack")
	}

	// check variable length should be same with packager

	// check bitmap active
	keyActive := make([]int, len(m.BitActive))
	copy(keyActive, m.BitActive)
	bitmap := "0"
	for i := 2; i <= 128; i++ {
		for a, v := range keyActive {
			if v == i {
				bitmap += "1"
				keyActive = append(keyActive[:a], keyActive[a+1:]...)
				break
			}
		}
		if len(bitmap) == i {
			continue
		}
		bitmap += "0"
	}

	if len(bitmap) == 128 {
		bitmap = "1" + bitmap[1:]
		m.actSecondBitmap = true
		secondbmp := bitmap[64:128]

		// get second bitmap
		SecondBmp, err := strconv.ParseUint(secondbmp, 2, 64)
		if err != nil {
			return errors.New("Error Parse Uint")
		}
		m.secondBitmap = strings.ToUpper(fmt.Sprintf("%016x", SecondBmp))
	}

	firstbmp := bitmap[:64]
	// get firstbitmap
	FirstBmp, err := strconv.ParseUint(firstbmp, 2, 64)
	if err != nil {
		return errors.New("Error Parse Uint")
	}
	m.firstBitmap = strings.ToUpper(fmt.Sprintf("%016x", FirstBmp))

	// Append message
	// mti + first bitmap + ifhave(secondbitmap) + message
	result = m.MTI + m.firstBitmap
	if m.actSecondBitmap {
		result += m.secondBitmap
	}

	for _, i := range m.BitActive {
		leng := 0
		if IsoPack[i-1] > 0 {
			result += m.IsoMessageMap[i]
		} else if IsoPack[i-1] == -2 {
			leng = len(m.IsoMessageMap[i])
			l := strconv.Itoa(leng)
			result += fmt.Sprintf("%02v", l) + m.IsoMessageMap[i]
		} else if IsoPack[i-1] == -3 {
			leng = len(m.IsoMessageMap[i])
			l := strconv.Itoa(leng)
			result += fmt.Sprintf("%03v", l) + m.IsoMessageMap[i]
		}
	}

	m.AsString = result
	return nil
}

//  checkMessageISO for checking Message, MTI or  assume there is ISO Header
func (m *Message) checkMessageISO() (string, error) {
	data := []rune(m.AsString)
	// check ISO header
	if isoHeader {
		data = data[isoPackHeader:]
		return string(data), nil
	} else if string(data[0:3]) == "ISO" {
		data = data[isoPackHeader:]
		return string(data), nil
	}
	// check MTI
	mti := string(data[0:4])
	listmti := []string{"0200", "0210", "0400", "0410", "0800", "0810"}
	for _, list := range listmti {
		if mti == list {
			return m.AsString, nil
		}
		continue
	}
	return m.AsString, errors.New("Not Default MTI")
}
