package proxy

import (
	"fmt"
	"log"
)

type MessageType int

const (
	Server MessageType = iota
	Error
	Warn
)

func printLog(messageType MessageType, message ...any) {
	switch messageType {
	case Server:
		log.Println(append([]any{Green + Bold + "[Caching Proxy]" + Reset}, message...)...)
	case Error:
		log.Fatalln(append([]any{Red + Bold + "[Caching Proxy]" + Reset}, message...)...)
	default:
		log.Println(append([]any{Yellow + Bold + "[Caching Proxy]" + Reset}, message...)...)
	}

}

func formatRequestLog(url string, milliseconds int64) string {
	var color string
	switch {
	case milliseconds > 1000:
		color = Red
	case milliseconds > 500:
		color = Yellow
	default:
		color = Blue
	}
	return fmt.Sprintf("---> Request: %s (%s%stook %dms%s)", url, Bold, color, milliseconds, Reset)
}

func WarnLog(message ...any) {
	printLog(Server, message...)
}

func ServerLog(message ...any) {
	printLog(Server, message...)
}

func ErrorLog(message ...any) {
	printLog(Error, message...)
}
