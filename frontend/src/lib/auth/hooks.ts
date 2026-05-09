'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';

import { api } from '@/lib/api/client';
import { clearAuthState, readAuthState, saveAuthState } from '@/lib/auth/store';
import type { ChangePasswordRequest, LoginRequest, UpdateProfileRequest } from '@/types/api';

export function useAuthState() {
  return readAuthState();
}

export function useMe() {
  const auth = readAuthState();
  return useQuery({
    queryKey: ['auth', 'me'],
    queryFn: async () => {
      if (!auth.accessToken) {
        throw new Error('Authentication required');
      }
      return api.auth.me(auth.accessToken);
    },
    retry: false,
    enabled: Boolean(auth.accessToken),
  });
}

export function useLogin() {
  const router = useRouter();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: LoginRequest) => api.auth.login(body),
    onSuccess: (tokens) => {
      saveAuthState(tokens, tokens.user);
      queryClient.invalidateQueries({ queryKey: ['auth'] });
      router.replace('/admin');
    },
  });
}

export function useLogout() {
  const router = useRouter();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const auth = readAuthState();
      if (auth.accessToken && auth.refreshToken) {
        await api.auth.logout(auth.refreshToken, auth.accessToken);
      }
    },
    onSettled: () => {
      clearAuthState();
      queryClient.clear();
      router.replace('/login');
    },
  });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: UpdateProfileRequest) => {
      const auth = readAuthState();
      if (!auth.accessToken) {
        throw new Error('Authentication required');
      }
      return api.auth.updateProfile(body, auth.accessToken);
    },
    onSuccess: (user) => {
      const auth = readAuthState();
      if (auth.accessToken && auth.refreshToken) {
        saveAuthState({ accessToken: auth.accessToken, refreshToken: auth.refreshToken, expiresIn: 0, user }, user);
      }
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
}

export function useChangePassword() {
  const logout = useLogout();
  return useMutation({
    mutationFn: (body: ChangePasswordRequest) => {
      const auth = readAuthState();
      if (!auth.accessToken) {
        throw new Error('Authentication required');
      }
      return api.auth.changePassword(body, auth.accessToken);
    },
    onSuccess: () => logout.mutate(),
  });
}
