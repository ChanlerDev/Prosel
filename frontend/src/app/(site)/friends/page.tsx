import { FriendGrid } from '@/components/features/page/friend-grid';
import { EmptyState } from '@/components/features/system/states';
import { SiteContainer } from '@/components/layout/site-container';
import { SiteFooter } from '@/components/layout/site-footer';
import { SiteHeader } from '@/components/layout/site-header';
import { pagesApi } from '@/lib/api/pages';

export const revalidate = 60;

export default async function FriendsPage() {
  const friends = await pagesApi.friends().catch(() => []);

  return (
    <>
      <SiteHeader />
      <main className="py-12">
        <SiteContainer>
          <div className="mx-auto mb-10 max-w-3xl">
            <h1 className="text-5xl font-semibold tracking-tight">Friends</h1>
            <p className="mt-4 text-lg text-[var(--muted-foreground)]">People and sites worth visiting.</p>
          </div>
          <div className="mx-auto max-w-4xl">
            {friends.length > 0 ? <FriendGrid friends={friends} /> : <EmptyState title="No friends yet" description="Friend links will appear here once added." />}
          </div>
        </SiteContainer>
      </main>
      <SiteFooter />
    </>
  );
}
