import { env } from '@/lib/env';
import type { ApiResponse, HealthStatus, PublicSettings } from '@/types/api';

export class ApiClientError extends Error {
  constructor(
    message: string,
    public readonly code: string,
    public readonly status: number,
  ) {
    super(message);
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
    next: init?.cache === 'no-store' ? undefined : { revalidate: 60 },
  });

  const body = (await response.json()) as ApiResponse<T>;
  if (!response.ok || body.error) {
    throw new ApiClientError(
      body.error?.message ?? 'Request failed',
      body.error?.code ?? 'REQUEST_FAILED',
      response.status,
    );
  }

  return body.data as T;
}

export const api = {
  system: {
    health: () => request<HealthStatus>('/health', { cache: 'no-store' }),
  },
  settings: {
    public: () => request<PublicSettings>('/settings/public'),
  },
};
