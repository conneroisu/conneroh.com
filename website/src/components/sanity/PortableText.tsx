import { For, Show, type Component } from 'solid-js';
import type { BlockContent } from '@/lib/sanity-types';

/**
 * Simple Portable Text renderer for Solid.js
 *
 * Renders Sanity's BlockContent structure into HTML
 */

interface PortableTextProps {
  content: BlockContent | undefined;
  class?: string;
}

interface TextSpan {
  text?: string;
  marks?: string[];
  _type: 'span';
  _key: string;
}

interface MarkDef {
  _key: string;
  _type: string;
  href?: string;
  reference?: {
    _ref: string;
    _type: string;
  };
}

interface Block {
  _type: 'block';
  _key: string;
  style?: 'normal' | 'h1' | 'h2' | 'h3' | 'blockquote';
  listItem?: 'bullet' | 'number';
  children?: TextSpan[];
  markDefs?: MarkDef[];
  level?: number;
}

interface ImageBlock {
  _type: 'image';
  _key: string;
  asset?: {
    _ref: string;
  };
  alt?: string;
  caption?: string;
}

type ContentBlock = Block | ImageBlock;

/**
 * Render a text span with marks (bold, italic, links, etc.)
 */
const RenderSpan: Component<{ span: TextSpan; markDefs?: MarkDef[] }> = (props) => {
  const hasMarks = () => props.span.marks && props.span.marks.length > 0;

  const renderWithMarks = (text: string, marks: string[], markDefs?: MarkDef[]): any => {
    if (!marks || marks.length === 0) {
      return text;
    }

    const [firstMark, ...restMarks] = marks;

    // Handle link marks
    if (markDefs) {
      const markDef = markDefs.find(def => def._key === firstMark);
      if (markDef) {
        if (markDef._type === 'link' && markDef.href) {
          return (
            <a
              href={markDef.href}
              class="text-blue-600 hover:text-blue-800 underline"
              target="_blank"
              rel="noopener noreferrer"
            >
              {renderWithMarks(text, restMarks, markDefs)}
            </a>
          );
        }
        if (markDef._type === 'internalLink' && markDef.reference) {
          // TODO: Implement internal link routing
          return (
            <a
              href={`#${markDef.reference._ref}`}
              class="text-blue-600 hover:text-blue-800 underline"
            >
              {renderWithMarks(text, restMarks, markDefs)}
            </a>
          );
        }
      }
    }

    // Handle standard marks
    switch (firstMark) {
      case 'strong':
        return <strong>{renderWithMarks(text, restMarks, markDefs)}</strong>;
      case 'em':
        return <em>{renderWithMarks(text, restMarks, markDefs)}</em>;
      case 'code':
        return (
          <code class="bg-gray-100 px-1 py-0.5 rounded text-sm font-mono">
            {renderWithMarks(text, restMarks, markDefs)}
          </code>
        );
      case 'underline':
        return <u>{renderWithMarks(text, restMarks, markDefs)}</u>;
      case 'strike-through':
        return <s>{renderWithMarks(text, restMarks, markDefs)}</s>;
      default:
        return renderWithMarks(text, restMarks, markDefs);
    }
  };

  return (
    <Show
      when={hasMarks()}
      fallback={<>{props.span.text}</>}
    >
      {renderWithMarks(props.span.text || '', props.span.marks || [], props.markDefs)}
    </Show>
  );
};

/**
 * Render a block element (paragraph, heading, blockquote, etc.)
 */
const RenderBlock: Component<{ block: Block }> = (props) => {
  const style = () => props.block.style || 'normal';
  const children = () => (
    <For each={props.block.children}>
      {(child) => (
        <RenderSpan span={child} markDefs={props.block.markDefs} />
      )}
    </For>
  );

  return (
    <Show
      when={style() === 'h1'}
      fallback={
        <Show
          when={style() === 'h2'}
          fallback={
            <Show
              when={style() === 'h3'}
              fallback={
                <Show
                  when={style() === 'blockquote'}
                  fallback={
                    <Show
                      when={props.block.listItem === 'bullet'}
                      fallback={
                        <Show
                          when={props.block.listItem === 'number'}
                          fallback={
                            <p class="mb-4 leading-relaxed">{children()}</p>
                          }
                        >
                          <li class="ml-4">{children()}</li>
                        </Show>
                      }
                    >
                      <li class="ml-4">{children()}</li>
                    </Show>
                  }
                >
                  <blockquote class="border-l-4 border-gray-300 pl-4 italic mb-4">
                    {children()}
                  </blockquote>
                </Show>
              }
            >
              <h3 class="text-2xl font-semibold mb-3 mt-6">{children()}</h3>
            </Show>
          }
        >
          <h2 class="text-3xl font-semibold mb-4 mt-8">{children()}</h2>
        </Show>
      }
    >
      <h1 class="text-4xl font-bold mb-6 mt-10">{children()}</h1>
    </Show>
  );
};

/**
 * Render an image block
 */
const RenderImage: Component<{ block: ImageBlock }> = (props) => {
  // For now, we'll just render a placeholder
  // In production, you'd use Sanity's image URL builder
  return (
    <figure class="my-6">
      <div class="bg-gray-200 rounded aspect-video flex items-center justify-center">
        <span class="text-gray-500">Image: {props.block.asset?._ref}</span>
      </div>
      <Show when={props.block.caption}>
        <figcaption class="text-sm text-gray-600 mt-2 text-center">
          {props.block.caption}
        </figcaption>
      </Show>
    </figure>
  );
};

/**
 * Main Portable Text component
 */
export const PortableText: Component<PortableTextProps> = (props) => {
  return (
    <Show
      when={props.content && props.content.length > 0}
      fallback={<p class="text-gray-500 italic">No content available</p>}
    >
      <div class={props.class || 'prose prose-lg max-w-none'}>
        <For each={props.content as ContentBlock[]}>
          {(block) => (
            <Show
              when={block._type === 'block'}
              fallback={
                <Show when={block._type === 'image'}>
                  <RenderImage block={block as ImageBlock} />
                </Show>
              }
            >
              <RenderBlock block={block as Block} />
            </Show>
          )}
        </For>
      </div>
    </Show>
  );
};

export default PortableText;
