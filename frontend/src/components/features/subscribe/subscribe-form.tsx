'use client';

import { useState } from 'react';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useSubscribe } from '@/lib/subscribe/hooks';

export function SubscribeForm() {
  const subscribe = useSubscribe();
  const [email, setEmail] = useState('');
  const [name, setName] = useState('');

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!email.trim()) return;
    await subscribe.mutateAsync({ email, name });
    setEmail('');
    setName('');
  }

  return (
    <form className="grid gap-3 sm:grid-cols-[1fr_1fr_auto]" onSubmit={onSubmit}>
      <Input aria-label="Name" onChange={(event) => setName(event.target.value)} placeholder="Name (optional)" value={name} />
      <Input aria-label="Email" onChange={(event) => setEmail(event.target.value)} placeholder="you@example.com" type="email" value={email} />
      <Button disabled={subscribe.isPending || !email.trim()} type="submit">
        {subscribe.isPending ? 'Subscribing...' : 'Subscribe'}
      </Button>
      {subscribe.isSuccess ? <p className="text-sm text-emerald-700 sm:col-span-3">Check your inbox to verify the subscription.</p> : null}
      {subscribe.isError ? <p className="text-sm text-red-700 sm:col-span-3">{subscribe.error.message}</p> : null}
    </form>
  );
}
