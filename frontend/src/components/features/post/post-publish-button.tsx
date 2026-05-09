'use client';

import { Button } from '@/components/ui/button';
import { usePublishPost, useUnpublishPost } from '@/lib/posts/hooks';
import type { Post } from '@/types/post';

export function PostPublishButton({ post }: { post: Post }) {
  const publish = usePublishPost();
  const unpublish = useUnpublishPost();
  const isPublished = post.status === 'published';
  const isPending = publish.isPending || unpublish.isPending;

  return (
    <Button
      className="px-3 py-1 text-xs"
      disabled={isPending}
      onClick={() => (isPublished ? unpublish.mutate(post.id) : publish.mutate(post.id))}
      type="button"
    >
      {isPublished ? 'Unpublish' : 'Publish'}
    </Button>
  );
}
