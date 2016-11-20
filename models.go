package main

import (
	"fmt"
	"time"
)

type HookDatum struct {
	Method  string
	Headers map[string][]string
	Body    string
	Time    time.Time
}

const sessionsKey = "hooker/session_ids"

func sessionKey(guid string) string {
	return fmt.Sprintf("hooker/session/%s", guid)
}
