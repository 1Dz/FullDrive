package handlers

import (
	"sync"
	"io"
	"encoding/base64"
	"crypto/rand"
	"net/http"
	"net/url"
	"time"
	"fmt"
)

var provides = make(map[string]Provider)

type Manager struct{
	cookieName string
	lock sync.Mutex
	provider Provider
	maxLifeTime int64
}

type Provider interface{
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type ProviderControl struct{
	lock sync.Mutex
}

type SessionControl struct{
	sid string
	timeAcceced time.Time
	value map[interface{}]interface{}
}

type Session interface{
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxlifetime}, nil
}

func Register(name string, provider Provider){
	if provider == nil{
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice with " + name)
	}
	provides[name] = provider
}

func (m *Manager) sessionId() string{
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (m *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == ""{
		sid := m.sessionId()
		session, _ = m.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name: m.cookieName,
			Value:url.QueryEscape(sid),
			HttpOnly:true,
			MaxAge:int(m.maxLifeTime)}
		http.SetCookie(w, &cookie)
	}else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = m.provider.SessionRead(sid)
	}
	return
}

func (m *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == ""{
		return
	}else{
		m.lock.Lock()
		defer m.lock.Unlock()
		m.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{
			Name: m.cookieName,
			HttpOnly:true,
			Expires:expiration,
			MaxAge: - 1}
			http.SetCookie(w, &cookie)
	}
}

func (m *Manager) GC() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.provider.SessionGC(m.maxLifeTime)
	time.AfterFunc(time.Duration(m.maxLifeTime), func() {
		m.GC()
	})
}


func (p *ProviderControl) SessionInit(sid string)(Session, error){
	p.lock.Lock()
	defer p.lock.Unlock()
	sess := &SessionControl{sid, time.Now(), nil}

}