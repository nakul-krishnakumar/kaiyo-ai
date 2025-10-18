import { motion } from 'framer-motion';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Bot, User } from 'lucide-react';

export interface Message {
  id: string;
  role: 'user' | 'bot';
  content: string;
  timestamp: Date;
}

interface ChatMessagesProps {
  messages: Message[];
}

export const ChatMessages = ({ messages }: ChatMessagesProps) => {
  return (
    <ScrollArea className="flex-1 p-6">
      <div className="space-y-4 max-w-4xl mx-auto">
        {messages.map((message, index) => (
          <motion.div
            key={message.id}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
            className={`flex gap-3 ${message.role === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            {message.role === 'bot' && (
              <div className="w-8 h-8 rounded-full bg-primary flex items-center justify-center flex-shrink-0">
                <Bot className="h-5 w-5 text-primary-foreground" />
              </div>
            )}
            <div
              className={`rounded-2xl px-4 py-3 max-w-[80%] shadow-card ${
                message.role === 'user'
                  ? 'bg-gradient-hero text-primary-foreground'
                  : 'bg-card'
              }`}
            >
              <p className="text-sm whitespace-pre-wrap">{message.content}</p>
            </div>
            {message.role === 'user' && (
              <div className="w-8 h-8 rounded-full bg-secondary flex items-center justify-center flex-shrink-0">
                <User className="h-5 w-5" />
              </div>
            )}
          </motion.div>
        ))}
      </div>
    </ScrollArea>
  );
};
