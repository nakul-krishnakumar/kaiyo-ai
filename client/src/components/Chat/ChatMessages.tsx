// ChatMessages.tsx
import { ScrollArea } from "@/components/ui/scroll-area";
import { motion } from "framer-motion";
import { Bot, User } from "lucide-react";
import { useEffect, useRef } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

export interface Message {
    id: string;
    role: "user" | "bot";
    content: string;
    timestamp: Date;
    isStreaming?: boolean; // Track if message is still streaming
}

interface ChatMessagesProps {
    messages: Message[];
    isLoading?: boolean;
}

export const ChatMessages = ({ messages, isLoading }: ChatMessagesProps) => {
    const scrollRef = useRef<HTMLDivElement>(null);

    // Auto-scroll to bottom when messages change
    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [messages]);

    return (
        <ScrollArea className="flex-1 p-6">
            <div className="space-y-4 max-w-4xl mx-auto">
                {messages.map((message, index) => (
                    <motion.div
                        key={message.id}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className={`flex gap-3 ${
                            message.role === "user"
                                ? "justify-end"
                                : "justify-start"
                        }`}
                    >
                        {message.role === "bot" && (
                            <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center flex-shrink-0">
                                <Bot className="h-5 w-5 text-primary-foreground" />
                            </div>
                        )}
                        <div
                            className={`rounded-2xl px-4 py-3 max-w-[80%] shadow-card ${
                                message.role === "user"
                                    ? "bg-gradient-hero text-primary-foreground"
                                    : "bg-card"
                            }`}
                        >
                            {message.role === "bot" ? (
                                message.content === "" && isLoading ? (
                                    // Loading animation
                                    <span className="inline-flex gap-1">
                                        <span
                                            className="w-2 h-2 bg-primary rounded-full animate-bounce"
                                            style={{ animationDelay: "0ms" }}
                                        />
                                        <span
                                            className="w-2 h-2 bg-primary rounded-full animate-bounce"
                                            style={{ animationDelay: "150ms" }}
                                        />
                                        <span
                                            className="w-2 h-2 bg-primary rounded-full animate-bounce"
                                            style={{ animationDelay: "300ms" }}
                                        />
                                    </span>
                                ) : message.isStreaming ? (
                                    // Show plain text while streaming
                                    <p className="text-sm whitespace-pre-wrap text-foreground">
                                        {message.content}
                                    </p>
                                ) : (
                                    // Render markdown only when streaming is complete
                                    <div className="text-sm text-foreground prose prose-sm dark:prose-invert max-w-none">
                                        <ReactMarkdown
                                            remarkPlugins={[remarkGfm]}
                                            components={{
                                                h1: ({ children }) => (
                                                    <h1 className="text-xl font-bold mb-3 mt-2 text-foreground">
                                                        {children}
                                                    </h1>
                                                ),
                                                h2: ({ children }) => (
                                                    <h2 className="text-lg font-bold mb-2 mt-4 text-foreground">
                                                        {children}
                                                    </h2>
                                                ),
                                                h3: ({ children }) => (
                                                    <h3 className="text-base font-semibold mb-2 mt-3 text-foreground">
                                                        {children}
                                                    </h3>
                                                ),
                                                p: ({ children }) => (
                                                    <p className="mb-2 leading-relaxed text-foreground">
                                                        {children}
                                                    </p>
                                                ),
                                                ul: ({ children }) => (
                                                    <ul className="list-disc list-inside mb-2 space-y-1 ml-2 text-foreground">
                                                        {children}
                                                    </ul>
                                                ),
                                                ol: ({ children }) => (
                                                    <ol className="list-decimal list-inside mb-2 space-y-1 ml-2 text-foreground">
                                                        {children}
                                                    </ol>
                                                ),
                                                li: ({ children }) => (
                                                    <li className="leading-relaxed text-foreground">
                                                        {children}
                                                    </li>
                                                ),
                                                table: ({ children }) => (
                                                    <div className="overflow-x-auto my-3">
                                                        <table className="min-w-full divide-y divide-border border border-border rounded-md text-xs">
                                                            {children}
                                                        </table>
                                                    </div>
                                                ),
                                                thead: ({ children }) => (
                                                    <thead className="bg-muted/50">
                                                        {children}
                                                    </thead>
                                                ),
                                                tbody: ({ children }) => (
                                                    <tbody className="divide-y divide-border bg-card">
                                                        {children}
                                                    </tbody>
                                                ),
                                                th: ({ children }) => (
                                                    <th className="px-3 py-2 text-left text-xs font-medium text-foreground">
                                                        {children}
                                                    </th>
                                                ),
                                                td: ({ children }) => (
                                                    <td className="px-3 py-2 text-xs text-foreground">
                                                        {children}
                                                    </td>
                                                ),
                                                tr: ({ children }) => (
                                                    <tr className="hover:bg-muted/30">
                                                        {children}
                                                    </tr>
                                                ),
                                                blockquote: ({ children }) => (
                                                    <blockquote className="border-l-4 border-primary pl-3 italic my-2 text-muted-foreground">
                                                        {children}
                                                    </blockquote>
                                                ),
                                                pre: ({ children }) => (
                                                    <pre className="bg-muted p-3 rounded-md text-xs font-mono overflow-x-auto my-2 text-foreground">
                                                        {children}
                                                    </pre>
                                                ),
                                                code: ({ children, className }) => {
                                                    const isInline = !className;
                                                    return isInline ? (
                                                        <code className="bg-muted px-1.5 py-0.5 rounded text-xs font-mono text-foreground">
                                                            {children}
                                                        </code>
                                                    ) : (
                                                        <code className={className}>
                                                            {children}
                                                        </code>
                                                    );
                                                },
                                                strong: ({ children }) => (
                                                    <strong className="font-bold text-foreground">
                                                        {children}
                                                    </strong>
                                                ),
                                                em: ({ children }) => (
                                                    <em className="italic text-foreground">
                                                        {children}
                                                    </em>
                                                ),
                                                hr: () => (
                                                    <hr className="my-4 border-border" />
                                                ),
                                                a: ({ children, href }) => (
                                                    <a
                                                        href={href}
                                                        className="text-primary underline hover:text-primary/80"
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                    >
                                                        {children}
                                                    </a>
                                                ),
                                            }}
                                        >
                                            {message.content}
                                        </ReactMarkdown>
                                    </div>
                                )
                            ) : (
                                <p className="text-sm whitespace-pre-wrap">
                                    {message.content}
                                </p>
                            )}
                        </div>
                        {message.role === "user" && (
                            <div className="w-8 h-8 rounded-full bg-secondary flex items-center justify-center flex-shrink-0">
                                <User className="h-5 w-5" />
                            </div>
                        )}
                    </motion.div>
                ))}
                <div ref={scrollRef} />
            </div>
        </ScrollArea>
    );
};
