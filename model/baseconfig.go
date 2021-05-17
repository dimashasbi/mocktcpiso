package model

import (
	"net"
	"time"
)

var (
	// FormatTimeLog for log HHMmmss
	FormatTimeLog = "150405"

	// FormatTimeforLog format for Logging File
	FormatTimeforLog = "2006-01-02 15:04:05.000"
)

type (
	// TCPconnection struct for server host
	TCPconnection struct {
		Host       string
		RemotePort int
		AsServer   bool
		ListenPort int
	}

	// BaseConfiguration for Running Service and Base Info
	BaseConfiguration struct {
		// Config Global
		ServiceID  string
		TimeoutTrx int
		PoolSize   int
		// config conditional
		TrxSession *TransactionTCP
		TCPConn    *net.TCPConn
	}

	// TransactionTCP is context for TCP handling session
	TransactionTCP struct {
		Isokey map[string]VarTCP
	}

	// VarTCP is struct to send for map
	VarTCP struct {
		Message string
		Time    time.Time
	}
)
