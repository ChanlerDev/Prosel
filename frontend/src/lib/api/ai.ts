import { request } from '@/lib/api/client';
import { readAuthState } from '@/lib/auth/store';
import type { AIStatus, AISummary, AITranslation, GenerateSummaryValues, GenerateTranslationValues } from '@/types/ai';

function accessToken() {
  const auth = readAuthState();
  if (!auth.accessToken) {
    throw new Error('Authentication required');
  }
  return auth.accessToken;
}

function queryString(params: Record<string, string>) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value) search.set(key, value);
  }
  const value = search.toString();
  return value ? `?${value}` : '';
}

export const aiApi = {
  summary: (refType: string, refId: string, lang = 'zh') => request<AISummary>(`/ai/summaries${queryString({ refType, refId, lang })}`, { cache: 'no-store' }),
  translation: (refType: string, refId: string, lang: string) => request<AITranslation>(`/ai/translations${queryString({ refType, refId, lang })}`, { cache: 'no-store' }),
  admin: {
    status: () => request<AIStatus>('/admin/ai/status', { cache: 'no-store', accessToken: accessToken() }),
    generateSummary: (body: GenerateSummaryValues) => request<AISummary>('/admin/ai/summaries', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
    generateTranslation: (body: GenerateTranslationValues) => request<AITranslation>('/admin/ai/translations', { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken: accessToken() }),
  },
};
