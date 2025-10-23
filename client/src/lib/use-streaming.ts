// useStreamingMarkdown.ts
import { useState, useEffect, useRef } from "react";

function findSafeBoundary(text: string): number {
  // 1. If there's a closing code fence, split after the trailing newline
  const fenceEnd = text.lastIndexOf("```");
  if (fenceEnd !== -1) {
    const afterFence = text.indexOf("\n", fenceEnd + 3);
    if (afterFence !== -1) {
      return afterFence + 1;
    }
  }
  // 2. Otherwise split on last double-newline
  const dblNl = text.lastIndexOf("\n\n");
  if (dblNl !== -1) {
    return dblNl + 2;
  }
  // 3. No safe split
  return 0;
}

export function useStreamingMarkdown(stream: ReadableStream<string>) {
  const [renderText, setRenderText] = useState("");
  const bufferRef = useRef("");

  useEffect(() => {
    const reader = stream.getReader();
    let cancelled = false;

    async function read() {
      while (!cancelled) {
        const { value, done } = await reader.read();
        if (done) break;
        if (value) {
          bufferRef.current += value;
          const safeLen = findSafeBoundary(bufferRef.current);
          if (safeLen > 0) {
            setRenderText((prev) => prev + bufferRef.current.slice(0, safeLen));
            bufferRef.current = bufferRef.current.slice(safeLen);
          }
        }
      }
      // Flush remaining buffer at end
      if (!cancelled && bufferRef.current) {
        setRenderText((prev) => prev + bufferRef.current);
      }
    }

    read().catch(console.error);
    return () => {
      cancelled = true;
      reader.cancel();
    };
  }, [stream]);

  return renderText;
}