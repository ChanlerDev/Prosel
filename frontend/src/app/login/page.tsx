'use client';

import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center px-6">
      <Card className="w-full max-w-md">
        <p className="text-sm text-[var(--muted-foreground)]">Admin login</p>
        <h1 className="mt-2 text-2xl font-semibold">Sign in to Prosel</h1>
        <form className="mt-6 grid gap-4">
          <Input disabled placeholder="Email (coming in Auth module)" type="email" />
          <Input disabled placeholder="Password" type="password" />
          <Button disabled type="button">
            Auth module pending
          </Button>
        </form>
      </Card>
    </main>
  );
}
