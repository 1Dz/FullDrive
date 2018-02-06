package persistence

import (
	"sync"
	"time"
	"io"
	"encoding/base64"
	"crypto/rand"
	"errors"
	"net/http"
	"net/url"
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

type Driver interface{
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

func (m *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (Session, error){
	m.Lock()
	defer m.Unlock()
	cookie, err := r.Cookie(m.cookiesName)
	if err != nil || cookie.Value == ""{
		sid := sessionId()
		session, err := m.SessionInit(sid)
		cookie := http.Cookie{
			Name: m.cookiesName,
			Value: url.QueryEscape(sid),
			HttpOnly: true,
			MaxAge: int(m.maxLifeTime)}
			http.SetCookie(w, &cookie)
			return session, err
	}else{
		sid, err := url.QueryUnescape(cookie.Value)
		if err != nil{
			return Session{}, err
		}
		session, err := m.SessionRead(sid)
		return session, err
	}
}

func (m *Manager) SessionInit(sid string) (Session, error){
	s := Session{sid, time.Now(), nil}
	err := SessionInit(&s)
	if err != nil{
		return Session{}, err
	}
	return s, nil
}

func (m *Manager) SessionRead(sid string) (Session, error){
	s, err := SessionRead(sid)
	if err != nil{
		return Session{}, err
	}
	m.SessionUpdate(s.sid)
	return s, nil
}

func (m *Manager) SessionDestroy(sid string) error{
	err := SessionDestroy(sid)
	return err
}

func (m *Manager) SessionUpdate(sid string) error{
	s, err := m.SessionRead(sid)
	if err != nil{
		return err
	}
	s.timeAcceced = time.Now()
	err = SessionUpdate(&s)
	return err
}

func (m *Manager) SessionGC() {
	m.Lock()
	defer m.Unlock()
	s, err := GetAllSessions()
	if err != nil{
		panic(err)
	}
	for _, j := range s{
		if j.timeAcceced.Unix() + m.maxLifeTime < time.Now().Unix(){
			m.SessionDestroy(j.sid)
		}
	}
	time.AfterFunc(time.Duration(m.maxLifeTime), func() {
		m.SessionGC()
	})
}

func sessionId() string{
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}