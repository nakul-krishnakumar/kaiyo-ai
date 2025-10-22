import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";
import { Sidebar } from "@/components/chat/Sidebar";
import { ChatArea } from "@/components/chat/ChatArea";
import { MapPanel } from "@/components/chat/MapPanel";
export default function ChatPage() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [currentChatId, setCurrentChatId] = useState("1");
  const [locations, setLocations] = useState([
    { name: "Coorg", lat: 12.3375, lng: 75.8069, type: "destination" },
    { name: "Raja's Seat", lat: 12.4244, lng: 75.7382, type: "attraction" },
    { name: "Tadiandamol Peak", lat: 12.2458, lng: 75.7167, type: "trekking" },
    { name: "Abbey Falls", lat: 12.4544, lng: 75.7167, type: "waterfall" },
  ]);

  const handleNewChat = () => {
    const newChatId = Date.now().toString();
    setCurrentChatId(newChatId);
    setSidebarOpen(false);
  };

  const handleSelectChat = (chatId: string) => {
    setCurrentChatId(chatId);
    setSidebarOpen(false);
  };

  const handleLocationUpdate = (newLocations: any[]) => {
    setLocations(newLocations);
  };

  return (
    <div className="chat-container flex bg-gray-50">
      {/* Sidebar - Desktop */}
      <div className={`sidebar-panel ${sidebarOpen ? "block" : ""} lg:block`}>
        <Sidebar
          isOpen={sidebarOpen}
          onToggle={() => setSidebarOpen(!sidebarOpen)}
          onNewChat={handleNewChat}
          onSelectChat={handleSelectChat}
          currentChatId={currentChatId}
        />
      </div>

      {/* Main Content */}
      <div className="flex-1 flex">
        {/* Chat Area */}
        <div className="flex-1 flex flex-col min-w-0">
          {/* Header - only show on mobile */}
          <div className="lg:hidden bg-white border-b p-4 flex items-center justify-between">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setSidebarOpen(true)}
            >
              <Menu className="w-5 h-5" />
            </Button>
            <h1 className="text-lg font-semibold">Kaiyo AI</h1>
            <div></div>
          </div>

          {/* Chat Messages - Full height on mobile, shared on desktop */}
          <div className="flex-1 lg:h-screen">
            <ChatArea onLocationUpdate={handleLocationUpdate} />
          </div>
        </div>

        {/* Map Panel - Desktop only */}
        <div className="map-panel hidden lg:flex lg:flex-col">
          <MapPanel locations={locations} />
        </div>
      </div>

      {/* Mobile Map Panel - Bottom Sheet */}
      <div className="lg:hidden fixed bottom-0 left-0 right-0 bg-white border-t-2 border-black rounded-t-2xl z-40">
        <div className="h-80 overflow-hidden">
          <div className="p-4 border-b border-gray-200">
            <div className="w-12 h-1 bg-gray-300 rounded-full mx-auto mb-2"></div>
            <h3 className="text-center font-semibold">Itinerary</h3>
          </div>
          <div className="h-64">
            <MapPanel locations={locations} />
          </div>
        </div>
      </div>
    </div>
  );
}
