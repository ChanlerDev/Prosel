export type AIRefType = 'post' | 'note' | 'page';

export interface AISummary {
  id: string;
  refType: AIRefType;
  refId: string;
  language: string;
  contentHash?: string;
  summary: string;
  keywords: string[];
  provider?: string;
  model?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AITranslation {
  id: string;
  refType: AIRefType;
  refId: string;
  sourceLanguage: string;
  targetLanguage: string;
  contentHash?: string;
  title?: string;
  summary?: string;
  contentMarkdown: string;
  provider?: string;
  model?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AIStatus {
  configured: boolean;
}

export interface GenerateSummaryValues {
  refType: AIRefType;
  refId: string;
  language: string;
}

export interface GenerateTranslationValues {
  refType: AIRefType;
  refId: string;
  sourceLanguage: string;
  targetLanguage: string;
}
