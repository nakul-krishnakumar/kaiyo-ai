import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Send, User, Bot } from "lucide-react";

interface Message {
  id: string;
  type: "user" | "bot";
  content: string;
  timestamp: string;
}

interface ChatAreaProps {
  onLocationUpdate?: (locations: any[]) => void;
}

export function ChatArea({ onLocationUpdate }: ChatAreaProps) {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: "1",
      type: "bot",
      content:
        "Hello! I'm your personal Travel Planner AI. Where are you headed, or should I help you discover your next adventure? 😊",
      timestamp: new Date().toISOString(),
    },
    {
      id: "2",
      type: "user",
      content:
        "Hey! I want to go on a relaxing trip somewhere in India. Any suggestions?",
      timestamp: new Date().toISOString(),
    },
    {
      id: "3",
      type: "bot",
      content:
        "Absolutely! 😊 For a peaceful and scenic getaway, here are a few handpicked options:\n\n• Coorg, Karnataka – lush coffee plantations and misty hills\n• Alleppey, Kerala – serene backwaters and houseboats\n• Wayanad, Kerala – calm monasteries and Himalayan views\n\nWant me to check best travel dates or stay options in Alleppey?",
      timestamp: new Date().toISOString(),
    },
    {
      id: "4",
      type: "user",
      content:
        "Coorg sounds good! I want a 3-day plan with nature, light trekking, and good food.",
      timestamp: new Date().toISOString(),
    },
    {
      id: "5",
      type: "bot",
      content:
        "Perfect choice! Here's a sample 3-day plan for Coorg:\n\n🌟 Day 1: Arrive and relax at a homestay. Sunset at Raja's Seat.\n🌟 Day 2: Morning trek to Tadiandamol peak, local Coorg lunch, spice plantation tour\n🌟 Day 3: Visit Abbey Falls and Namdroling Monastery before departure\n\nShall I prepare a detailed itinerary with estimated costs?",
      timestamp: new Date().toISOString(),
    },
  ]);

  const [inputValue, setInputValue] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSend = async () => {
    if (!inputValue.trim()) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: "user",
      content: inputValue,
      timestamp: new Date().toISOString(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputValue("");
    setIsTyping(true);

    // Simulate AI response with SSE
    setTimeout(() => {
      const botMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "bot",
        content:
          "Of course! Here's what I found for you:\n\n🏨 Top-rated homestays (under ₹1500/night):\n• Green Woods Stay – near Madikeri town, traditional meals included\n• River Mist Home – scenic river view, WiFi, peaceful location\n\n🚌 Bus options from Bangalore:\n• KSRTC Volvo departs 10:00 PM, arrives 6:00 AM\n• Private Volvo (₹950) – more flexible timings\n\nWould you like me to reserve a homestay and alert you when bus bookings open?",
        timestamp: new Date().toISOString(),
      };

      setMessages((prev) => [...prev, botMessage]);
      setIsTyping(false);

      // Update map with Coorg location
      if (onLocationUpdate) {
        onLocationUpdate([
          { name: "Coorg", lat: 12.3375, lng: 75.8069, type: "destination" },
          {
            name: "Raja's Seat",
            lat: 12.4244,
            lng: 75.7382,
            type: "attraction",
          },
          {
            name: "Tadiandamol Peak",
            lat: 12.2458,
            lng: 75.7167,
            type: "trekking",
          },
          {
            name: "Abbey Falls",
            lat: 12.4544,
            lng: 75.7167,
            type: "waterfall",
          },
        ]);
      }
    }, 2000);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="flex flex-col h-full bg-white">
      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        {messages.map((message) => (
          <div
            key={message.id}
            className={`flex ${message.type === "user" ? "justify-end" : "justify-start"}`}
          >
            <div
              className={`flex items-start space-x-3 max-w-[85%] ${
                message.type === "user"
                  ? "flex-row-reverse space-x-reverse"
                  : ""
              }`}
            >
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 ${
                  message.type === "user"
                    ? "bg-black text-white"
                    : "bg-black text-white"
                }`}
              >
                {message.type === "user" ? (
                  <User className="w-4 h-4" />
                ) : (
                  <Bot className="w-4 h-4" />
                )}
              </div>

              <div
                className={`p-4 rounded-2xl ${
                  message.type === "user"
                    ? "bg-purple-500 text-white comic-border border-purple-600"
                    : "bg-gray-50 text-gray-900 comic-border border-gray-200"
                }`}
              >
                <p className="text-sm whitespace-pre-wrap leading-relaxed">
                  {message.content}
                </p>
              </div>
            </div>
          </div>
        ))}

        {isTyping && (
          <div className="flex justify-start">
            <div className="flex items-start space-x-3">
              <div className="w-8 h-8 rounded-full bg-black flex items-center justify-center">
                <Bot className="w-4 h-4 text-white" />
              </div>
              <div className="p-4 rounded-2xl bg-gray-50 comic-border border-gray-200">
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                  <div
                    className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"
                    style={{ animationDelay: "0.1s" }}
                  ></div>
                  <div
                    className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"
                    style={{ animationDelay: "0.2s" }}
                  ></div>
                </div>
              </div>
            </div>
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div className="p-6 border-t border-gray-100">
        <div className="flex space-x-3">
          <Input
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type your message..."
            className="flex-1 comic-border border-gray-200 rounded-2xl px-4 py-3"
          />
          <Button
            onClick={handleSend}
            disabled={!inputValue.trim()}
            className="comic-button rounded-2xl px-6"
          >
            <Send className="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
