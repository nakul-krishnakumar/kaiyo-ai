import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { MessageCircle, Heart, Share, Search, Plus } from "lucide-react";

export default function CommunityPage() {
  const posts = [
    {
      id: 1,
      author: "Sarah Johnson",
      avatar: "SJ",
      title: "Amazing 5-day Kerala backwaters experience!",
      content:
        "Just got back from an incredible houseboat journey through Alleppey and Kumarakom. The sunset views were absolutely breathtaking! Here are my top recommendations...",
      likes: 24,
      comments: 8,
      timestamp: "2 hours ago",
      tags: ["Kerala", "Backwaters", "Houseboat"],
    },
    {
      id: 2,
      author: "Mike Chen",
      avatar: "MC",
      title: "Budget-friendly Goa itinerary under ₹15,000",
      content:
        "Spent 4 days in Goa without breaking the bank! Stayed in hostels, ate at local joints, and still had an amazing time. Here's how I did it...",
      likes: 42,
      comments: 15,
      timestamp: "5 hours ago",
      tags: ["Goa", "Budget", "Beach"],
    },
    {
      id: 3,
      author: "Priya Sharma",
      avatar: "PS",
      title: "Solo female travel in Himachal - Safety tips",
      content:
        "Just completed a solo trek in Himachal Pradesh. Sharing some important safety tips and beautiful spots that are perfect for solo female travelers...",
      likes: 67,
      comments: 23,
      timestamp: "1 day ago",
      tags: ["Himachal", "Solo Travel", "Trekking", "Safety"],
    },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white comic-border border-b-4 p-4">
        <div className="max-w-6xl mx-auto flex items-center justify-between">
          <h1 className="text-3xl font-bold">Travel Community</h1>
          <Button>
            <Plus className="w-4 h-4 mr-2" />
            New Post
          </Button>
        </div>
      </div>

      <div className="max-w-6xl mx-auto p-4">
        <div className="grid lg:grid-cols-4 gap-6">
          {/* Sidebar */}
          <div className="lg:col-span-1">
            <Card className="mb-4">
              <CardHeader>
                <CardTitle className="text-lg">Popular Topics</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {[
                    "Kerala Backwaters",
                    "Goa Beaches",
                    "Himachal Trekking",
                    "Rajasthan Culture",
                    "South India Food",
                  ].map((topic) => (
                    <Button
                      key={topic}
                      variant="outline"
                      size="sm"
                      className="w-full justify-start"
                    >
                      #{topic.replace(" ", "")}
                    </Button>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Travel Stats</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div>
                    <p className="text-sm text-gray-600">Active Travelers</p>
                    <p className="text-2xl font-bold text-purple-600">2,847</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Posts This Week</p>
                    <p className="text-2xl font-bold text-blue-600">156</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Countries Covered</p>
                    <p className="text-2xl font-bold text-green-600">23</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Main Content */}
          <div className="lg:col-span-3">
            {/* Search Bar */}
            <Card className="mb-6">
              <CardContent className="p-4">
                <div className="flex space-x-2">
                  <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                    <Input
                      placeholder="Search posts, destinations, or travelers..."
                      className="pl-10"
                    />
                  </div>
                  <Button>Search</Button>
                </div>
              </CardContent>
            </Card>

            {/* Posts */}
            <div className="space-y-6">
              {posts.map((post) => (
                <Card
                  key={post.id}
                  className="hover:shadow-lg transition-shadow"
                >
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex items-center space-x-3">
                        <div className="w-10 h-10 bg-purple-100 rounded-full flex items-center justify-center">
                          <span className="text-sm font-semibold text-purple-600">
                            {post.avatar}
                          </span>
                        </div>
                        <div>
                          <p className="font-semibold">{post.author}</p>
                          <p className="text-sm text-gray-500">
                            {post.timestamp}
                          </p>
                        </div>
                      </div>
                    </div>
                    <CardTitle className="text-xl mt-3">{post.title}</CardTitle>
                  </CardHeader>

                  <CardContent>
                    <p className="text-gray-700 mb-4">{post.content}</p>

                    <div className="flex flex-wrap gap-2 mb-4">
                      {post.tags.map((tag) => (
                        <span
                          key={tag}
                          className="px-2 py-1 bg-purple-100 text-purple-700 text-xs rounded-full"
                        >
                          #{tag}
                        </span>
                      ))}
                    </div>

                    <div className="flex items-center space-x-4 pt-4 border-t">
                      <Button variant="ghost" size="sm">
                        <Heart className="w-4 h-4 mr-1" />
                        {post.likes}
                      </Button>
                      <Button variant="ghost" size="sm">
                        <MessageCircle className="w-4 h-4 mr-1" />
                        {post.comments}
                      </Button>
                      <Button variant="ghost" size="sm">
                        <Share className="w-4 h-4 mr-1" />
                        Share
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
