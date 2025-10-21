import { ChatInput } from "@/components/Chat/ChatInput";
import { ChatMessages, Message } from "@/components/Chat/ChatMessages";
import { ChatSidebar } from "@/components/Chat/ChatSidebar";
import { TravelData, TravelResults } from "@/components/Chat/TravelResults";
import { getTokens, logout } from "@/lib/tanstack-query";
import { useRef, useState } from "react";

const API_BASE_URL =
    (import.meta as any).env?.VITE_API_URL ?? "http://0.0.0.0:8081/api/v1";

const ChatPage = () => {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            role: "bot",
            content:
                "Hello! I'm your personal Travel Planner AI. Where are you headed, or should I help you discover your next adventure? 🌍",
            timestamp: new Date(),
        },
    ]);

    const [chatHistory] = useState([
        { id: "1", title: "Trip to Coorg", date: "Today" },
        { id: "2", title: "Kerala Backwaters", date: "Yesterday" },
    ]);

    const [currentChatId, setCurrentChatId] = useState<string>("1");
    const [travelData, setTravelData] = useState<TravelData | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const abortControllerRef = useRef<AbortController | null>(null);

    const handleSendMessage = async (content: string) => {
        const userMessage: Message = {
            id: Date.now().toString(),
            role: "user",
            content,
            timestamp: new Date(),
        };

        setMessages((prev) => [...prev, userMessage]);
        setIsLoading(true);

        // Create a placeholder message for the bot response
        const botMessageId = (Date.now() + 1).toString();
        const botMessage: Message = {
            id: botMessageId,
            role: "bot",
            content: "",
            timestamp: new Date(),
            isStreaming: true, // Mark as streaming
        };

        setMessages((prev) => [...prev, botMessage]);

        // Create an AbortController for this request
        abortControllerRef.current = new AbortController();

        try {
            // Get authentication token from localStorage
            const { access } = getTokens();

            if (!access) {
                // No token - redirect to login
                await logout();
                throw new Error("Not authenticated. Redirecting to login...");
            }

            const response = await fetch(`${API_BASE_URL}/chats/`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${access}`,
                },
                body: JSON.stringify({
                    ChatID: currentChatId,
                    Content: content,
                }),
                signal: abortControllerRef.current.signal,
            });

            if (!response.ok) {
                if (response.status === 401) {
                    // Authentication failed - redirect to login
                    await logout();
                    throw new Error("Session expired. Redirecting to login...");
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            if (!response.body) {
                throw new Error("Response body is null");
            }

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let accumulatedContent = "";

            while (true) {
                const { done, value } = await reader.read();

                if (done) {
                    break;
                }

                const chunk = decoder.decode(value, { stream: true });
                const lines = chunk.split("\n");

                for (const line of lines) {
                    if (line.startsWith("data: ")) {
                        const data = line.slice(6); // Remove "data: " prefix

                        // Decode escaped newlines back to actual newlines
                        const decodedData = data.replace(/\\n/g, "\n");

                        accumulatedContent += decodedData;

                        // Update the bot message with accumulated content (still streaming)
                        setMessages((prev) =>
                            prev.map((msg) =>
                                msg.id === botMessageId
                                    ? {
                                          ...msg,
                                          content: accumulatedContent,
                                          isStreaming: true,
                                      }
                                    : msg
                            )
                        );
                    }
                }
            }

            // Streaming complete - mark as finished to trigger markdown rendering
            setMessages((prev) =>
                prev.map((msg) =>
                    msg.id === botMessageId
                        ? { ...msg, isStreaming: false }
                        : msg
                )
            );

            setIsLoading(false);
        } catch (error) {
            if (error instanceof Error) {
                if (error.name === "AbortError") {
                    console.log("Request was aborted");
                } else {
                    console.error("Error streaming message:", error);
                    // Update bot message with error
                    setMessages((prev) =>
                        prev.map((msg) =>
                            msg.id === botMessageId
                                ? {
                                      ...msg,
                                      content:
                                          error.message ||
                                          "Sorry, I encountered an error. Please try again.",
                                  }
                                : msg
                        )
                    );
                }
            }
            setIsLoading(false);
        }
    };

    const handleNewChat = () => {
        // Abort any ongoing request
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }

        setMessages([
            {
                id: Date.now().toString(),
                role: "bot",
                content:
                    "Hello! I'm your personal Travel Planner AI. Where are you headed, or should I help you discover your next adventure? 🌍",
                timestamp: new Date(),
            },
        ]);
        setTravelData(null);
        setCurrentChatId(Date.now().toString());
        setIsLoading(false);
    };

    const handleSelectChat = (id: string) => {
        setCurrentChatId(id);
        // In a real app, load chat history from storage
    };

    return (
        <div className="flex h-screen overflow-hidden">
            <ChatSidebar
                chatHistory={chatHistory}
                onNewChat={handleNewChat}
                onSelectChat={handleSelectChat}
                currentChatId={currentChatId}
            />

            <div className="flex-1 flex flex-col">
                <ChatMessages messages={messages} isLoading={isLoading} />
                <ChatInput onSend={handleSendMessage} disabled={isLoading} />
            </div>

            <TravelResults travelData={travelData} />
        </div>
    );
};

export default ChatPage;
