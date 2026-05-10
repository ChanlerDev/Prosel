'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { commentsApi } from '@/lib/api/comments';
import type { AdminCommentListParams, CommentRefType, CommentStatus, SubmitCommentValues } from '@/types/comment';

export const commentKeys = {
  publicList: (refType: CommentRefType, refId: string) => ['comments', refType, refId] as const,
  adminList: (params: AdminCommentListParams) => ['admin', 'comments', params] as const,
};

export function useComments(refType: CommentRefType, refId: string) {
  return useQuery({
    queryKey: commentKeys.publicList(refType, refId),
    queryFn: () => commentsApi.list(refType, refId),
    enabled: Boolean(refId),
  });
}

export function useSubmitComment(refType: CommentRefType, refId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: SubmitCommentValues) => commentsApi.submit(body),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: commentKeys.publicList(refType, refId) });
      await queryClient.invalidateQueries({ queryKey: ['admin', 'comments'] });
    },
  });
}

export function useAdminComments(params: AdminCommentListParams) {
  return useQuery({
    queryKey: commentKeys.adminList(params),
    queryFn: () => commentsApi.admin.list(params),
  });
}

export function useModerateComment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: CommentStatus }) => commentsApi.admin.moderate(id, status),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'comments'] });
      await queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
  });
}

export function useAdminReply() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, content }: { id: string; content: string }) => commentsApi.admin.reply(id, content),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'comments'] });
      await queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
  });
}

export function useDeleteComment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => commentsApi.admin.delete(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'comments'] });
      await queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
  });
}
