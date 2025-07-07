package documents

import (
	"sync"
	"time"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Manager handles document lifecycle and content management
type Manager struct {
	documents      map[protocol.DocumentUri]string
	lastUpdate     map[protocol.DocumentUri]time.Time
	mutex          sync.RWMutex
}

// NewManager creates a new document manager
func NewManager() *Manager {
	return &Manager{
		documents:  make(map[protocol.DocumentUri]string),
		lastUpdate: make(map[protocol.DocumentUri]time.Time),
	}
}

// DidOpen handles document open events
func (m *Manager) DidOpen(uri protocol.DocumentUri, content string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.documents[uri] = content
	m.lastUpdate[uri] = time.Now()
}

// DidChange handles document change events
func (m *Manager) DidChange(uri protocol.DocumentUri, changes []protocol.TextDocumentContentChangeEvent) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	content, exists := m.documents[uri]
	if !exists {
		return ""
	}
	
	// Apply changes (for incremental changes, we'd need more sophisticated logic)
	// For now, assume full document replacement
	for _, change := range changes {
		if change.Range == nil {
			// Full document change
			content = change.Text
		} else {
			// Range-based change (more complex implementation needed)
			// For simplicity, replace full content
			content = change.Text
		}
	}
	
	m.documents[uri] = content
	m.lastUpdate[uri] = time.Now()
	
	return content
}

// UpdateContent updates document content directly
func (m *Manager) UpdateContent(uri protocol.DocumentUri, content string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.documents[uri] = content
	m.lastUpdate[uri] = time.Now()
}

// GetContent retrieves document content
func (m *Manager) GetContent(uri protocol.DocumentUri) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	content, exists := m.documents[uri]
	return content, exists
}

// GetLastUpdate returns the last update time for a document
func (m *Manager) GetLastUpdate(uri protocol.DocumentUri) (time.Time, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	lastUpdate, exists := m.lastUpdate[uri]
	return lastUpdate, exists
}

// DidClose handles document close events
func (m *Manager) DidClose(uri protocol.DocumentUri) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	delete(m.documents, uri)
	delete(m.lastUpdate, uri)
}

// GetAllDocuments returns all managed documents
func (m *Manager) GetAllDocuments() map[protocol.DocumentUri]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	result := make(map[protocol.DocumentUri]string)
	for uri, content := range m.documents {
		result[uri] = content
	}
	
	return result
}