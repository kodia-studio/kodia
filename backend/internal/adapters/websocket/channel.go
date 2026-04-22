package websocket

import (
	"strings"
	"sync"
)

// ChannelType classifies the access model of a channel.
type ChannelType string

const (
	// ChannelTypePublic is accessible by any client.
	ChannelTypePublic ChannelType = "public"
	// ChannelTypePrivate requires the client to be authenticated.
	ChannelTypePrivate ChannelType = "private"
	// ChannelTypePresence is a private channel that also tracks member presence.
	ChannelTypePresence ChannelType = "presence"
)

// ChannelMiddlewareFunc is a function that intercepts messages before delivery.
// Return false to drop the message.
type ChannelMiddlewareFunc func(conn *Connection, msg *Message) bool

// Channel represents a named, typed broadcast channel within the Hub.
type Channel struct {
	Name        string
	Type        ChannelType
	middleware  []ChannelMiddlewareFunc
	subscribers map[*Connection]bool
	mu          sync.RWMutex
}

// NewChannel creates a new Channel.
func NewChannel(name string, channelType ChannelType) *Channel {
	return &Channel{
		Name:        name,
		Type:        channelType,
		subscribers: make(map[*Connection]bool),
	}
}

// Use adds a middleware function to the channel's processing pipeline.
func (c *Channel) Use(fn ChannelMiddlewareFunc) *Channel {
	c.middleware = append(c.middleware, fn)
	return c
}

// Subscribe adds a connection to this channel.
func (c *Channel) Subscribe(conn *Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subscribers[conn] = true
}

// Unsubscribe removes a connection from this channel.
func (c *Channel) Unsubscribe(conn *Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.subscribers, conn)
}

// Broadcast delivers a message to all subscribers, running middleware first.
func (c *Channel) Broadcast(msg *Message) {
	c.mu.RLock()
	subs := make([]*Connection, 0, len(c.subscribers))
	for conn := range c.subscribers {
		subs = append(subs, conn)
	}
	c.mu.RUnlock()

	for _, conn := range subs {
		// Run middleware chain
		allowed := true
		for _, mw := range c.middleware {
			if !mw(conn, msg) {
				allowed = false
				break
			}
		}
		if !allowed {
			continue
		}
		select {
		case conn.send <- msg:
		default:
			// Drop if buffer is full
		}
	}
}

// SubscriberCount returns the number of active subscribers.
func (c *Channel) SubscriberCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.subscribers)
}

// GetPresenceMembers returns user IDs of all connected members (presence channels only).
func (c *Channel) GetPresenceMembers() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	seen := make(map[string]bool)
	members := make([]string, 0)
	for conn := range c.subscribers {
		if conn.UserID != "" && !seen[conn.UserID] {
			seen[conn.UserID] = true
			members = append(members, conn.UserID)
		}
	}
	return members
}

// ParseChannelType infers a ChannelType from a channel name prefix.
// "private-*" → Private, "presence-*" → Presence, others → Public.
func ParseChannelType(name string) ChannelType {
	switch {
	case strings.HasPrefix(name, "private-"):
		return ChannelTypePrivate
	case strings.HasPrefix(name, "presence-"):
		return ChannelTypePresence
	default:
		return ChannelTypePublic
	}
}
