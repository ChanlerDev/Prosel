import Image from 'next/image';

import { Card } from '@/components/ui/card';
import type { Friend } from '@/types/page';

export function FriendGrid({ friends }: { friends: Friend[] }) {
  return (
    <div className="grid gap-4 md:grid-cols-2">
      {friends.map((friend) => (
        <FriendCard friend={friend} key={friend.id} />
      ))}
    </div>
  );
}

function FriendCard({ friend }: { friend: Friend }) {
  return (
    <a href={friend.url} rel="noreferrer" target="_blank">
      <Card className="h-full transition hover:-translate-y-0.5 hover:border-[var(--primary)]">
        <div className="flex items-center gap-4">
          {friend.avatarUrl ? <Image alt="" className="h-12 w-12 rounded-full object-cover" height={48} src={friend.avatarUrl} width={48} /> : <div className="flex h-12 w-12 items-center justify-center rounded-full bg-[var(--secondary)] text-lg font-semibold">{friend.name.slice(0, 1)}</div>}
          <div>
            <h2 className="font-semibold">{friend.name}</h2>
            <p className="text-xs text-[var(--muted-foreground)]">{new URL(friend.url).hostname}</p>
          </div>
        </div>
        {friend.description ? <p className="mt-4 text-sm leading-6 text-[var(--muted-foreground)]">{friend.description}</p> : null}
      </Card>
    </a>
  );
}
