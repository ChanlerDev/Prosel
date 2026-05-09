package taxonomy

import "context"

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	ListTree(ctx context.Context) ([]CategoryNode, error)
	SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error)
}

type TagRepository interface {
	Create(ctx context.Context, tag *Tag) error
	Update(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Tag, error)
	List(ctx context.Context) ([]TagWithCount, error)
	SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error)
	ReplacePostTags(ctx context.Context, postID string, tagIDs []string) error
}

type TopicRepository interface {
	Create(ctx context.Context, topic *Topic) error
	Update(ctx context.Context, topic *Topic) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Topic, error)
	GetBySlug(ctx context.Context, slug string) (*Topic, error)
	List(ctx context.Context) ([]Topic, error)
	ReplaceItems(ctx context.Context, topicID string, items []TopicItem) error
	ListItems(ctx context.Context, topicID string) ([]TopicItem, error)
	SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error)
}
