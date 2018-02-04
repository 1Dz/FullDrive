package handlers

import (
	"sync"
	"time"
	"io"
	"encoding/base64"
	"crypto/rand"
	"errors"
	"Conus/persistence"
)

type Manager struct {
	sync.Mutex
	cookiesName string
	maxLifeTime int64
}

type Session struct {
	sid string
	timeAcceced time.Time
	values map[string]interface{}
}

type Provider interface{
	Get(key string, value interface{}) (interface{}, error)
	Set(key string, value interface{})
	Delete(key string) error
	SessionId()string
}

func NewManager(cookieName string, maxlifetime int64) (*Manager, error) {

	return &Manager{cookiesName:cookieName, maxLifeTime:maxlifetime}, nil
}

func (s *Session) Get(key string, value interface{}) (interface{}, error){
	res, ok := s.values[key]
	if ok{
		return res, nil
	}
	return nil, errors.New("No such element with key: " + key)
}

func (s *Session) Set(key string, value interface{}) {
	s.values[key] = value
}

func (s *Session) Delete(key string) error{
	if _, ok := s.values[key]; !ok{
		return errors.New("No such element with key: " + key)
	}
	delete(s.values, key)
	return nil
}

func (s *Session) SessionId() string{
	return s.sid
}

func (s *Session) TimeAcceced() time.Time{
	return s.timeAcceced
}

func (s *Session) Values() map[string]interface{}{
	return s.values
}

func (m *Manager) SessionInit(sid string) (Session, error){
	s := Session{sid, time.Now(), nil}
	err := persistence.SessionInit(&s)
	if err != nil{
		return Session{}, err
	}
	return s, nil
}

func (m *Manager) SessionRead(sid string) (*Session, error){
	s, err := persistence.SessionRead(sid)
	if err != nil{
		return nil, err
	}
	return s, nil
}

func (m *Manager) SessionDestroy(sid string) error{

}

func (m *Manager) SessionUpdate(sid string) error{

}

func (m *Manager) SessionGC() {
	m.Lock()
	defer m.Unlock()
	for _, j := range m.sessions{
		if j.timeAcceced.Unix() + m.maxLifeTime < time.Now().Unix(){
			m.SessionDestroy(j.sid)
		}
	}
	time.AfterFunc(time.Duration(maxLifeTime), func() {
		m.SessionGC(maxLifeTime)
	})
}

func sessionId() string{
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}