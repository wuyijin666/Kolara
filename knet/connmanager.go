package knet

import (
	"Kolara/kiface"
)

import (
	"sync"
	"errors"
	"fmt"
)

type ConnManager struct {
	connections map[uint32] kiface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager{
	return &ConnManager{
		connections: make(map[uint32] kiface.IConnection),
	}
}

func(cm *ConnManager) Add (conn kiface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn
	fmt.Printf("[DEBUG] Connection added, connId: %d, current connections: %d\n", conn.GetConnID(), len(cm.connections))

}

func(cm *ConnManager) Remove (conn kiface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID())
	fmt.Printf("[DEBUG] Connection removed, connId: %d, remaining connections: %d\n", conn.GetConnID(), len(cm.connections))
}



// 根据连接id， 获取连接
func (cm *ConnManager) Get(connId uint32) (kiface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	
	if conn, ok := cm.connections[connId]; ok {
		fmt.Printf("[DEBUG] Connection found, connId: %d\n", connId)
		return conn, nil
	}
	fmt.Printf("[DEBUG] Connection not found, connId: %d\n", connId)
	return nil, errors.New("connection not found")
}
// 获取当前连接总数
func (cm *ConnManager) Len() int {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	
	fmt.Printf("[DEBUG] Current connection count: %d\n", len(cm.connections))
	return len(cm.connections)
	
}
	// 清理所有连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for k := range cm.connections {
		delete(cm.connections, k)
	}
	fmt.Println("[DEBUG] All connections cleared")
}
