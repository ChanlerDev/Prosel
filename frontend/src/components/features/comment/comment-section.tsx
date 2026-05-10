'use client';

import { useState } from 'react';

import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useComments, useSubmitComment } from '@/lib/comments/hooks';
import type { CommentNode, CommentRefType, SubmitCommentValues } from '@/types/comment';

type CommentFormValues = Omit<SubmitCommentValues, 'refType' | 'refId' | 'parentId'>;

const emptyForm: CommentFormValues = { authorName: '', authorEmail: '', authorWebsite: '', content: '' };

export function CommentSection({ refType, refId }: { refType: CommentRefType; refId: string }) {
  const comments = useComments(refType, refId);
  const submit = useSubmitComment(refType, refId);
  const [notice, setNotice] = useState('');

  async function handleSubmit(values: CommentFormValues, parentId?: string) {
    await submit.mutateAsync({ refType, refId, parentId, ...values });
    setNotice('Thanks! Your comment is waiting for approval.');
  }

  return (
    <section className="mx-auto mt-12 max-w-3xl" id="comments">
      <div className="mb-5">
        <h2 className="text-2xl font-semibold tracking-tight">Comments</h2>
        <p className="mt-2 text-sm text-[var(--muted-foreground)]">Join the discussion. Comments appear after moderation.</p>
      </div>
      <Card className="mb-5">
        <CommentForm isPending={submit.isPending} onSubmit={(values) => handleSubmit(values)} />
        {notice ? <p className="mt-3 text-sm text-emerald-700">{notice}</p> : null}
        {submit.isError ? <p className="mt-3 text-sm text-red-600">{submit.error.message}</p> : null}
      </Card>
      {comments.isLoading ? <LoadingState /> : null}
      {comments.isError ? <ApiErrorState message={comments.error.message} /> : null}
      {comments.data && comments.data.length === 0 ? <EmptyState title="No approved comments yet" description="Be the first to share a thought." /> : null}
      {comments.data && comments.data.length > 0 ? <CommentList comments={comments.data} isPending={submit.isPending} onReply={handleSubmit} /> : null}
    </section>
  );
}

function CommentList({ comments, isPending, onReply }: { comments: CommentNode[]; isPending: boolean; onReply: (values: CommentFormValues, parentId: string) => Promise<void> }) {
  return (
    <div className="grid gap-4">
      {comments.map((comment) => (
        <CommentItem comment={comment} isPending={isPending} key={comment.id} onReply={onReply} />
      ))}
    </div>
  );
}

function CommentItem({ comment, isPending, onReply }: { comment: CommentNode; isPending: boolean; onReply: (values: CommentFormValues, parentId: string) => Promise<void> }) {
  const [replying, setReplying] = useState(false);
  const children = comment.children ?? [];

  return (
    <Card className={comment.isAdminReply ? 'border-[var(--primary)] bg-[var(--secondary)]' : ''}>
      <div className="flex flex-wrap items-center gap-2">
        {comment.authorWebsite ? (
          <a className="font-semibold hover:text-[var(--primary)]" href={comment.authorWebsite} rel="noreferrer" target="_blank">
            {comment.authorName}
          </a>
        ) : (
          <p className="font-semibold">{comment.authorName}</p>
        )}
        {comment.isAdminReply ? <AdminReplyBadge /> : null}
        <span className="text-xs text-[var(--muted-foreground)]">{formatDate(comment.createdAt)}</span>
      </div>
      <p className="mt-3 whitespace-pre-wrap text-sm leading-7 text-[var(--foreground)]">{comment.content}</p>
      <Button className="mt-4 px-3 py-1 text-xs" onClick={() => setReplying((value) => !value)} type="button">
        Reply
      </Button>
      {replying ? (
        <div className="mt-4 border-l border-[var(--border)] pl-4">
          <CommentForm isPending={isPending} onSubmit={async (values) => { await onReply(values, comment.id); setReplying(false); }} />
        </div>
      ) : null}
      {children.length > 0 ? (
        <div className="mt-4 grid gap-4 border-l border-[var(--border)] pl-4">
          {children.map((child) => (
            <CommentItem comment={child} isPending={isPending} key={child.id} onReply={onReply} />
          ))}
        </div>
      ) : null}
    </Card>
  );
}

function CommentForm({ isPending, onSubmit }: { isPending: boolean; onSubmit: (values: CommentFormValues) => Promise<void> }) {
  const [values, setValues] = useState<CommentFormValues>(emptyForm);

  return (
    <form className="grid gap-3" onSubmit={async (event) => { event.preventDefault(); await onSubmit(values); setValues(emptyForm); }}>
      <div className="grid gap-3 md:grid-cols-2">
        <Input onChange={(event) => setValues({ ...values, authorName: event.target.value })} placeholder="Name" required value={values.authorName} />
        <Input onChange={(event) => setValues({ ...values, authorEmail: event.target.value })} placeholder="Email" required type="email" value={values.authorEmail} />
      </div>
      <Input onChange={(event) => setValues({ ...values, authorWebsite: event.target.value })} placeholder="Website (optional)" type="url" value={values.authorWebsite} />
      <Textarea onChange={(event) => setValues({ ...values, content: event.target.value })} placeholder="Write a thoughtful comment" required rows={5} value={values.content} />
      <Button disabled={isPending} type="submit">{isPending ? 'Submitting...' : 'Submit comment'}</Button>
    </form>
  );
}

function AdminReplyBadge() {
  return <span className="rounded-full bg-[var(--primary)] px-2 py-0.5 text-xs font-medium text-[var(--primary-foreground)]">Author</span>;
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium' }).format(new Date(value));
}
