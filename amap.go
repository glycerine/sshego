package sshego

import "sync"

// atomic map from string to *User

//go:generate msgp

type AtomicUserMap struct {
	U   map[string]*User
	tex sync.RWMutex
}

func NewAtomicUserMap() *AtomicUserMap {
	return &AtomicUserMap{
		U: make(map[string]*User),
	}
}

func (m *AtomicUserMap) Get(key string) *User {
	m.tex.RLock()
	defer m.tex.RUnlock()
	return m.U[key]
}

func (m *AtomicUserMap) Get2(key string) (*User, bool) {
	m.tex.RLock()
	defer m.tex.RUnlock()
	v, ok := m.U[key]
	return v, ok
}

func (m *AtomicUserMap) Set(key string, val *User) {
	m.tex.Lock()
	defer m.tex.Unlock()
	m.U[key] = val
}

func (m *AtomicUserMap) Del(key string) {
	m.tex.Lock()
	defer m.tex.Unlock()
	delete(m.U, key)
}
