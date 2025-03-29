package chatstate

import "sync"

// chatState управляет состояниями пользователей в чатах.
type ChatState struct {
	states map[int64][]string // Карта, где ключ — chatID, а значение — стек состояний.
	mu     sync.Mutex         // Мьютекс для потокобезопасного доступа.
}

var (
	chatStateInstance *ChatState
	once              sync.Once
)


// GetInstance возвращает единственный экземпляр chatState (реализация Singleton).
func GetInstance() *ChatState {
	once.Do(func() {
		chatStateInstance = &ChatState{
			states: make(map[int64][]string),
		}
	})
	return chatStateInstance
}

// PopState удаляет последнее состояние из стека и возвращает предыдущее состояние.
// Если стек пуст (или содержит только одно состояние), возвращает "start".
func (b *ChatState) PopState(chatID int64) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.states[chatID]) > 1 { // TODO: check 0 later!
		b.states[chatID] = b.states[chatID][:len(b.states[chatID])-1] // remove last element
		return b.states[chatID][len(b.states[chatID])-1]
	}
	return "start"
}

// PushState добавляет новое состояние в стек состояний пользователя.
func (b *ChatState) PushState(chatID int64, state string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.states[chatID] = append(b.states[chatID], state)
}

// InitState инициализирует стек состояний пользователя с состоянием state.
func (b *ChatState) InitState(chatID int64, state string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.states[chatID] = []string{state}
}
