function escapeHtml(text: string): string {
  return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

/**
 * Renders a small, safe subset of Markdown (bold, italic, inline code,
 * "- " bullet lists, line breaks) to HTML. Input is HTML-escaped first, so
 * this is safe to feed into v-html even for untrusted/LLM-generated text -
 * no raw tag from the source can ever survive into the output.
 */
export function renderChatMarkdown(text: string): string {
  const escaped = escapeHtml(text)

  const withInline = escaped
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/(?<!\*)\*(?!\*)(.+?)\*(?!\*)/g, '<em>$1</em>')
    .replace(/`(.+?)`/g, '<code class="rounded bg-gray-100 px-1 py-0.5 text-xs">$1</code>')

  const lines = withInline.split('\n')
  const html: string[] = []
  let inList = false

  for (const line of lines) {
    const bullet = /^\s*[-*]\s+(.*)$/.exec(line)
    if (bullet) {
      if (!inList) {
        html.push('<ul class="list-disc pl-4">')
        inList = true
      }
      html.push(`<li>${bullet[1]}</li>`)
      continue
    }
    if (inList) {
      html.push('</ul>')
      inList = false
    }
    html.push(line.length > 0 ? `<p>${line}</p>` : '<br>')
  }
  if (inList) html.push('</ul>')

  return html.join('')
}
