'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import { aiApi } from '@/lib/api/ai';
import type { GenerateSummaryValues, GenerateTranslationValues } from '@/types/ai';

export const aiKeys = {
  status: ['admin', 'ai', 'status'] as const,
  summary: (refType: string, refId: string, lang: string) => ['ai', 'summary', refType, refId, lang] as const,
  translation: (refType: string, refId: string, lang: string) => ['ai', 'translation', refType, refId, lang] as const,
};

export function useAIStatus() {
  return useQuery({ queryKey: aiKeys.status, queryFn: aiApi.admin.status });
}

export function usePublicAISummary(refType: string, refId: string, lang = 'zh') {
  return useQuery({ queryKey: aiKeys.summary(refType, refId, lang), queryFn: () => aiApi.summary(refType, refId, lang), retry: false, enabled: Boolean(refId) });
}

export function usePublicAITranslation(refType: string, refId: string, lang: string) {
  return useQuery({ queryKey: aiKeys.translation(refType, refId, lang), queryFn: () => aiApi.translation(refType, refId, lang), retry: false, enabled: Boolean(refId && lang) });
}

export function useGenerateSummary() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: GenerateSummaryValues) => aiApi.admin.generateSummary(body),
    onSuccess: async (summary) => {
      await queryClient.invalidateQueries({ queryKey: aiKeys.summary(summary.refType, summary.refId, summary.language) });
    },
  });
}

export function useGenerateTranslation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: GenerateTranslationValues) => aiApi.admin.generateTranslation(body),
    onSuccess: async (translation) => {
      await queryClient.invalidateQueries({ queryKey: aiKeys.translation(translation.refType, translation.refId, translation.targetLanguage) });
    },
  });
}
