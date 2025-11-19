import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import { Copy, CheckCircle } from 'lucide-react';
import clsx from 'clsx';

interface EnhancedMarkdownRendererProps {
  content: string;
  className?: string;
  enableCopy?: boolean;
  maxHeight?: string;
}

export const EnhancedMarkdownRenderer: React.FC<EnhancedMarkdownRendererProps> = ({
  content,
  className,
  enableCopy = true,
  maxHeight = 'none',
}) => {
  const [copied, setCopied] = React.useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(content);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('复制失败:', err);
    }
  };

  const components = {
    h1: ({ children, ...props }: any) => (
      <h1 className={clsx('text-3xl font-bold text-white mb-4 mt-6', className)} {...props}>
        {children}
      </h1>
    ),
    h2: ({ children, ...props }: any) => (
      <h2 className={clsx('text-2xl font-semibold text-blue-400 mb-3 mt-5', className)} {...props}>
        {children}
      </h2>
    ),
    h3: ({ children, ...props }: any) => (
      <h3 className={clsx('text-xl font-semibold text-green-400 mb-2 mt-4', className)} {...props}>
        {children}
      </h3>
    ),
    h4: ({ children, ...props }: any) => (
      <h4 className={clsx('text-lg font-semibold text-yellow-400 mb-2 mt-3', className)} {...props}>
        {children}
      </h4>
    ),
    p: ({ children, ...props }: any) => (
      <p className={clsx('text-gray-300 leading-relaxed mb-3', className)} {...props}>
        {children}
      </p>
    ),
    ul: ({ children, ...props }: any) => (
      <ul className={clsx('list-disc list-inside text-gray-300 space-y-2 mb-4', className)} {...props}>
        {children}
      </ul>
    ),
    ol: ({ children, ...props }: any) => (
      <ol className={clsx('list-decimal list-inside text-gray-300 space-y-2 mb-4', className)} {...props}>
        {children}
      </ol>
    ),
    li: ({ children, ...props }: any) => (
      <li className={clsx('leading-relaxed', className)} {...props}>
        {children}
      </li>
    ),
    strong: ({ children, ...props }: any) => (
      <strong className={clsx('text-pink-400 font-semibold', className)} {...props}>
        {children}
      </strong>
    ),
    em: ({ children, ...props }: any) => (
      <em className={clsx('text-purple-400 italic', className)} {...props}>
        {children}
      </em>
    ),
    a: ({ href, children, ...props }: any) => (
      <a
        href={href}
        className={clsx('text-blue-400 hover:text-blue-300 underline transition-colors', className)}
        target="_blank"
        rel="noopener noreferrer"
        {...props}
      >
        {children}
      </a>
    ),
    blockquote: ({ children, ...props }: any) => (
      <blockquote
        className={clsx(
          'border-l-4 border-blue-500 pl-4 py-2 my-4 bg-blue-500/10 rounded-r-lg',
          className
        )}
        {...props}
      >
        {children}
      </blockquote>
    ),
    code: ({ inline, children, ...props }: any) => {
      if (inline) {
        return (
          <code
            className={clsx(
              'bg-gray-800 text-yellow-400 px-2 py-1 rounded text-sm font-mono',
              className
            )}
            {...props}
          >
            {children}
          </code>
        );
      }
      return null;
    },
    pre: ({ children, ...props }: any) => (
      <pre
        className={clsx(
          'bg-gray-900 p-4 rounded-lg overflow-x-auto border border-gray-700 mb-4',
          className
        )}
        {...props}
      >
        {children}
      </pre>
    ),
    table: ({ children, ...props }: any) => (
      <div className={clsx('overflow-x-auto mb-4', className)}>
        <table className={clsx('w-full border-collapse border border-gray-700', className)} {...props}>
          {children}
        </table>
      </div>
    ),
    thead: ({ children, ...props }: any) => (
      <thead className={clsx('bg-gray-800', className)} {...props}>
        {children}
      </thead>
    ),
    th: ({ children, ...props }: any) => (
      <th
        className={clsx('border border-gray-700 px-4 py-2 text-left text-gray-300 font-semibold', className)}
        {...props}
      >
        {children}
      </th>
    ),
    td: ({ children, ...props }: any) => (
      <td className={clsx('border border-gray-700 px-4 py-2 text-gray-300', className)} {...props}>
        {children}
      </td>
    ),
    hr: ({ ...props }: any) => (
      <hr className={clsx('border-gray-700 my-6', className)} {...props} />
    ),
  };

  return (
    <div className={clsx('relative', className)}>
      {enableCopy && (
        <button
          onClick={handleCopy}
          className="absolute top-2 right-2 p-2 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-md transition-colors"
          title="复制内容"
        >
          {copied ? (
            <CheckCircle className="w-4 h-4 text-green-400" />
          ) : (
            <Copy className="w-4 h-4" />
          )}
        </button>
      )}
      <div
        className={clsx('prose prose-invert max-w-none', {
          'overflow-y-auto': maxHeight !== 'none',
        })}
        style={{ maxHeight }}
      >
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          rehypePlugins={[rehypeHighlight]}
          components={components}
        >
          {content}
        </ReactMarkdown>
      </div>
    </div>
  );
};

export default EnhancedMarkdownRenderer;