package persistence

import (
	"sync"
	"time"
	"errors"
	"encoding/json"
	"io"
	"encoding/base64"
	"crypto/rand"
	"net/url"
	"net/http"
)

type Manager struct{
	sync.Mutex
	cookiesName string
	maxLifeTime int64
	Driver Driver
}

type Storage struct {

}

type Driver interface{
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionUpdate(sid string) error
	SessionDestroy(sid string) error
}

type SessionProvider interface{
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	SessionID() string
}

type Session struct{
	sid string
	timeAcceced time.Time
	values map[string]interface{}
}

func (s *Session) Values() map[string] interface{}{
	return s.values
}

func (s *Session) Set(key string, value interface{}) error{
	s.values[key] = value
	err := save(s)
	return err
}

func (s *Session) Get(key string) (interface{}, error){
	res, ok := s.values[key]
	if ok{
		return res, nil
	}
	return nil, errors.New("There is no value in values map with key: " + key)
}

func (s *Session) Delete(key string) error {
	delete(s.values, key)
	err := save(s)
	return err
}

func (s *Session) SessionID() string{
	return s.sid
}

func (s *Storage) SessionInit(sid string) (Session, error) {
	req, err := getRequestByName("initSession")
	if err != nil {
		return Session{}, err
	}
	sess := Session{sid, time.Now(), make(map[string]interface{}, 0)}
	js, err := json.Marshal(sess.values)
	if err != nil{
		return Session{}, err
	}
	_, err = db.Exec(req, sess.sid, sess.timeAcceced, string(js))
	if err != nil {
		return Session{}, err
	}
	return sess, nil
}

func (s *Storage)SessionRead(sid string) (Session, error) {
	rows := makeUserQuery([]string{"readSession", sid})
	var ssid string
	var timeAcceced time.Time
	var values []byte
	for rows.Next() {
		err := rows.Scan(&ssid, &timeAcceced, &values)
		if err != nil {
			return Session{}, err
		}
	}
	valuesMap, err := Unmarshal(values)
	if err != nil {
		return Session{}, err
	}
	return Session{ssid, timeAcceced, *valuesMap}, nil
}

func (s *Storage) SessionUpdate(sid string) error{
	sess, err := s.SessionRead(sid)
	if err != nil{
		return err
	}
	sess.timeAcceced = time.Now()
	err = save(&sess)
	return err
}

func (s *Storage)SessionDestroy(sid string) error{
	req, err := getRequestByName("deleteSessions")
	if err != nil{
		return err
	}
	_, err = db.Exec(req, sid)
	if err != nil{
		return err
	}
	return nil
}

func save(s *Session) error{
	req, err := getRequestByName("updateSession")
	if err != nil{
		return err
	}
	jsn, err := json.Marshal(&s.values)
	if err != nil{
		return err
	}
	_, err = db.Exec(req, s.sid, s.timeAcceced, jsn)
	if err != nil{
		return err
	}
	return nil
}

func Unmarshal(values []byte) (*map[string]interface{}, error) {
	valuesMap := make(map[string]interface{})
	err := json.Unmarshal(values, &valuesMap)
	return &valuesMap, err
}

func NewManager(cookiesName string, maxLifeTime int64) *Manager {
	return &Manager{cookiesName:cookiesName, maxLifeTime:maxLifeTime, Driver: new(Storage)}
}

func (m *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (Session, error) {
	m.Lock()
	defer m.Unlock()
	cookie, err := r.Cookie(m.cookiesName)
	if err != nil || cookie.Value == ""{
		sid := sessionId()
		session, err := m.Driver.SessionInit(sid)
		cookie := http.Cookie{
			Name: m.cookiesName,
			Value: url.QueryEscape(sid),
			HttpOnly: true,
			Path: "/",
			MaxAge: int(m.maxLifeTime),
			Expires:time.Now().Add(time.Duration(m.maxLifeTime))}
		http.SetCookie(w, &cookie)
		return session, err
	}else{
		sid, err := url.QueryUnescape(cookie.Value)
		if err != nil{
			return Session{}, err
		}
		session, err := m.Driver.SessionRead(sid)
		return session, err
	}
}

func (m *Manager) SessionGC() {
	m.Lock()
	defer m.Unlock()
	sessions, err := GetAllSessions()
	if err != nil{
		panic(err)
	}
	for _, j := range sessions{
		if j.timeAcceced.Unix() + m.maxLifeTime < time.Now().Unix(){
			m.Driver.SessionDestroy(j.sid)
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