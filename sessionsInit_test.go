package main

import (
	"testing"
	"Conus/persistence"
	"io"
	"encoding/base64"
	"crypto/rand"
)

func TestSessionsInit(t *testing.T){
	persistence.Init()
	manager, err := persistence.NewManager("pgm", 3600)
	if err != nil{
		t.Error(err)
	}
	s, err := manager.SessionInit(sessionId())
	if err != nil{
		t.Error(err)
	}
	if s.Values() == nil{
		t.Error("map value is nil")
	}
	t.Log(s.Values())
}

func sessionId() string{
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}