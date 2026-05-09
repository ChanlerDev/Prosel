'use client';

import { useState } from 'react';

import { ApiErrorState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useMe, useUpdateProfile } from '@/lib/auth/hooks';
import type { User } from '@/types/api';

export function ProfileForm() {
  const me = useMe();
  const update = useUpdateProfile();

  if (me.isLoading) return <LoadingState />;
  if (me.isError || !me.data) return <ApiErrorState message="Unable to load profile." />;

  return <ProfileFormFields key={me.data.id} initialUser={me.data} isPending={update.isPending} onSubmit={(body) => update.mutate(body)} error={update.isError ? update.error.message : null} />;
}

function ProfileFormFields({
  initialUser,
  isPending,
  error,
  onSubmit,
}: {
  initialUser: User;
  isPending: boolean;
  error: string | null;
  onSubmit: (body: { displayName: string; avatarUrl: string; bio: string }) => void;
}) {
  const [displayName, setDisplayName] = useState(initialUser.displayName);
  const [avatarUrl, setAvatarUrl] = useState(initialUser.avatarUrl ?? '');
  const [bio, setBio] = useState(initialUser.bio ?? '');

  return (
    <Card className="max-w-xl">
      <form
        className="grid gap-4"
        onSubmit={(event) => {
          event.preventDefault();
          onSubmit({ displayName, avatarUrl, bio });
        }}
      >
        <label className="grid gap-2 text-sm">
          Display name
          <Input onChange={(event) => setDisplayName(event.target.value)} value={displayName} />
        </label>
        <label className="grid gap-2 text-sm">
          Avatar URL
          <Input onChange={(event) => setAvatarUrl(event.target.value)} value={avatarUrl} />
        </label>
        <label className="grid gap-2 text-sm">
          Bio
          <Input onChange={(event) => setBio(event.target.value)} value={bio} />
        </label>
        {error ? <ApiErrorState message={error} /> : null}
        <Button disabled={isPending} type="submit">
          {isPending ? 'Saving...' : 'Save profile'}
        </Button>
      </form>
    </Card>
  );
}
