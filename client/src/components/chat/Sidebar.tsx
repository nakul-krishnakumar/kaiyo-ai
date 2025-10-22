import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Plus, MessageSquare, X, PanelLeftClose, PanelLeft, LogOut } from "lucide-react";
import { logout } from "@/lib/tanstack-query";

interface Chat {
  id: string;
  title: string;
  timestamp: string;
}
import { useNavigate } from "@tanstack/react-router";
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

  const handleNewChat = () => {
    onNewChat();
  };
const navigate = useNavigate();
    const handleLogout = () => {
        localStorage.clear();
        sessionStorage.clear();
        navigate({ to: "/login" });
    };

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
        {isOpen ? (
          // Full sidebar
          <>
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-2xl font-bold text-black">Kaiyo AI</h2>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={onToggle}
                  className="hover:bg-gray-200"
                  title="Collapse sidebar"
                >
                  <PanelLeftClose className="w-5 h-5" />
                </Button>
              </div>

              <Button onClick={handleNewChat} className="w-full mb-6 comic-button bg-purple-600 hover:bg-purple-500">
                <Plus className="w-4 h-4 mr-2" />
                New Chat
              </Button>
            </div>

            <div className="flex-1 px-6 pb-6 overflow-y-auto">
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

            {/* Logout button */}
            <div className="p-6 border-t border-gray-200">
              <Button
                onClick={handleLogout}
                variant="outline"
                className="w-full flex items-center justify-center space-x-2 hover:bg-red-50 hover:text-red-600 hover:border-red-300"
              >
                <LogOut className="w-4 h-4" />
                <span>Logout</span>
              </Button>
            </div>
          </>
        ) : (
          // Collapsed sidebar - icons only
          <div className="flex flex-col items-center py-6 space-y-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={onToggle}
              className="hover:bg-gray-200"
              title="Expand sidebar"
            >
              <PanelLeft className="w-5 h-5" />
            </Button>

            <Button
              onClick={handleNewChat}
              size="icon"
              className="comic-button"
              title="New Chat"
            >
              <Plus className="w-4 h-4" />
            </Button>

            <div className="flex-1 flex flex-col space-y-2 mt-4">
              {chats.map((chat) => (
                <Button
                  key={chat.id}
                  variant="ghost"
                  size="icon"
                  onClick={() => onSelectChat(chat.id)}
                  className={`hover:bg-gray-200 ${
                    currentChatId === chat.id ? "bg-gray-200" : ""
                  }`}
                  title={chat.title}
                >
                  <MessageSquare className="w-4 h-4" />
                </Button>
              ))}
            </div>

            {/* Logout button for collapsed view */}
            <div className="mt-auto pt-4 border-t border-gray-200">
              <Button
                onClick={handleLogout}
                variant="ghost"
                size="icon"
                className="cursor-pointer text-destructive hover:bg-red-50 hover:text-red-600"
                title="Logout"
              >
                <LogOut className="w-4 h-4" />
              </Button>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
