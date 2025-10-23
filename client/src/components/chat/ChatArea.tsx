import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Send, User, Bot } from "lucide-react";
import { fetchClient } from "@/lib/tanstack-query";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface Message {
  id: string;
  type: "user" | "bot";
  content: string;
  timestamp: string;
  isStreaming?: boolean;
}

interface ChatAreaProps {
  chatId: string;
  userId: string;
}

export function ChatArea({ chatId, userId }: ChatAreaProps) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom
  const scrollToBottom = () => messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  useEffect(() => scrollToBottom(), [messages]);

  // Welcome message on new chat
  useEffect(() => {
    setMessages([
      {
        id: "welcome",
        type: "bot",
        content:
          "Hello! I'm your personal Travel Planner AI. Where are you headed, or should I help you discover your next adventure? 😊",
        timestamp: new Date().toISOString(),
        isStreaming: false,
      },
    ]);
    setIsTyping(false);
  }, [chatId]);

  const handleSend = async () => {
    if (!inputValue.trim() || isTyping) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: "user",
      content: inputValue,
      timestamp: new Date().toISOString(),
      isStreaming: false,
    };
    setMessages((prev) => [...prev, userMessage]);

    const messageContent = inputValue;
    setInputValue("");
    setIsTyping(true);

    const botMessageId = (Date.now() + 1).toString();
    const streamingMessage: Message = {
      id: botMessageId,
      type: "bot",
      content: "",
      timestamp: new Date().toISOString(),
      isStreaming: true,
    };
    setMessages((prev) => [...prev, streamingMessage]);

    try {
      const response = await fetchClient("/api/v1/chats/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content: messageContent, chatId, userId }),
      });

      if (!response.ok || !response.body) throw new Error("Failed to receive streaming response");

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let fullResponse = "";
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const content = line.slice(6);
            if (content && content !== "[DONE]") {
              fullResponse += content;

              // Update streaming content
              setMessages((prev) =>
                prev.map((msg) =>
                  msg.id === botMessageId
                    ? { ...msg, content: fullResponse, isStreaming: true }
                    : msg
                )
              );
            }
          }
        }
      }

      // Mark streaming complete
      setMessages((prev) =>
        prev.map((msg) =>
          msg.id === botMessageId
            ? { ...msg, content: fullResponse.replace(/\\n/g, "\n"), isStreaming: false }
            : msg
        )
      );

      setIsTyping(false);
    } catch (error) {
      console.error("Error streaming message:", error);
      setMessages((prev) =>
        prev.map((msg) =>
          msg.id === botMessageId ? { ...msg, content: "Error", isStreaming: false } : msg
        )
      );
      setIsTyping(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="flex flex-col h-full bg-white">
      {/* Chat Messages */}
      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        {messages.map((message) => (
          <div
            key={message.id}
            className={`flex ${message.type === "user" ? "justify-end" : "justify-start"}`}
          >
            <div
              className={`flex items-start space-x-3 max-w-[85%] ${
                message.type === "user" ? "flex-row-reverse space-x-reverse" : ""
              }`}
            >
              <div className="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 bg-black text-white">
                {message.type === "user" ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
              </div>

              <div
                className={`p-4 rounded-2xl ${
                  message.type === "user"
                    ? "bg-purple-500 text-white border border-purple-600"
                    : "bg-gray-50 text-gray-900 border border-gray-200"
                }`}
              >
                {message.type === "bot" ? (
                  message.isStreaming ? (
                    <p className="text-sm whitespace-pre-wrap leading-relaxed">
                      {message.content}
                    </p>
                  ) : (
                    <ReactMarkdown
                      remarkPlugins={[remarkGfm]}
                      components={{
                        h1: ({ children }) => <h1 className="text-xl font-bold mb-3 mt-2">{children}</h1>,
                        h2: ({ children }) => <h2 className="text-lg font-bold mb-2 mt-4">{children}</h2>,
                        h3: ({ children }) => <h3 className="text-base font-semibold mb-2 mt-3">{children}</h3>,
                        p: ({ children }) => <p className="mb-2 leading-relaxed">{children}</p>,
                        ul: ({ children }) => <ul className="list-disc list-inside mb-2 space-y-1 ml-2">{children}</ul>,
                        ol: ({ children }) => <ol className="list-decimal list-inside mb-2 space-y-1 ml-2">{children}</ol>,
                        li: ({ children }) => <li className="leading-relaxed">{children}</li>,
                        table: ({ children }) => <div className="overflow-x-auto my-3"><table className="min-w-full divide-y divide-border border border-border rounded-md text-xs">{children}</table></div>,
                        thead: ({ children }) => <thead className="bg-muted/50">{children}</thead>,
                        tbody: ({ children }) => <tbody className="divide-y divide-border bg-card">{children}</tbody>,
                        th: ({ children }) => <th className="px-3 py-2 text-left text-xs font-medium text-foreground">{children}</th>,
                        td: ({ children }) => <td className="px-3 py-2 text-xs text-foreground">{children}</td>,
                        tr: ({ children }) => <tr className="hover:bg-muted/30">{children}</tr>,
                        blockquote: ({ children }) => <blockquote className="border-l-4 border-primary pl-3 italic my-2 text-muted-foreground">{children}</blockquote>,
                        pre: ({ children }) => <pre className="bg-muted p-3 rounded-md text-xs font-mono overflow-x-auto my-2 text-foreground">{children}</pre>,
                        code: ({ children, className }) => {
                          const isInline = !className;
                          return isInline ? (
                            <code className="bg-muted px-1.5 py-0.5 rounded text-xs font-mono text-foreground">{children}</code>
                          ) : (
                            <code className={className}>{children}</code>
                          );
                        },
                        strong: ({ children }) => <strong className="font-bold text-foreground">{children}</strong>,
                        em: ({ children }) => <em className="italic text-foreground">{children}</em>,
                        hr: () => <hr className="my-4 border-border" />,
                        a: ({ children, href }) => <a href={href} className="text-primary underline hover:text-primary/80" target="_blank" rel="noopener noreferrer">{children}</a>,
                      }}
                    >
                      {message.content}
                    </ReactMarkdown>
                  )
                ) : (
                  <p className="text-sm whitespace-pre-wrap leading-relaxed">{message.content}</p>
                )}
              </div>
            </div>
          </div>
        ))}

        <div ref={messagesEndRef} />
      </div>

      {/* Input Field */}
      <div className="p-6 border-t border-gray-100">
        <div className="flex space-x-3">
          <Input
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Type your message..."
            className="flex-1 border border-gray-200 rounded-2xl px-4 py-3"
            disabled={isTyping}
          />
          <Button
            onClick={handleSend}
            disabled={!inputValue.trim() || isTyping}
            className="rounded-2xl px-6"
          >
            <Send className="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
