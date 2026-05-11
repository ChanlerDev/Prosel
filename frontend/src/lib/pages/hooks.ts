'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';

import { pagesApi } from '@/lib/api/pages';
import type { AdminPageListParams, FriendValues, PageEditorValues } from '@/types/page';

export const pageKeys = {
  adminList: (params: AdminPageListParams) => ['admin', 'pages', params] as const,
  adminDetail: (id: string) => ['admin', 'pages', id] as const,
  adminFriends: (status: string) => ['admin', 'friends', status] as const,
};

export function useAdminPages(params: AdminPageListParams) {
  return useQuery({ queryKey: pageKeys.adminList(params), queryFn: () => pagesApi.admin.list(params) });
}

export function useAdminPage(id: string) {
  return useQuery({ queryKey: pageKeys.adminDetail(id), queryFn: () => pagesApi.admin.detail(id), enabled: Boolean(id) });
}

export function useCreatePage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: PageEditorValues) => pagesApi.admin.create(body),
    onSuccess: async (page) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'pages'] });
      router.replace(`/admin/pages/${page.id}/edit`);
    },
  });
}

export function useUpdatePage(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: PageEditorValues) => pagesApi.admin.update(id, body),
    onSuccess: async (page) => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'pages'] });
      await queryClient.invalidateQueries({ queryKey: pageKeys.adminDetail(page.id) });
    },
  });
}

export function useDeletePage() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (id: string) => pagesApi.admin.delete(id), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['admin', 'pages'] }) });
}

export function useAdminFriends(status = '') {
  return useQuery({ queryKey: pageKeys.adminFriends(status), queryFn: () => pagesApi.admin.friends(status) });
}

export function useCreateFriend() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (body: FriendValues) => pagesApi.admin.createFriend(body), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['admin', 'friends'] }) });
}

export function useUpdateFriend() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: ({ id, body }: { id: string; body: FriendValues }) => pagesApi.admin.updateFriend(id, body), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['admin', 'friends'] }) });
}

export function useDeleteFriend() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (id: string) => pagesApi.admin.deleteFriend(id), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['admin', 'friends'] }) });
}
