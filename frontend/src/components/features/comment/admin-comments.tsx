'use client';

import { useState } from 'react';

import { ApiErrorState, EmptyState, LoadingState } from '@/components/features/system/states';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useAdminComments, useAdminReply, useDeleteComment, useModerateComment } from '@/lib/comments/hooks';
import type { CommentNode, CommentRefType, CommentStatus } from '@/types/comment';

export function AdminComments() {
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<CommentStatus | ''>('');
  const [refType, setRefType] = useState<CommentRefType | ''>('');
  const comments = useAdminComments({ search, status, refType });
  const moderate = useModerateComment();
  const remove = useDeleteComment();
  const reply = useAdminReply();
  const [replyingTo, setReplyingTo] = useState<string | null>(null);
  const [replyContent, setReplyContent] = useState('');

  async function submitReply(commentID: string) {
    if (!replyContent.trim()) return;
    await reply.mutateAsync({ id: commentID, content: replyContent });
    setReplyContent('');
    setReplyingTo(null);
  }

  return (
    <div className="grid gap-5">
      <Card className="grid gap-3 md:grid-cols-[1fr_auto_auto] md:items-center">
        <Input onChange={(event) => setSearch(event.target.value)} placeholder="Search author, email, or content" value={search} />
        <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setStatus(event.target.value as CommentStatus | '')} value={status}>
          <option value="">All status</option>
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
          <option value="spam">Spam</option>
        </select>
        <select className="rounded-xl border border-[var(--border)] bg-transparent px-3 py-2 text-sm" onChange={(event) => setRefType(event.target.value as CommentRefType | '')} value={refType}>
          <option value="">All refs</option>
          <option value="post">Posts</option>
          <option value="note">Notes</option>
          <option value="page">Pages</option>
        </select>
      </Card>
      {comments.isLoading ? <LoadingState /> : null}
      {comments.isError ? <ApiErrorState message={comments.error.message} /> : null}
      {comments.data && comments.data.comments.length === 0 ? <EmptyState title="No comments" description="New visitor comments will appear here for moderation." /> : null}
      {comments.data?.comments.map((comment) => (
        <CommentModerationCard comment={comment} isBusy={moderate.isPending || remove.isPending || reply.isPending} key={comment.id} onDelete={() => remove.mutate(comment.id)} onModerate={(nextStatus) => moderate.mutate({ id: comment.id, status: nextStatus })} onReply={() => submitReply(comment.id)} replyContent={replyContent} replying={replyingTo === comment.id} setReplyContent={setReplyContent} setReplying={(value) => setReplyingTo(value ? comment.id : null)} />
      ))}
      {comments.data ? <p className="text-sm text-[var(--muted-foreground)]">Page {comments.data.meta.page} of {comments.data.meta.totalPages || 1} · {comments.data.meta.total} comments</p> : null}
    </div>
  );
}

function CommentModerationCard({ comment, isBusy, onDelete, onModerate, onReply, replyContent, replying, setReplyContent, setReplying }: { comment: CommentNode; isBusy: boolean; onDelete: () => void; onModerate: (status: CommentStatus) => void; onReply: () => void; replyContent: string; replying: boolean; setReplyContent: (value: string) => void; setReplying: (value: boolean) => void }) {
  return (
    <Card className="grid gap-4">
      <div className="flex flex-wrap items-center gap-2">
        <StatusBadge status={comment.status} />
        {comment.isAdminReply ? <span className="rounded-full bg-[var(--primary)] px-2 py-0.5 text-xs font-medium text-[var(--primary-foreground)]">Author</span> : null}
        <p className="font-semibold">{comment.authorName}</p>
        {comment.authorEmail ? <p className="text-xs text-[var(--muted-foreground)]">{comment.authorEmail}</p> : null}
      </div>
      <p className="whitespace-pre-wrap text-sm leading-7">{comment.content}</p>
      <p className="text-xs text-[var(--muted-foreground)]">{comment.refType}:{comment.refId} · {formatDate(comment.createdAt)}</p>
      <div className="flex flex-wrap gap-2">
        <Button className="px-3 py-1 text-xs" disabled={isBusy} onClick={() => onModerate('approved')} type="button">Approve</Button>
        <Button className="px-3 py-1 text-xs" disabled={isBusy} onClick={() => onModerate('rejected')} type="button">Reject</Button>
        <Button className="px-3 py-1 text-xs" disabled={isBusy} onClick={() => onModerate('spam')} type="button">Spam</Button>
        <Button className="px-3 py-1 text-xs" disabled={isBusy} onClick={() => setReplying(!replying)} type="button">Reply</Button>
        <Button className="bg-red-600 px-3 py-1 text-xs text-white" disabled={isBusy} onClick={onDelete} type="button">Delete</Button>
      </div>
      {replying ? (
        <div className="grid gap-3 border-l border-[var(--border)] pl-4">
          <Textarea onChange={(event) => setReplyContent(event.target.value)} placeholder="Write author reply" value={replyContent} />
          <Button disabled={isBusy || !replyContent.trim()} onClick={onReply} type="button">Send reply</Button>
        </div>
      ) : null}
    </Card>
  );
}

function StatusBadge({ status }: { status: CommentStatus }) {
  const color = status === 'approved' ? 'bg-emerald-100 text-emerald-800' : status === 'pending' ? 'bg-amber-100 text-amber-900' : 'bg-zinc-100 text-zinc-700';
  return <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${color}`}>{status}</span>;
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value));
}
