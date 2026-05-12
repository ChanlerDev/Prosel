'use client';

import { AdminAITools } from '@/components/features/ai/admin-ai-tools';
import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { PostEditor } from '@/components/features/post/post-editor';
import { PostPublishButton } from '@/components/features/post/post-publish-button';
import { PostStatusBadge } from '@/components/features/post/post-status-badge';
import { useAdminPost, useUpdatePost } from '@/lib/posts/hooks';

export function AdminPostEdit({ id }: { id: string }) {
  const post = useAdminPost(id);
  const update = useUpdatePost(id);

  if (post.isLoading) return <LoadingState />;
  if (post.isError || !post.data) return <ApiErrorState message="Unable to load post." />;

  return (
    <div className="grid gap-5">
      <div className="flex items-center gap-3">
        <PostStatusBadge status={post.data.status} />
        <PostPublishButton post={post.data} refId={post.data.id} />
      </div>
      <PostEditor error={update.isError ? update.error.message : undefined} isPending={update.isPending} onSubmit={(values) => update.mutate(values)} post={post.data} refId={post.data.id} />
      <AdminAITools post={post.data} />
    </div>
  );
}
