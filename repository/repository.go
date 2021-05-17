package repository

import "mockTCPISO/model/isodata"

// IRepositoryTCP forwarding Repo
type IRepositoryTCP interface {
	// RequestMessageTrx(traceID string, isoData isodata.Message) (isodata.Message, error)
	ForwardMessage(isoData isodata.Message) error
}
