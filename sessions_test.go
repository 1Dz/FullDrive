package main

import (
	"testing"
	"Conus/persistence"
	"io"
	"encoding/base64"
	"crypto/rand"
	"encoding/json"
)

var id string

func TestSessionsInit(t *testing.T){
	persistence.Init()
	manager := persistence.NewManager("pgm", 3600)
	s, err := manager.Driver.SessionInit(sessionId())
	if err != nil{
		t.Error(err)
	}
	id = s.SessionID()
	if s.Values() == nil{
		t.Error("init values map is nil")
	}
	t.Log(s.Values())
}

func TestUnmarshalMap(t *testing.T){
	m := make(map[string]interface{})
	m["a"] = "asd"
	js, _ := json.Marshal(m)
	result := make(map[string]interface{}, 0)
	json.Unmarshal(js, &result)
	if len(result) == 0{
		t.Failed()
	}
	t.Log(len(result))
	t.Log(result["a"])
}

func TestSessionsReadAndSet(t *testing.T){
	persistence.Init()
	manager := persistence.NewManager("pgm", 3600)
	s, err := manager.Driver.SessionRead(id)
	if err != nil{
		t.Error(err)
		t.Error("session read test failed")
	}
	err = s.Set("username", "username")
	if err != nil{
		t.Error("error while setting into session")
		t.Error(err)
	}
	s, err = manager.Driver.SessionRead(s.SessionID())
	m := s.Values()["username"]
	if m == nil{
		t.Error("READ VALUES map is nil")
	}
	t.Log(m)
}

func sessionId() string{
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}