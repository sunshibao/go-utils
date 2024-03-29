package uuid

import (
	"fmt"
	"testing"
)

//******************************************************************************

func TestGetNodeID(t *testing.T) {
	nodeId, err := GetNodeID()
	if err != nil {
		t.Fatalf("error GetNodeID, %s", err)
	}

	fmt.Println(nodeId)
}

func TestGetUUID(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	id := node.Generate()

	t.Logf("Int64    : %#v", id)
}

func TestBatchGetUUID(t *testing.T) {

	node, _ := NewNode(1)

	go func() {
		for i := 0; i < 1000000000; i++ {

			NewNode(1)
		}
	}()
	var batchUUID []int64

	for i := 0; i < 4000; i++ {
		generate := node.Generate()
		batchUUID = append(batchUUID, generate)
	}
	fmt.Println(batchUUID)

}
