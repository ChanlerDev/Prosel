'use client';

import { useState } from 'react';

import { ApiErrorState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useChangePassword } from '@/lib/auth/hooks';

export function PasswordForm() {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const mutation = useChangePassword();

  return (
    <Card className="max-w-xl">
      <form
        className="grid gap-4"
        onSubmit={(event) => {
          event.preventDefault();
          mutation.mutate({ oldPassword, newPassword });
        }}
      >
        <label className="grid gap-2 text-sm">
          Current password
          <Input autoComplete="current-password" onChange={(event) => setOldPassword(event.target.value)} type="password" value={oldPassword} />
        </label>
        <label className="grid gap-2 text-sm">
          New password
          <Input autoComplete="new-password" onChange={(event) => setNewPassword(event.target.value)} type="password" value={newPassword} />
        </label>
        {mutation.isError ? <ApiErrorState message={mutation.error.message} /> : null}
        <Button disabled={mutation.isPending} type="submit">
          {mutation.isPending ? 'Updating...' : 'Change password'}
        </Button>
      </form>
    </Card>
  );
}
