package persistence

import (
	"testing"
)

func TestSessionsInit(t *testing.T){
	manager, err := NewManager("pgm", 3600)
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
}
