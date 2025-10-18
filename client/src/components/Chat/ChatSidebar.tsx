import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Plus, MessageSquare, LogOut } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

interface ChatHistory {
  id: string;
  title: string;
  date: string;
}

interface ChatSidebarProps {
  chatHistory: ChatHistory[];
  onNewChat: () => void;
  onSelectChat: (id: string) => void;
  currentChatId?: string;
}

export const ChatSidebar = ({
  chatHistory,
  onNewChat,
  onSelectChat,
  currentChatId,
}: ChatSidebarProps) => {
  const navigate = useNavigate();
  const { logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  return (
    <div className="w-64 bg-card border-r border-border flex flex-col h-screen">
      <div className="p-4 border-b border-border">
        <h2 className="text-xl font-bold bg-gradient-hero bg-clip-text text-transparent">
          Kaiyo AI
        </h2>
      </div>

      <div className="p-4">
        <Button onClick={onNewChat} variant="hero" className="w-full" size="lg">
          <Plus className="mr-2 h-4 w-4" />
          New Chat
        </Button>
      </div>

      <ScrollArea className="flex-1 px-4">
        <div className="space-y-2">
          {chatHistory.map((chat) => (
            <button
              key={chat.id}
              onClick={() => onSelectChat(chat.id)}
              className={`w-full text-left p-3 rounded-lg transition-colors ${
                currentChatId === chat.id
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary hover:bg-secondary/80"
              }`}
            >
              <div className="flex items-center gap-2 mb-1">
                <MessageSquare className="h-4 w-4" />
                <p className="font-medium text-sm truncate">{chat.title}</p>
              </div>
              <p className="text-xs opacity-70">{chat.date}</p>
            </button>
          ))}
        </div>
      </ScrollArea>

      <div className="p-4 border-t border-border">
        <Button onClick={handleLogout} variant="ghost" className="w-full">
          <LogOut className="mr-2 h-4 w-4" />
          Logout
        </Button>
      </div>
    </div>
  );
};
