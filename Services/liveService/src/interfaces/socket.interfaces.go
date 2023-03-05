package interfaces

// socket interface
type Socket struct {
	SessionID string
	HashedURL string
	SocketURL string
}

// Message interface x
type Message struct {
	Type        string `json:"type"`
	UserID      string `json:"userID"`
	Description string `json:"description"`
	Candidate   string `json:"candidate"`
	To          string `json:"to"`
}
