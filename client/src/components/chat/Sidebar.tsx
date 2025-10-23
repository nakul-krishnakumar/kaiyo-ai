import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Plus, MessageSquare, X, PanelLeftClose, PanelLeft, LogOut } from "lucide-react";
import { logout } from "@/lib/tanstack-query";

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

  const handleNewChat = () => {
    onNewChat();
  };

  const handleLogout = async () => {
    await logout();
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

      {/* Sidebar - Full width when open, icon-only when closed */}
      <div className={`h-full bg-gray-50 flex flex-col transition-all duration-300 ${isOpen ? 'w-full' : 'w-20'}`}>
        {isOpen ? (
          // Expanded Sidebar
          <>
            <div className="py-6 px-4">
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

              <Button onClick={handleNewChat} className="w-full mb-6 comic-button bg-purple-600 hover:bg-purple-700">
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

            <div className="p-6 border-t border-gray-200">
              <Button 
                onClick={handleLogout} 
                variant="outline" 
                className="w-full text-purple-600 hover:bg-purple-200 hover:text-purple-800 comic-button"
              >
                <LogOut className="w-4 h-4 mr-2" />
                Logout
              </Button>
            </div>
          </>
        ) : (
          // Collapsed Sidebar - Icon only
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

            <div className="border-t border-gray-300 w-12" />

            <Button
              variant="ghost"
              size="icon"
              onClick={handleNewChat}
              className="hover:bg-gray-200"
              title="New Chat"
            >
              <Plus className="w-5 h-5" />
            </Button>

            <div className="flex-1 flex flex-col items-center space-y-3 pt-4">
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
                  <MessageSquare className="w-5 h-5" />
                </Button>
              ))}
            </div>

            <div className="border-t border-gray-300 w-12 mt-auto" />

            <Button
              variant="ghost"
              size="icon"
              onClick={handleLogout}
              className="hover:bg-purple-200"
              title="Logout"
            >
              <LogOut className="w-5 h-5" />
            </Button>
          </div>
        )}
      </div>
    </>
  );
}
