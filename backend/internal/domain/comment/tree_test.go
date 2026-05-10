package comment

import (
	"testing"
	"time"
)

func TestBuildTreeKeepsGrandchildren(t *testing.T) {
	now := time.Now().UTC()
	root := Comment{ID: "root", RefType: RefPost, RefID: "post-1", Status: CommentApproved, CreatedAt: now}
	reply := Comment{ID: "reply", RefType: RefPost, RefID: "post-1", ParentID: &root.ID, RootID: &root.ID, Status: CommentApproved, CreatedAt: now.Add(time.Second)}
	grandchild := Comment{ID: "grandchild", RefType: RefPost, RefID: "post-1", ParentID: &reply.ID, RootID: &root.ID, Status: CommentApproved, CreatedAt: now.Add(2 * time.Second)}

	nodes := BuildTree([]Comment{root, reply, grandchild})
	if len(nodes) != 1 || len(nodes[0].Children) != 1 || len(nodes[0].Children[0].Children) != 1 {
		t.Fatalf("tree = %#v, want root -> reply -> grandchild", nodes)
	}
	if nodes[0].Children[0].Children[0].ID != "grandchild" {
		t.Fatalf("grandchild = %#v", nodes[0].Children[0].Children[0])
	}
}
