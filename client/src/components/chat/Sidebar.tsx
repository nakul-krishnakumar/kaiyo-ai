import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Plus, MessageSquare, X } from "lucide-react";

interface Chat {
  id: string;
  title: string;
  timestamp: string;
}

interface SidebarProps {
  isOpen: boolean;
  onToggle: () => void;
  onNewChat: () => void;
  onSelectChat: (chatId: string) => void;
  currentChatId?: string;
}

export function Sidebar({
  isOpen,
  onToggle,
  onNewChat,
  onSelectChat,
  currentChatId,
}: SidebarProps) {
  const [chats] = useState<Chat[]>([
    { id: "1", title: "3-Day Itinerary in Coorg", timestamp: "Feb 10, 2025" },
    { id: "2", title: "Weekend in Goa", timestamp: "Feb 8, 2025" },
    { id: "3", title: "Kerala Backwaters Trip", timestamp: "Feb 5, 2025" },
  ]);

  return (
    <>
      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={onToggle}
        />
      )}

      {/* Sidebar */}
      <div className="h-full bg-gray-50 flex flex-col">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold text-black">Kaiyo AI</h2>
            <Button
              variant="ghost"
              size="icon"
              onClick={onToggle}
              className="lg:hidden"
            >
              <X className="w-5 h-5" />
            </Button>
          </div>

          <Button onClick={onNewChat} className="w-full mb-6 comic-button">
            <Plus className="w-4 h-4 mr-2" />
            New Chat
          </Button>
        </div>

        <div className="flex-1 px-6 pb-6">
          <h3 className="text-sm font-medium text-gray-500 mb-4">
            Previous chats
          </h3>
          <div className="space-y-2">
            {chats.map((chat) => (
              <div
                key={chat.id}
                className={`p-3 rounded-lg cursor-pointer hover:bg-gray-100 transition-colors ${
                  currentChatId === chat.id ? "bg-gray-100" : ""
                }`}
                onClick={() => onSelectChat(chat.id)}
              >
                <div className="flex items-start space-x-3">
                  <MessageSquare className="w-4 h-4 mt-1 text-gray-400" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      {chat.title}
                    </p>
                    <p className="text-xs text-gray-500">{chat.timestamp}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </>
  );
}
