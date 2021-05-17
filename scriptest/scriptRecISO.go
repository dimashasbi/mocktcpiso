package scriptest

import (
	"fmt"
	"mockTCPISO/model/isodata"
)

// ReceiveISO For ReceiveISO respon
func ReceiveISO(req isodata.Message) isodata.Message {
	packing := func(req isodata.Message) isodata.Message {
		err := req.PackISO()
		if err != nil {
			fmt.Printf("%+v", err)
		}
		fmt.Printf("RESP : %v, %+v", req.MTI, req.IsoMessageMap)
		fmt.Println()
		return req
	}

	fmt.Println("---------------")
	fmt.Printf("REQ : %v, %+v", req.MTI, req.IsoMessageMap)
	fmt.Println()

	// ------------------------------ EDIT RESPONSE START
	if req.MTI == "0200" {
		req.MTI = "0210"
	} else if req.MTI == "0400" {
		req.MTI = "0410"
		req.IsoMessageMap[39] = "00"
		packing(req)
		return req
	}

	// delete(req.IsoMessageMap, 18)
	if req.IsoMessageMap[3] == "381000" {
		req.IsoMessageMap[39] = "00"
		req.IsoMessageMap[48] += "TA12007000000020BN30NASABAH TOYOTA ALPHARD        X10545678X00590123DD1715-DESEMBER-2020 "
	} else if req.IsoMessageMap[3] == "181000" {
		req.IsoMessageMap[39] = "00"
		req.IsoMessageMap[48] += "BR08billreff"
	}

	// ------------------------------ EDIT RESPONSE STOP
	req = packing(req)
	return req
}
