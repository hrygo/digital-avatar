import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import { Copy, Check } from 'lucide-react';
import clsx from 'clsx';
import 'highlight.js/styles/github-dark.css';

interface MarkdownRendererProps {
  content: string;
  className?: string;
  enableCopy?: boolean;
}

export const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({
  content,
  className = '',
  enableCopy = true,
}) => {
  const [copiedCode, setCopiedCode] = React.useState<string>('');

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      setCopiedCode(text);
      setTimeout(() => setCopiedCode(''), 2000);
    } catch (err) {
      console.error('Failed to copy text: ', err);
    }
  };

  const components = {
    // 标题组件
    h1: ({ children, ...props }: any) => (
      <h1
        className={clsx(
          'text-3xl font-bold text-white mb-6 mt-8 first:mt-0',
          'border-b border-gray-700 pb-3'
        )}
        {...props}
      >
        {children}
      </h1>
    ),

    h2: ({ children, ...props }: any) => (
      <h2
        className={clsx(
          'text-2xl font-bold text-white mb-4 mt-6 first:mt-0',
          'border-l-4 border-primary-500 pl-4'
        )}
        {...props}
      >
        {children}
      </h2>
    ),

    h3: ({ children, ...props }: any) => (
      <h3
        className={clsx(
          'text-xl font-semibold text-white mb-3 mt-5 first:mt-0',
          'text-gray-200'
        )}
        {...props}
      >
        {children}
      </h3>
    ),

    h4: ({ children, ...props }: any) => (
      <h4
        className={clsx(
          'text-lg font-medium text-white mb-2 mt-4 first:mt-0',
          'text-gray-300'
        )}
        {...props}
      >
        {children}
      </h4>
    ),

    // 段落组件
    p: ({ children, ...props }: any) => (
      <p
        className={clsx(
          'text-gray-300 leading-relaxed mb-4 first:mt-0',
          'text-justify'
        )}
        {...props}
      >
        {children}
      </p>
    ),

    // 列表组件
    ul: ({ children, ...props }: any) => (
      <ul
        className={clsx(
          'list-disc list-inside text-gray-300 mb-4 space-y-2',
          'ml-4'
        )}
        {...props}
      >
        {children}
      </ul>
    ),

    ol: ({ children, ...props }: any) => (
      <ol
        className={clsx(
          'list-decimal list-inside text-gray-300 mb-4 space-y-2',
          'ml-4'
        )}
        {...props}
      >
        {children}
      </ol>
    ),

    li: ({ children, ...props }: any) => (
      <li
        className={clsx(
          'leading-relaxed',
          'marker:text-primary-400'
        )}
        {...props}
      >
        {children}
      </li>
    ),

    // 强调组件
    strong: ({ children, ...props }: any) => (
      <strong
        className={clsx(
          'text-primary-400 font-semibold',
          'font-bold'
        )}
        {...props}
      >
        {children}
      </strong>
    ),

    em: ({ children, ...props }: any) => (
      <em
        className={clsx(
          'text-accent-400 italic'
        )}
        {...props}
      >
        {children}
      </em>
    ),

    // 代码组件
    code: ({ inline, children, ...props }: any) => {
      const codeContent = String(children).replace(/\n$/, '');

      if (inline) {
        return (
          <code
            className={clsx(
              'bg-gray-800 text-pink-400 px-2 py-1 rounded',
              'text-sm font-mono'
            )}
            {...props}
          >
            {children}
          </code>
        );
      }

      return (
        <div className="relative group">
          {enableCopy && (
            <div className="absolute right-2 top-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button
                onClick={() => copyToClipboard(codeContent)}
                className={clsx(
                  'p-2 bg-gray-700 hover:bg-gray-600 rounded',
                  'text-gray-300 hover:text-white',
                  'transition-colors text-xs'
                )}
                title="复制代码"
              >
                {copiedCode === codeContent ? (
                  <Check className="w-4 h-4 text-success-400" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </button>
            </div>
          )}
          <code
            className={clsx(
              'block bg-gray-900 text-gray-300 p-4 rounded-lg',
              'overflow-x-auto font-mono text-sm',
              'border border-gray-700'
            )}
            {...props}
          >
            {children}
          </code>
        </div>
      );
    },

    // 代码块
    pre: ({ children, ...props }: any) => (
      <pre
        className={clsx(
          'bg-gray-900 p-4 rounded-lg overflow-x-auto mb-4',
          'border border-gray-700'
        )}
        {...props}
      >
        {children}
      </pre>
    ),

    // 引用块
    blockquote: ({ children, ...props }: any) => (
      <blockquote
        className={clsx(
          'border-l-4 border-primary-500 pl-4 py-2 my-4',
          'bg-gray-800/50 rounded-r-lg',
          'text-gray-300 italic'
        )}
        {...props}
      >
        {children}
      </blockquote>
    ),

    // 链接
    a: ({ children, href, ...props }: any) => (
      <a
        href={href}
        target="_blank"
        rel="noopener noreferrer"
        className={clsx(
          'text-primary-400 hover:text-primary-300 underline',
          'transition-colors'
        )}
        {...props}
      >
        {children}
      </a>
    ),

    // 表格
    table: ({ children, ...props }: any) => (
      <div className="overflow-x-auto mb-4">
        <table
          className={clsx(
            'min-w-full border-collapse',
            'border border-gray-700'
          )}
          {...props}
        >
          {children}
        </table>
      </div>
    ),

    thead: ({ children, ...props }: any) => (
      <thead
        className={clsx(
          'bg-gray-800'
        )}
        {...props}
      >
        {children}
      </thead>
    ),

    th: ({ children, ...props }: any) => (
      <th
        className={clsx(
          'border border-gray-700 px-4 py-2 text-left',
          'text-gray-200 font-semibold',
          'text-sm'
        )}
        {...props}
      >
        {children}
      </th>
    ),

    td: ({ children, ...props }: any) => (
      <td
        className={clsx(
          'border border-gray-700 px-4 py-2',
          'text-gray-300 text-sm'
        )}
        {...props}
      >
        {children}
      </td>
    ),

    // 分隔线
    hr: ({ ...props }: any) => (
      <hr
        className={clsx(
          'border-gray-700 my-6',
          'border-t'
        )}
        {...props}
      />
    ),

    // 删除线
    del: ({ children, ...props }: any) => (
      <del
        className={clsx(
          'text-gray-500 line-through'
        )}
        {...props}
      >
        {children}
      </del>
    ),

    // 任务列表
    input: ({ type, checked, ...props }: any) => {
      if (type === 'checkbox') {
        return (
          <input
            type="checkbox"
            checked={checked}
            readOnly
            className={clsx(
              'mr-2 w-4 h-4 text-primary-500',
              'bg-gray-700 border-gray-600 rounded',
              'focus:ring-primary-500 focus:ring-2'
            )}
            {...props}
          />
        );
      }
      return <input {...props} />;
    },
  };

  if (!content || content.trim() === '') {
    return (
      <div className={clsx('text-gray-500 text-center py-8', className)}>
        暂无内容
      </div>
    );
  }

  return (
    <div className={clsx('prose prose-invert max-w-none', className)}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight]}
        components={components}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
};

export default MarkdownRenderer;