package chatstate

func (b *ChatState) GetLastState(chatID int64) string {
	b.mu.Lock()
	defer b.mu.Unlock()

	if states, exists := b.states[chatID]; exists && len(states) > 0{
		return	states[len(states)-1]
	}
	return ""
}



func (b *ChatState) GetDefState(){

}