'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { subscribeApi } from '@/lib/api/subscribe';
import type { SubscribeRequest, SubscriberListParams } from '@/types/subscribe';

export const subscribeKeys = {
  adminList: (params: SubscriberListParams) => ['admin', 'subscribers', params] as const,
};

export function useSubscribe() {
  return useMutation({
    mutationFn: (body: SubscribeRequest) => subscribeApi.create(body),
  });
}

export function useAdminSubscribers(params: SubscriberListParams) {
  return useQuery({
    queryKey: subscribeKeys.adminList(params),
    queryFn: () => subscribeApi.admin.list(params),
  });
}

export function useNotifyPostSubscribers() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (postId: string) => subscribeApi.admin.notifyPost(postId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['admin', 'subscribers'] });
    },
  });
}
