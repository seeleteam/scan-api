package node

import (
	"github.com/seeleteam/scan-api/database"
	"testing"
	"fmt"
	"sync"
)

func newTestNodeService()  *NodeService{
	return &NodeService{
		nodeMap: make(map[string]database.DBNodeInfo),
	}
}

func Test_writeAndRead(t *testing.T) {
	service := newTestNodeService()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 100000; i++ {
			service.nodeMapLock.Lock()
			service.nodeMap[string(i)] = database.DBNodeInfo{
				LastSeen:int64(i),
			}
			service.nodeMapLock.Unlock()
		}
		wg.Done()

	}()

	func() {
		for i := 0; i < 100000; i++ {
			service.nodeMapLock.RLock()
			_, ok:= service.nodeMap[string(i)]
			if !ok{
				fmt.Println("read map error:", i)
			}
			service.nodeMapLock.RUnlock()
		}
		wg.Done()
	}()
	wg.Wait()

}
