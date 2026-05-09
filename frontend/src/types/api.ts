export interface ApiResponse<T> {
  data?: T;
  meta?: PaginationMeta;
  error?: ApiError;
}

export interface PaginationMeta {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
}

export interface ApiError {
  code: string;
  message: string;
  details?: unknown;
}

export interface SiteSetting {
  key: string;
  value: string | number | boolean | null;
  valueType: 'string' | 'number' | 'boolean' | 'json';
}

export type PublicSettings = {
  site_name?: string;
  site_description?: string;
  site_url?: string;
  posts_per_page?: number;
  comment_moderation?: boolean;
  analytics_enabled?: boolean;
};

export interface HealthStatus {
  status: 'healthy' | 'unhealthy';
  databaseOk: boolean;
  redisOk: boolean;
  version: string;
  checkedAt: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  displayName: string;
  avatarUrl?: string;
  bio?: string;
  role: 'admin';
}

export interface TokenResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user?: User;
}

export interface LoginRequest {
  login: string;
  password: string;
}

export interface UpdateProfileRequest {
  displayName: string;
  avatarUrl: string;
  bio: string;
}

export interface ChangePasswordRequest {
  oldPassword: string;
  newPassword: string;
}
