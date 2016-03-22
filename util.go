// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"fmt"
	"log"
	"os"
)

type Logger struct {
	prefix                    string
	logger                    *log.Logger
	InfoEnabled, DebugEnabled bool
}

func NewLogger(prefix string, debug bool) *Logger {
	l := &Logger{
		prefix: prefix,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		InfoEnabled: true,
		DebugEnabled: debug,
	}
	return l
}

func (l *Logger) Infof(f string, args ...interface{}) {
	if l.InfoEnabled {
		l.logger.Printf(l.prefix+": "+f, args...)
	}
}

func (l *Logger) Debugf(f string, args ...interface{}) {
	if l.DebugEnabled {
		l.logger.Printf(l.prefix+": "+f, args...)
	}
}

func (l *Logger) Errorf(f string, args ...interface{}) error {
	return fmt.Errorf(l.prefix+": "+f, args...)
}

// tokenListContainsValue returns true if the 1#token header with the given
// name contains token.
func tokenListContainsValue(header http.Header, name string, value string) bool {
	for _, v := range header[name] {
		for _, s := range strings.Split(v, ",") {
			if strings.EqualFold(value, strings.TrimSpace(s)) {
				return true
			}
		}
	}
	return false
}

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

func computeAcceptKey(challengeKey string) string {
	h := sha1.New()
	h.Write([]byte(challengeKey))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func generateChallengeKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}
