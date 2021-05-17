package logfile

import (
	"fmt"
	"mockTCPISO/model"
	"mockTCPISO/model/isodata"
	"sort"
	"time"

	logg "github.com/jeanphorn/log4go"
	"github.com/pkg/errors"
)

// LogInfo used for Logging Application Condition, Get New Configuration, etc
func LogInfo(info string) {
	logg.LOGGER("Event").Debug(info)
}

// LogNMM used for Logging Message Transaction pure
func LogNMM(traceID, flow, serviceID string, isoData *isodata.Message) {
	logg.LOGGER("NMM").Trace("[%v][%v][%v][%v][%v][%v]{%v}", time.Now().Format(model.FormatTimeforLog),
		traceID, flow, "SERVICE_"+serviceID, isoData.IsoMessageMap[39], "TCP", isoData.AsString)
}

// LogMessageTCP used for Logging Message Transaction pure
func LogMessageTCP(serviceID, flow string, isoData *isodata.Message) {
	raw := fmt.Sprintf("[%v]", isoData.MTI)
	sort.Ints(isoData.BitActive)
	for _, v := range isoData.BitActive {
		raw += fmt.Sprintf("[%v]", isoData.IsoMessageMap[v])
	}
	// [TIME][FLOW][ISO-Conn][RC]{RAW}
	logg.LOGGER("Message").Trace("[%v][%v][%v][%v]{%v}", time.Now().Format(model.FormatTimeforLog),
		flow, "ISO-"+serviceID, isoData.IsoMessageMap[39], raw)
}

// LogTrxCondition used for Transaction Condition
func LogTrxCondition(traceID, flow, condition string) {
	// [TIME][TRACEID][FLOW][CONDITION]
	logg.LOGGER("Message").Trace("[%v][%v][%v][%v]", time.Now().Format(model.FormatTimeforLog),
		traceID, flow, condition)
}

// LogError used for Logging Some Problem in Error Transaction or Application Configuration
// For Error in Transaction Process or Other Usecase
func LogError(err error) {
	logg.LOGGER("Error").Error("%+v", errors.Wrap(err, ""))
}

// LogErrorMessage used for Logging Some Problem in Error Transaction or Application Configuration
// For Error in Transaction Process or Other Usecase
func LogErrorMessage(info string, err error) {
	logg.LOGGER("Error").Error(info)
	logg.LOGGER("Error").Error("%+v", errors.Wrap(err, info))
}

// LogCritical used for Logging Base Application that critical, like Connection DB, TCP
func LogCritical(err error) {
	logg.LOGGER("Critical").Critical("%+v", errors.Wrap(err, "Critical Error"))
}
