import { EmptyState } from '@/components/features/system/states';
import { PostCard } from '@/components/features/post/post-card';
import type { Post } from '@/types/post';

export function PostList({ posts }: { posts: Post[] }) {
  if (posts.length === 0) {
    return <EmptyState title="No posts yet" description="Published posts will appear here." />;
  }
  return (
    <div className="grid gap-5">
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
}
