'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { searchApi } from '@/lib/api/search';
import type { SearchParams } from '@/types/search';

export const searchKeys = {
  query: (params: SearchParams) => ['search', params] as const,
  adminStatus: () => ['admin', 'search', 'status'] as const,
};

export function useSearch(params: SearchParams) {
  return useQuery({ queryKey: searchKeys.query(params), queryFn: () => searchApi.query(params), enabled: params.q.trim().length > 0 });
}

export function useSearchIndexStatus() {
  return useQuery({ queryKey: searchKeys.adminStatus(), queryFn: () => searchApi.admin.status() });
}

export function useRebuildSearchIndex() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: () => searchApi.admin.rebuild(), onSuccess: async () => queryClient.invalidateQueries({ queryKey: searchKeys.adminStatus() }) });
}
