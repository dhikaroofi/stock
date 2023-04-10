package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// LogKey is the key used by logger to save in context
const LogKey = "log"

type LogKeyType string

// LogPayload for custom logger
type LogPayload struct {
	Time           string                 `json:"time"`
	LogType        string                 `json:"logType"`
	Caller         string                 `json:"caller"`
	LogID          string                 `json:"logID,omitempty"`
	TrxID          string                 `json:"trxID,omitempty"`
	Message        string                 `json:"message"`
	Error          string                 `json:"error"`
	AdditionalData map[string]interface{} `json:"additionalData,omitempty"`
}

func generateID() string {
	return uuid.New().String()
}

// Call is used for initiate the logger specially on transaction
func Call() *LogPayload {
	var payload = LogPayload{}

	payload.LogID = generateID()
	payload.Time = time.Now().Format(time.RFC3339)

	return &payload
}

// TDR this function is used for end the logger and printed into console
// TDR function has purpose to logging every incoming and outgoing transaction thats happen on the system
func (l *LogPayload) TDR(ctx context.Context, message string) {
	l.LogType = "tdr"
	l.Message = message
	l.done()
}

// Info this function is used for end the logger and printed into console,
// this function is has purpose to logging all of detail that correlated with TDR
func (l *LogPayload) Info(ctx context.Context, message string) {
	l.LogType = "business info"
	l.Message = message
	l.done()
}

func (l *LogPayload) done() {
	_, file, line, _ := runtime.Caller(2)
	l.Caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	jsonPayload, _ := json.Marshal(l)
	fmt.Printf("%s\n", string(jsonPayload))
}

// SetError this function is used for add error on logger
func (l *LogPayload) SetError(err error) *LogPayload {
	if err != nil {
		l.Error = err.Error()
	}
	return l
}

// SetAdditionalInfo this function is used for add additional info if needed for logger
func (l *LogPayload) SetAdditionalInfo(key string, value interface{}) *LogPayload {
	if len(l.AdditionalData) < 1 {
		l.AdditionalData = make(map[string]interface{})
	}
	l.AdditionalData[key] = value

	return l
}

// SysInfo this function is used for logging the information from internal system
func SysInfo(message string) {
	var payload = LogPayload{}
	_, file, line, _ := runtime.Caller(1)

	payload.LogType = "system info"
	payload.Time = time.Now().Format(time.RFC3339)
	payload.Caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	payload.Message = message

	jsonPayload, _ := json.Marshal(payload)
	fmt.Printf("%s\n", string(jsonPayload))
}

// Fatal is replacement for log.Fatalf
func Fatal(message string) {
	var payload = LogPayload{}
	_, file, line, _ := runtime.Caller(1)

	payload.LogType = "system info"
	payload.Time = time.Now().Format(time.RFC3339)
	payload.Caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	payload.Message = message

	jsonPayload, _ := json.Marshal(payload)
	fmt.Printf("%s\n", string(jsonPayload))
	os.Exit(1)
}
