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
var manager persistence.Manager

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
	id = s.SessionId()
	if s.Values() == nil{
		t.Error("map value is nil")
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
	s, err := manager.SessionRead(id)
	if err != nil{
		t.Error(err)
		t.Error("session read test failed")
	}
	s.Set("username", "username")
	s, err = manager.SessionRead(s.SessionId())
	if s.Values()["username"] == nil{
		t.Error("secondary session read test failed")
	}
}

func sessionId() string{
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}