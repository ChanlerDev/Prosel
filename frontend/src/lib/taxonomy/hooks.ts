'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { taxonomyApi } from '@/lib/api/taxonomy';
import type { CategoryValues, TagValues, TopicValues } from '@/types/taxonomy';

export const taxonomyKeys = {
  categories: ['categories'] as const,
  tags: ['tags'] as const,
  topics: ['topics'] as const,
};

export function useCategories() {
  return useQuery({ queryKey: taxonomyKeys.categories, queryFn: taxonomyApi.categories });
}

export function useTags() {
  return useQuery({ queryKey: taxonomyKeys.tags, queryFn: taxonomyApi.tags });
}

export function useTopics() {
  return useQuery({ queryKey: taxonomyKeys.topics, queryFn: taxonomyApi.topics });
}

export function useCreateCategory() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (body: CategoryValues) => taxonomyApi.admin.createCategory(body), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.categories }) });
}

export function useUpdateCategory() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: ({ id, body }: { id: string; body: CategoryValues }) => taxonomyApi.admin.updateCategory(id, body), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.categories }) });
}

export function useDeleteCategory() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (id: string) => taxonomyApi.admin.deleteCategory(id), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.categories }) });
}

export function useCreateTag() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (body: TagValues) => taxonomyApi.admin.createTag(body), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.tags }) });
}

export function useUpdateTag() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: ({ id, body }: { id: string; body: TagValues }) => taxonomyApi.admin.updateTag(id, body), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.tags }) });
}

export function useDeleteTag() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (id: string) => taxonomyApi.admin.deleteTag(id), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.tags }) });
}

export function useCreateTopic() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (body: TopicValues) => taxonomyApi.admin.createTopic(body), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.topics }) });
}

export function useUpdateTopic() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: ({ id, body }: { id: string; body: TopicValues }) => taxonomyApi.admin.updateTopic(id, body), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.topics }) });
}

export function useDeleteTopic() {
  const queryClient = useQueryClient();
  return useMutation({ mutationFn: (id: string) => taxonomyApi.admin.deleteTopic(id), onSuccess: () => queryClient.invalidateQueries({ queryKey: taxonomyKeys.topics }) });
}
