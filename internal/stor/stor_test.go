package stor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Path string
}

func TestIndex(t *testing.T) {
	col := NewCollection()

	node1, _ := col.NewNode(uuid.New().String(), Data{Path: "/node1"}, true)
	node2, _ := col.NewNode(uuid.New().String(), Data{Path: "/node2"}, false)
	node3, _ := col.NewNode(uuid.New().String(), Data{Path: "/node3"}, false)
	node4, _ := col.NewNode(uuid.New().String(), Data{Path: "/node4"}, true)
	node5, _ := col.NewNode(uuid.New().String(), Data{Path: "/node5"}, false)
	node6, _ := col.NewNode(uuid.New().String(), Data{Path: "/node6"}, false)

	assert.Equal(t, node1.Id, node2.Prev.Id)
	assert.Equal(t, node1.Next.Id, node2.Id)
	assert.Equal(t, node2.Id, node3.Prev.Id)
	assert.Equal(t, node2.Next.Id, node3.Id)

	assert.Equal(t, node4.Id, node5.Prev.Id)
	assert.Equal(t, node4.Next.Id, node5.Id)
	assert.Equal(t, node5.Id, node6.Prev.Id)
	assert.Equal(t, node5.Next.Id, node6.Id)

	n := col.GetNode(node5.Id)
	assert.Equal(t, node5.Id, n.Id)

	n = col.GetStartNode(node5.Id)
	assert.Equal(t, node4.Id, n.Id)

	ns := col.GetBeforeNodes(node6.Id)
	assert.Equal(t, []*Node{node4, node5, node6}, ns)
	ns = col.GetBeforeNodes(node5.Id)
	assert.Equal(t, []*Node{node4, node5}, ns)

	err := Serialize(col, "./index")
	assert.Nil(t, err)

	col1 := Collection{}
	err = Deserialize(&col1, "./index")
	assert.Nil(t, err)
}
