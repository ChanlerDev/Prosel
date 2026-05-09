import type { Post } from '@/types/post';

export function PostContent({ post }: { post: Post }) {
  return (
    <article className="prose prose-neutral max-w-none">
      {(post.contentMarkdown ?? '').split('\n').map((line, index) => (
        <MarkdownLine key={`${index}-${line}`} line={line} />
      ))}
    </article>
  );
}

function MarkdownLine({ line }: { line: string }) {
  const trimmed = line.trim();
  if (!trimmed) return <div className="h-4" />;
  if (trimmed.startsWith('### ')) return <h3 className="mt-8 text-2xl font-semibold">{trimmed.slice(4)}</h3>;
  if (trimmed.startsWith('## ')) return <h2 className="mt-10 text-3xl font-semibold">{trimmed.slice(3)}</h2>;
  if (trimmed.startsWith('# ')) return <h1 className="mt-10 text-4xl font-semibold">{trimmed.slice(2)}</h1>;
  if (trimmed.startsWith('- ')) return <p className="pl-4 text-[var(--muted-foreground)]">• {trimmed.slice(2)}</p>;
  return <p className="mt-4 leading-8 text-[var(--foreground)]">{trimmed.replaceAll('**', '').replaceAll('`', '')}</p>;
}
