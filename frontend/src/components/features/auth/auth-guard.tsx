'use client';

import { useRouter } from 'next/navigation';
import { useEffect, type ReactNode } from 'react';

import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { useMe } from '@/lib/auth/hooks';
import { readAuthState } from '@/lib/auth/store';

export function AuthGuard({ children }: { children: ReactNode }) {
  const router = useRouter();
  const auth = readAuthState();
  const me = useMe();

  useEffect(() => {
    if (!auth.accessToken || me.isError) {
      router.replace('/login');
    }
  }, [auth.accessToken, me.isError, router]);

  if (!auth.accessToken) {
    return <LoadingState />;
  }
  if (me.isLoading) {
    return <LoadingState />;
  }
  if (me.isError) {
    return <ApiErrorState message="Your session has expired. Redirecting to login." />;
  }
  return <>{children}</>;
}
