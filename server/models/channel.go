package models

import (
	"time"
)

// Message represents a message in a channel
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Channel represents a communication channel
type Channel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewChannel creates a new channel
func NewChannel(name string) *Channel {
	return &Channel{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddMessage adds a new message to the channel
func (c *Channel) AddMessage(content string) {
	message := Message{
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	c.Messages = append(c.Messages, message)
	c.UpdatedAt = time.Now()
}

// UpdateMessage updates an existing message in the channel
func (c *Channel) UpdateMessage(id int, content string) {
	for i, msg := range c.Messages {
		if msg.ID == id {
			c.Messages[i].Content = content
			c.Messages[i].UpdatedAt = time.Now()
			c.UpdatedAt = time.Now()
			break
		}
	}
}

// DeleteMessage deletes a message from the channel
func (c *Channel) DeleteMessage(id int) {
	for i, msg := range c.Messages {
		if msg.ID == id {
			c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
			c.UpdatedAt = time.Now()
			break
		}
	}
}
