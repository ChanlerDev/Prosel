'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';

import { postsApi } from '@/lib/api/posts';
import type { AdminPostListParams, PostEditorValues } from '@/types/post';

export const postKeys = {
  list: (params: Record<string, unknown>) => ['posts', params] as const,
  detail: (slug: string) => ['posts', slug] as const,
  adminList: (params: AdminPostListParams) => ['admin', 'posts', params] as const,
  adminDetail: (id: string) => ['admin', 'posts', id] as const,
};

export function useAdminPosts(params: AdminPostListParams) {
  return useQuery({
    queryKey: postKeys.adminList(params),
    queryFn: () => postsApi.admin.list(params),
  });
}

export function useAdminPost(id: string) {
  return useQuery({
    queryKey: postKeys.adminDetail(id),
    queryFn: () => postsApi.admin.detail(id),
    enabled: Boolean(id),
  });
}

export function useCreatePost() {
  const router = useRouter();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: PostEditorValues) => postsApi.admin.create(body),
    onSuccess: async (post) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'posts'] });
      router.replace(`/admin/posts/${post.id}/edit`);
    },
  });
}

export function useUpdatePost(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: PostEditorValues) => postsApi.admin.update(id, body),
    onSuccess: async (post) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'posts'] });
      await queryClient.invalidateQueries({ queryKey: postKeys.adminDetail(post.id) });
      await queryClient.invalidateQueries({ queryKey: ['posts'] });
    },
  });
}

export function useDeletePost() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => postsApi.admin.delete(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'posts'] });
      await queryClient.invalidateQueries({ queryKey: ['posts'] });
    },
  });
}

export function usePublishPost() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => postsApi.admin.publish(id),
    onSuccess: async (post) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'posts'] });
      await queryClient.invalidateQueries({ queryKey: postKeys.adminDetail(post.id) });
      await queryClient.invalidateQueries({ queryKey: ['posts'] });
    },
  });
}

export function useUnpublishPost() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => postsApi.admin.unpublish(id),
    onSuccess: async (post) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'posts'] });
      await queryClient.invalidateQueries({ queryKey: postKeys.adminDetail(post.id) });
      await queryClient.invalidateQueries({ queryKey: ['posts'] });
    },
  });
}
