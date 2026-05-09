'use client';

import { PostEditor } from '@/components/features/post/post-editor';
import { useCreatePost } from '@/lib/posts/hooks';

export function AdminPostCreate() {
  const create = useCreatePost();
  return <PostEditor error={create.isError ? create.error.message : undefined} isPending={create.isPending} onSubmit={(values) => create.mutate(values)} />;
}
