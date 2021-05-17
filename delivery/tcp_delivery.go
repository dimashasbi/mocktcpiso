package delivery

import (
	"errors"
	"fmt"
	"mockTCPISO/model"
	"mockTCPISO/model/isodata"
	"mockTCPISO/repository"
	"net"
	"os"
	"strconv"
	"time"
	// "github.com/panjf2000/ants"
)

var (
	// ServiceID to show name of ID Service
	ServiceID string
	// InitNMM for
	InitNMM = true
	// InitScheduler for
	InitScheduler = true
	// TCPReady for
	TCPReady chan bool
)

type tcpdelivery struct {
	responfunc func(isodata.Message) isodata.Message
	baseconf   *model.BaseConfiguration
	send       repository.IRepositoryTCP
}

// TCPAsServer is Start Application if TCPconn.AsServer is true
func TCPAsServer(m *model.BaseConfiguration, tc model.TCPconnection) {
	TCPReady = make(chan bool)
	hosting := &net.TCPAddr{
		Port: tc.ListenPort,
	}
	l, err := net.ListenTCP("tcp", hosting)
	if err != nil {
		err = errors.New("Connection TCP is error : " + err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
	go func() {
		for {
			m.TCPConn, err = l.AcceptTCP()
			fmt.Println()
			fmt.Println("New Client Connection is Accepted")
			TCPReady <- true
			if err != nil {
				err = errors.New("Listen TCP is error : " + err.Error())
				fmt.Println(err)
				os.Exit(2)
			}
			time.Sleep(1 * time.Second)
		}
	}()

}

// TCPAsClient is Start Application if TCPconn.AsServer is true
func TCPAsClient(m *model.BaseConfiguration, tc model.TCPconnection) {
	TCPReady = make(chan bool)
	var err error
	hosting := &net.TCPAddr{
		IP:   net.ParseIP(tc.Host),
		Port: tc.RemotePort,
	}
	m.TCPConn, err = net.DialTCP("tcp", nil, hosting)
	if err != nil {
		infLog := fmt.Sprintf("Connection TCP is error : %+v", err)
		fmt.Println(infLog)
		os.Exit(2)
	}
	go func() {
		fmt.Println("New Server Connection is Accepted")
		TCPReady <- true
	}()
	time.Sleep(1 * time.Second)
}

// StartTCPConnection for start TCP connection
func StartTCPConnection(m *model.BaseConfiguration, repo repository.IRepositoryTCP, respon func(isodata.Message) isodata.Message) {
	ServiceID = m.ServiceID
	delivery := tcpdelivery{
		baseconf:   m,
		responfunc: respon,
		send:       repo,
	}

	// create triggered Activation tcp
	// poolingConnection, _ := ants.NewPool(m.PoolSize)

	for {
		<-TCPReady
		fmt.Println()
		fmt.Println("TCP Connection is Ready")

		buff1 := make([]byte, 1)
		buff2 := make([]byte, 1)
		time.Sleep(2 * time.Second)
		for {
			// read Input
			// get first byte
			_, err := m.TCPConn.Read(buff1)
			if err != nil {
				err = errors.New("Error TCP Connection when Read First byte : " + err.Error())
				fmt.Print(err)
				break // testing
			}
			lenbit1, _ := strconv.ParseInt(fmt.Sprintf("%x", buff1), 16, 64)
			// get second byte
			_, err = m.TCPConn.Read(buff2)
			if err != nil {
				err = errors.New("Error TCP Connection when Read Second byte : " + err.Error())
				fmt.Print(err)
				break // testing
			}
			lenbit2, err := strconv.ParseInt(fmt.Sprintf("%x", buff2), 16, 64)
			if err != nil {
				err = errors.New("Error Read Length Byte : " + err.Error())
				fmt.Print(err)
			}
			isolen := lenbit1*256 + lenbit2
			if isolen > 4096 {
				isolen = 4096
			}
			if isolen <= 0 {
				break
			}
			bufferisomessage := make([]byte, isolen)
			_, err = m.TCPConn.Read(bufferisomessage)
			// fmt.Println(c)
			if err != nil {
				infoLog := fmt.Sprintf("Error TCP Connection when Read ISO message %+v", err)
				fmt.Print(infoLog)
			}

			// input to goroutine function
			timeIn := time.Now()

			go delivery.messageReceived(bufferisomessage, timeIn)

			// task := func(isoMsgBuffer []byte, timeIn time.Time) func() {
			// 	return func() {
			// 		delivery.messageReceived(bufferisomessage, timeIn)
			// 	}
			// }
			// poolingConnection.Submit(task(bufferisomessage, timeIn))
		}
		time.Sleep(2 * time.Second)
	}
}

// messageReceived to Follow up Message come from TCP Connection
func (d *tcpdelivery) messageReceived(isoMsgBuffer []byte, timeIn time.Time) {
	// read MessageISO
	isoMsg := fmt.Sprintf("%s", isoMsgBuffer)

	// Parsed ISO immadiately
	isoMessage := isodata.Message{
		AsString: isoMsg,
	}
	if err := isoMessage.ParseISO(); err != nil {
		errInf := fmt.Sprintf("Error Parsing ISO : %v", err)
		fmt.Printf(errInf)
		return
	}
	if isoMessage.MTI == "0200" || isoMessage.MTI == "0400" {
		resp := d.responfunc(isoMessage)
		err := d.send.ForwardMessage(resp)
		if err != nil {
			fmt.Printf("%+v", err)
		}
	} else if isoMessage.MTI == "0800" {
		fmt.Println()
		fmt.Printf("REQ : %v, %+v", isoMessage.MTI, isoMessage.IsoMessageMap)
		fmt.Println()
		isoMessage.MTI = "0810"
		isoMessage.IsoMessageMap[39] = "00"
		isoMessage.PackISO()
		fmt.Printf("RESP : %v, %+v", isoMessage.MTI, isoMessage.IsoMessageMap)
		fmt.Println()
		d.send.ForwardMessage(isoMessage)
	}
}

// ShowParsedISO to show Data
func showParsedISO(s string) {
	isoData := isodata.Message{
		AsString: s,
	}
	err := isoData.ParseISO()
	if err != nil {
		errorInf := errors.New("Error Parsing ISO : " + err.Error())
		fmt.Printf("%+v", errorInf)
	}
	infoLog := fmt.Sprintf("Parsed Iso Data : %+v", isoData.IsoMessageMap)
	fmt.Printf(infoLog)
}
