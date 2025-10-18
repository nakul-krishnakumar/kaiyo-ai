import { useState } from 'react';
import { ChatSidebar } from '@/components/Chat/ChatSidebar';
import { ChatMessages, Message } from '@/components/Chat/ChatMessages';
import { ChatInput } from '@/components/Chat/ChatInput';
import { TravelResults, TravelData } from '@/components/Chat/TravelResults';

const ChatPage = () => {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: '1',
      role: 'bot',
      content: "Hello! I'm your personal Travel Planner AI. Where are you headed, or should I help you discover your next adventure? 🌍",
      timestamp: new Date(),
    },
  ]);

  const [chatHistory, setChatHistory] = useState([
    { id: '1', title: 'Trip to Coorg', date: 'Today' },
    { id: '2', title: 'Kerala Backwaters', date: 'Yesterday' },
  ]);

  const [currentChatId, setCurrentChatId] = useState<string>('1');
  const [travelData, setTravelData] = useState<TravelData | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSendMessage = async (content: string) => {
    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setIsLoading(true);

    // Simulate API call to backend
    setTimeout(() => {
      const botResponse: Message = {
        id: (Date.now() + 1).toString(),
        role: 'bot',
        content: generateMockResponse(content),
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, botResponse]);
      
      // If the message is about a trip, generate mock travel data
      if (content.toLowerCase().includes('coorg') || content.toLowerCase().includes('trip')) {
        setTravelData({
          destination: '3-Day Itinerary in Coorg',
          totalCost: '10000Rs',
          dates: 'Feb 12 - Feb 19, 2025',
          coordinates: [12.4244, 75.7382],
          itinerary: [
            {
              day: 1,
              title: 'Arrival and Relaxation',
              activities: [
                'Arrive at Madikeri and check into homestay',
                'Sunset at Raja\'s Seat',
                'Explore local markets',
              ],
            },
            {
              day: 2,
              title: 'Nature and Adventure',
              activities: [
                'Morning trek to Tadiandamol peak',
                'Local Coorgi lunch at spice plantation',
                'Visit Abbey Falls',
              ],
            },
            {
              day: 3,
              title: 'Culture and Departure',
              activities: [
                'Visit Namdroling Monastery',
                'Coffee plantation tour',
                'Departure from Madikeri',
              ],
            },
          ],
        });
      }
      
      setIsLoading(false);
    }, 1000);
  };

  const generateMockResponse = (userMessage: string): string => {
    if (userMessage.toLowerCase().includes('coorg')) {
      return "Absolutely! 🌿 For a peaceful and scenic getaway, here are a few handpicked options:\n\n• Coorg, Karnataka – lush coffee plantations and misty hills\n• Alleppey, Kerala – serene backwaters and houseboats\n• Tawang, Arunachal Pradesh – calm monasteries and Himalayan views\n\nWant me to check best travel dates or stay options in one of these?";
    }
    if (userMessage.toLowerCase().includes('3-day')) {
      return "Perfect choice! Here's a sample 3-day plan for Coorg:\n\n🏔️ Day 1: Arrive and relax at a homestay in Madikeri. Sunset at Raja's Seat.\n🥾 Day 2: Morning trek to Tadiandamol peak, local Coorgi lunch, spice plantation tour\n🏛️ Day 3: Visit Abbey Falls and Namdroling Monastery before departure\n\nShall I prepare a detailed itinerary with estimated costs?";
    }
    return "I'd be happy to help you plan that! Could you tell me more about your preferences, budget, or the type of experience you're looking for?";
  };

  const handleNewChat = () => {
    setMessages([
      {
        id: Date.now().toString(),
        role: 'bot',
        content: "Hello! I'm your personal Travel Planner AI. Where are you headed, or should I help you discover your next adventure? 🌍",
        timestamp: new Date(),
      },
    ]);
    setTravelData(null);
    setCurrentChatId(Date.now().toString());
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
        <ChatMessages messages={messages} />
        <ChatInput onSend={handleSendMessage} disabled={isLoading} />
      </div>

      <TravelResults travelData={travelData} />
    </div>
  );
};

export default ChatPage;
