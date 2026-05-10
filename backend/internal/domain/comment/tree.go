package comment

import "sort"

func BuildTree(comments []Comment) []CommentNode {
	sorted := append([]Comment(nil), comments...)
	sort.SliceStable(sorted, func(i int, j int) bool {
		if sorted[i].IsPinned != sorted[j].IsPinned {
			return sorted[i].IsPinned
		}
		return sorted[i].CreatedAt.Before(sorted[j].CreatedAt)
	})

	nodes := make(map[string]*CommentNode, len(sorted))
	for _, comment := range sorted {
		copy := comment
		nodes[comment.ID] = &CommentNode{Comment: copy, Children: []CommentNode{}}
	}

	childrenByParent := make(map[string][]string)
	rootIDs := make([]string, 0)
	for _, comment := range sorted {
		if comment.ParentID != nil {
			if nodes[*comment.ParentID] != nil {
				childrenByParent[*comment.ParentID] = append(childrenByParent[*comment.ParentID], comment.ID)
				continue
			}
		}
		rootIDs = append(rootIDs, comment.ID)
	}

	var materialize func(string) CommentNode
	materialize = func(id string) CommentNode {
		node := *nodes[id]
		for _, childID := range childrenByParent[id] {
			node.Children = append(node.Children, materialize(childID))
		}
		return node
	}

	roots := make([]CommentNode, 0, len(rootIDs))
	for _, id := range rootIDs {
		roots = append(roots, materialize(id))
	}
	return roots
}
