package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const LogKey = "log"

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

func Call() *LogPayload {
	var payload = LogPayload{}

	payload.LogID = generateID()
	payload.Time = time.Now().Format(time.RFC3339)

	return &payload
}

func (l *LogPayload) TDR(ctx context.Context, message string) {
	l.LogType = "tdr"
	l.Message = message
	l.Done()
}

func (l *LogPayload) Info(ctx context.Context, message string) {
	l.LogType = "business info"
	l.Message = message
	l.Done()
}

func (l *LogPayload) Done() {
	_, file, line, _ := runtime.Caller(2)
	l.Caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	jsonPayload, _ := json.Marshal(l)
	fmt.Printf("%s\n", string(jsonPayload))
}

func (l *LogPayload) SetError(err error) *LogPayload {
	if err != nil {
		l.Error = err.Error()
	}
	return l
}

func (l *LogPayload) SetMessage(ctx context.Context, message string) *LogPayload {
	l.Message = message
	return l
}

func (l *LogPayload) SetAdditionalInfo(key string, value interface{}) *LogPayload {
	if len(l.AdditionalData) < 1 {
		l.AdditionalData = make(map[string]interface{})
	}
	l.AdditionalData[key] = value

	return l
}

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
