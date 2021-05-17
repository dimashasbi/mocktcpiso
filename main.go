package main

import (
	"fmt"
	logg "github.com/jeanphorn/log4go"
	"mockTCPISO/delivery"
	"mockTCPISO/model"
	"mockTCPISO/model/isodata"
	"mockTCPISO/repository/tcpsend"
	"mockTCPISO/scriptest"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {

	Basemod := model.BaseConfiguration{
		ServiceID:  "Trigger ISO Application",
		PoolSize:   10,
		TimeoutTrx: 30,
		TrxSession: &model.TransactionTCP{
			Isokey: make(map[string]model.VarTCP)},
	}

	args := os.Args[1:]
	port, _ := strconv.Atoi(args[0])
	// init TCP connection as Client
	Tcpmod := model.TCPconnection{
		AsServer:   true,
		Host:       "localhost",
		ListenPort: port,
	}

	// Setup log file
	pathlist := []string{"./log", "./log/message", "./log/event", "./log/nmm", "./log/error"}
	for _, path := range pathlist {
		// create folder
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, os.ModePerm)
		}
	}
	logg.LoadConfiguration("./log4go.json")

	// Init TCP Connection
	if Tcpmod.AsServer { // as server
		logdata := fmt.Sprintf("TCP Connection as Server is Listening in %d \n", Tcpmod.ListenPort)
		fmt.Printf(logdata)
		delivery.TCPAsServer(&Basemod, Tcpmod)

	} else { // as client
		logdata := fmt.Sprintf("TCP Connection as Client is Remote to %s:%d \n", Tcpmod.Host, Tcpmod.RemotePort)
		fmt.Printf(logdata)
		delivery.TCPAsClient(&Basemod, Tcpmod)
	}

	respons := func(req isodata.Message) isodata.Message { return isodata.Message{} }

	//  for MPAY
	if Tcpmod.ListenPort == 24107 {
		respons = scriptest.ReceiveISO
	}

	// prepare Repo
	repoTCP := tcpsend.NewTCPService(&Basemod)
	var wg sync.WaitGroup
	wg.Add(1)
	time.Sleep(4)
	go delivery.StartTCPConnection(&Basemod, repoTCP, respons)
	wg.Wait()
}

// JustResp For MandiriVA respon
func JustResp(req isodata.Message) isodata.Message {
	fmt.Println("---------------")
	fmt.Printf("REQ : %v, %+v", req.MTI, req.IsoMessageMap)
	fmt.Println()
	de2 := req.IsoMessageMap[2]
	de11 := req.IsoMessageMap[11]
	de37 := req.IsoMessageMap[37]
	resp := req
	resp.AsString = "0210F23E40018A81801E000000000400000016409766090910005438000000002625000009090914350000961614350909200909106012060001290600012991299991009600912999911900116902                            NI LUH DARMI                  BR SEDANG, DS SEDANG                      SYARAF, POLI        UMUM / BAYAR SENDIRI2020-04-02 09:39:44 - 00:00:00            13600076438946034UNTUK INFORMASI HUB. RSD MANGUSADA0011004013400"
	err := resp.ParseISO()
	if err != nil {
		fmt.Printf("%v", err)
	}
	resp.IsoMessageMap[2] = de2
	resp.IsoMessageMap[11] = de11
	resp.IsoMessageMap[37] = de37
	err = resp.PackISO()
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("RESP : %v, %+v", resp.MTI, resp.IsoMessageMap)
	return resp
}
