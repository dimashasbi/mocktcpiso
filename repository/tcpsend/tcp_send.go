package tcpsend

import (
	"mockTCPISO/model"
	"mockTCPISO/model/isodata"
	"mockTCPISO/repository"
	"encoding/hex"
	"errors"
	"fmt"
)

type tcpsend struct {
	baseconf *model.BaseConfiguration
	tcptrx   *model.TransactionTCP
}

var (
	// ServiceID for application ID Service
	ServiceID string
)

// NewTCPService Create Repository / Service for send Message via TCP
func NewTCPService(basemod *model.BaseConfiguration) repository.IRepositoryTCP {
	ServiceID = basemod.ServiceID
	return &tcpsend{
		baseconf: basemod,
	}
}

func (t *tcpsend) ForwardMessage(isoData isodata.Message) error {
	msg := isoData.AsString
	le := len(msg)

	l := fmt.Sprintf("%04x", le)
	leng := []byte(l)
	FullMessage := make([]byte, hex.DecodedLen(len(leng)))
	_, err := hex.Decode(FullMessage, leng)
	if err != nil {
		errInfo := fmt.Sprintf("Error Decode when sending TCP : %v", err)
		err = errors.New(errInfo)
		return err
	}
	FullMessage = append(FullMessage, msg...)

	_, err = t.baseconf.TCPConn.Write(FullMessage)
	if err != nil {
		errInfo := fmt.Sprintf("Error Write when sending TCP : %v", err)
		err = errors.New(errInfo)
		return err
	}
	return nil
}
