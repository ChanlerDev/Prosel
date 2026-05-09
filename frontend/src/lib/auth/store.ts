import type { TokenResponse, User } from '@/types/api';

const storageKey = 'prosel.auth';

export interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
}

export function readAuthState(): AuthState {
  if (typeof window === 'undefined') {
    return emptyAuthState();
  }
  const raw = window.localStorage.getItem(storageKey);
  if (!raw) {
    return emptyAuthState();
  }
  try {
    const parsed = JSON.parse(raw) as AuthState;
    return {
      user: parsed.user ?? null,
      accessToken: parsed.accessToken ?? null,
      refreshToken: parsed.refreshToken ?? null,
    };
  } catch {
    clearAuthState();
    return emptyAuthState();
  }
}

export function saveAuthState(tokens: TokenResponse, fallbackUser?: User | null) {
  if (typeof window === 'undefined') {
    return;
  }
  const state: AuthState = {
    user: tokens.user ?? fallbackUser ?? null,
    accessToken: tokens.accessToken,
    refreshToken: tokens.refreshToken,
  };
  window.localStorage.setItem(storageKey, JSON.stringify(state));
}

export function clearAuthState() {
  if (typeof window === 'undefined') {
    return;
  }
  window.localStorage.removeItem(storageKey);
}

function emptyAuthState(): AuthState {
  return { user: null, accessToken: null, refreshToken: null };
}
