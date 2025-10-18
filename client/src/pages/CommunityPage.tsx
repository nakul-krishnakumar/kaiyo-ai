import { motion } from "framer-motion";
import { useRouter } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import {
  ArrowLeft,
  MessageSquare,
  PenTool,
  TrendingUp,
  Heart,
  MessageCircle,
  Eye,
  Plus,
} from "lucide-react";
import { logout } from "@/lib/tanstack-query";

const CommunityPage = () => {
  const router = useRouter();

  const questions = [
    {
      title: "Best time to visit Japan for cherry blossoms?",
      author: "Sarah M.",
      replies: 23,
      views: 456,
      likes: 12,
      tags: ["Japan", "Spring", "Nature"],
    },
    {
      title: "Budget-friendly hotels in Paris?",
      author: "John D.",
      replies: 18,
      views: 342,
      likes: 8,
      tags: ["Paris", "Budget", "Hotels"],
    },
    {
      title: "Safety tips for solo female travelers in Bali",
      author: "Maria L.",
      replies: 31,
      views: 678,
      likes: 45,
      tags: ["Bali", "Safety", "Solo Travel"],
    },
    {
      title: "Visa requirements for traveling to Thailand",
      author: "Alex K.",
      replies: 15,
      views: 290,
      likes: 6,
      tags: ["Thailand", "Visa", "Documentation"],
    },
  ];

  const blogs = [
    {
      title: "10 Hidden Gems in Southeast Asia",
      author: "Travel Explorer",
      excerpt:
        "Discover lesser-known destinations that will take your breath away...",
      likes: 156,
      comments: 42,
      readTime: "8 min read",
      image: "🏝️",
    },
    {
      title: "My 3-Month European Backpacking Journey",
      author: "Nomad Sarah",
      excerpt:
        "From the streets of Barcelona to the fjords of Norway, here's everything I learned...",
      likes: 203,
      comments: 67,
      readTime: "12 min read",
      image: "🎒",
    },
    {
      title: "Foodie Guide: Street Food in Vietnam",
      author: "Chef Wanderer",
      excerpt:
        "The ultimate guide to experiencing authentic Vietnamese cuisine...",
      likes: 189,
      comments: 54,
      readTime: "6 min read",
      image: "🍜",
    },
  ];

  const trending = [
    { name: "Tokyo, Japan", travelers: 1240, trend: "+15%", image: "🗼" },
    { name: "Paris, France", travelers: 980, trend: "+8%", image: "🗼" },
    { name: "Bali, Indonesia", travelers: 856, trend: "+22%", image: "🏝️" },
    { name: "Dubai, UAE", travelers: 742, trend: "+12%", image: "🏙️" },
  ];

  const handleLogout = () => {
    logout();
    router.navigate({ to: "/" });
  };

  return (
    <div className="min-h-screen p-4 md:p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="max-w-7xl mx-auto"
      >
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between mb-8 gap-4">
          <div>
            <h1 className="text-3xl md:text-4xl font-bold mb-2 bg-gradient-hero bg-clip-text text-transparent">
              Kaiyo Community
            </h1>
            <p className="text-muted-foreground">
              Share experiences, ask questions, and inspire fellow travelers
            </p>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => router.navigate({ to: "/chat" })}
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Chat
            </Button>
            <Button variant="ghost" onClick={handleLogout}>
              Logout
            </Button>
          </div>
        </div>

        <Tabs defaultValue="questions" className="w-full">
          <TabsList className="grid w-full md:w-auto grid-cols-3 mb-6">
            <TabsTrigger value="questions">
              <MessageSquare className="mr-2 h-4 w-4" />
              Questions
            </TabsTrigger>
            <TabsTrigger value="blogs">
              <PenTool className="mr-2 h-4 w-4" />
              Blogs
            </TabsTrigger>
            <TabsTrigger value="trending">
              <TrendingUp className="mr-2 h-4 w-4" />
              Trending
            </TabsTrigger>
          </TabsList>

          <TabsContent value="questions" className="space-y-4">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">Community Questions</h2>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Ask Question
              </Button>
            </div>
            {questions.map((question, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <Card className="p-6 hover-scale cursor-pointer">
                  <h3 className="text-lg font-semibold mb-3">
                    {question.title}
                  </h3>
                  <div className="flex flex-wrap gap-2 mb-4">
                    {question.tags.map((tag, i) => (
                      <Badge key={i} variant="secondary">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                  <div className="flex items-center justify-between text-sm text-muted-foreground">
                    <span className="font-medium">by {question.author}</span>
                    <div className="flex items-center gap-4">
                      <span className="flex items-center gap-1">
                        <Heart className="h-4 w-4" />
                        {question.likes}
                      </span>
                      <span className="flex items-center gap-1">
                        <MessageCircle className="h-4 w-4" />
                        {question.replies}
                      </span>
                      <span className="flex items-center gap-1">
                        <Eye className="h-4 w-4" />
                        {question.views}
                      </span>
                    </div>
                  </div>
                </Card>
              </motion.div>
            ))}
          </TabsContent>

          <TabsContent value="blogs" className="space-y-4">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">Travel Blogs</h2>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Write Blog
              </Button>
            </div>
            {blogs.map((blog, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <Card className="p-6 hover-scale cursor-pointer">
                  <div className="flex gap-4">
                    <div className="text-5xl">{blog.image}</div>
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold mb-2">
                        {blog.title}
                      </h3>
                      <p className="text-muted-foreground mb-3">
                        {blog.excerpt}
                      </p>
                      <div className="flex items-center justify-between text-sm">
                        <div className="flex items-center gap-4 text-muted-foreground">
                          <span className="font-medium">by {blog.author}</span>
                          <span>{blog.readTime}</span>
                        </div>
                        <div className="flex items-center gap-4 text-muted-foreground">
                          <span className="flex items-center gap-1">
                            <Heart className="h-4 w-4" />
                            {blog.likes}
                          </span>
                          <span className="flex items-center gap-1">
                            <MessageCircle className="h-4 w-4" />
                            {blog.comments}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </Card>
              </motion.div>
            ))}
          </TabsContent>

          <TabsContent value="trending" className="space-y-4">
            <h2 className="text-xl font-semibold mb-4">
              Trending Destinations
            </h2>
            <div className="grid md:grid-cols-2 gap-4">
              {trending.map((dest, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, scale: 0.95 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <Card className="p-6 hover-scale cursor-pointer">
                    <div className="flex items-center gap-4">
                      <span className="text-4xl">{dest.image}</span>
                      <div className="flex-1">
                        <h3 className="font-semibold text-lg mb-1">
                          {dest.name}
                        </h3>
                        <div className="flex items-center gap-3 text-sm text-muted-foreground">
                          <span>{dest.travelers} travelers</span>
                          <Badge
                            variant="secondary"
                            className="bg-primary/10 text-primary"
                          >
                            {dest.trend}
                          </Badge>
                        </div>
                      </div>
                    </div>
                  </Card>
                </motion.div>
              ))}
            </div>
          </TabsContent>
        </Tabs>
      </motion.div>
    </div>
  );
};

export default CommunityPage;
