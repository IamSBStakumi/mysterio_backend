package repository

import (
	"errors"
	"sync"

	"github.com/IamSBStakumi/mysterio_backend/internal/domain"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

// SessionRepository はセッションデータの永続化を担当する
type SessionRepository interface {
	Save(session *domain.Session) error
	FindByID(sessionID string) (*domain.Session, error)
	Update(session *domain.Session) error
	Delete(sessionID string) error
}

// InMemorySessionRepository はインメモリでセッションを管理する
type InMemorySessionRepository struct {
	sessions map[string]*domain.Session
	mu       sync.RWMutex
}

// NewInMemorySessionRepository は新しいインメモリリポジトリを作成する
func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: make(map[string]*domain.Session),
	}
}

// Save はセッションを保存する
func (r *InMemorySessionRepository) Save(session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	return nil
}

// FindByID はセッションIDからセッションを取得する
func (r *InMemorySessionRepository) FindByID(sessionID string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

// Update はセッション情報を更新する
func (r *InMemorySessionRepository) Update(session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[session.ID]; !exists {
		return ErrSessionNotFound
	}

	r.sessions[session.ID] = session
	return nil
}

// Delete はセッションを削除する
func (r *InMemorySessionRepository) Delete(sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[sessionID]; !exists {
		return ErrSessionNotFound
	}

	delete(r.sessions, sessionID)
	return nil
}
