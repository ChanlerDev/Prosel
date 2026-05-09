'use client';

import { useState } from 'react';

import { ApiErrorState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useLogin } from '@/lib/auth/hooks';

export function LoginForm() {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const mutation = useLogin();

  return (
    <form
      className="mt-6 grid gap-4"
      onSubmit={(event) => {
        event.preventDefault();
        mutation.mutate({ login, password });
      }}
    >
      <Input autoComplete="username" onChange={(event) => setLogin(event.target.value)} placeholder="Username or email" value={login} />
      <Input autoComplete="current-password" onChange={(event) => setPassword(event.target.value)} placeholder="Password" type="password" value={password} />
      {mutation.isError ? <ApiErrorState message={mutation.error.message} /> : null}
      <Button disabled={mutation.isPending} type="submit">
        {mutation.isPending ? 'Signing in...' : 'Sign in'}
      </Button>
    </form>
  );
}
