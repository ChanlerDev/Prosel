import { env } from '@/lib/env';
import type {
  ApiResponse,
  ChangePasswordRequest,
  HealthStatus,
  LoginRequest,
  PublicSettings,
  TokenResponse,
  UpdateProfileRequest,
  User,
} from '@/types/api';

export class ApiClientError extends Error {
  constructor(
    message: string,
    public readonly code: string,
    public readonly status: number,
  ) {
    super(message);
  }
}

async function request<T>(path: string, init?: RequestInit & { accessToken?: string }): Promise<T> {
  const headers = new Headers(init?.headers);
  headers.set('Content-Type', 'application/json');
  if (init?.accessToken) {
    headers.set('Authorization', `Bearer ${init.accessToken}`);
  }

  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    ...init,
    headers,
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

function post<T>(path: string, body: unknown, accessToken?: string) {
  return request<T>(path, { method: 'POST', body: JSON.stringify(body), cache: 'no-store', accessToken });
}

function patch<T>(path: string, body: unknown, accessToken: string) {
  return request<T>(path, { method: 'PATCH', body: JSON.stringify(body), cache: 'no-store', accessToken });
}

export const api = {
  system: {
    health: () => request<HealthStatus>('/health', { cache: 'no-store' }),
  },
  settings: {
    public: () => request<PublicSettings>('/settings/public'),
  },
  auth: {
    login: (body: LoginRequest) => post<TokenResponse>('/auth/login', body),
    refresh: (refreshToken: string) => post<TokenResponse>('/auth/refresh', { refreshToken }),
    logout: (refreshToken: string, accessToken: string) => post<{ ok: boolean }>('/auth/logout', { refreshToken }, accessToken),
    me: (accessToken: string) => request<User>('/auth/me', { cache: 'no-store', accessToken }),
    updateProfile: (body: UpdateProfileRequest, accessToken: string) => patch<User>('/admin/profile', body, accessToken),
    changePassword: (body: ChangePasswordRequest, accessToken: string) => patch<{ ok: boolean }>('/admin/password', body, accessToken),
  },
};
