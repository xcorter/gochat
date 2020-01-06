package chat

import (
	"net"
)

type Connection struct {
	conn     net.Conn
	threadId int
}

type ConnectionMap struct {
	connections []Connection
}

func newConnection(conn net.Conn, threadId int) Connection {
	connection := Connection{
		conn:     conn,
		threadId: threadId,
	}
	return connection
}

func (c *ConnectionMap) NewConnection(conn net.Conn, threadId int) {
	connection := newConnection(conn, threadId)
	c.connections = append(c.connections, connection)
}

func (c *ConnectionMap) GetConnectionsByThreadId(threadId int) []Connection {
	var result []Connection
	for _, connection := range c.connections {
		if connection.threadId == threadId {
			result = append(result, connection)
		}
	}
	return result
}

func NewConnectionMap() ConnectionMap {
	cm := ConnectionMap{}
	return cm
}
