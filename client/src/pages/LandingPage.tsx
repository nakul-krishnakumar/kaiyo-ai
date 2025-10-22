import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { MessageCircle, MapPin } from "lucide-react";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-50 to-blue-50 flex items-center justify-center p-4">
      <div className="max-w-4xl w-full">
        <div className="text-center mb-12">
          <h1 className="text-6xl font-bold text-gray-900 mb-4">Kaiyo AI</h1>
          <p className="text-xl text-gray-600 mb-8">
            Your Personal Travel Planning Assistant
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          <Card className="hover:shadow-lg transition-all duration-300 cursor-pointer">
            <CardHeader className="text-center">
              <div className="mx-auto mb-4 p-4 bg-purple-100 rounded-full w-16 h-16 flex items-center justify-center">
                <MapPin className="w-8 h-8 text-purple-600" />
              </div>
              <CardTitle className="text-2xl">Plan Your Trip</CardTitle>
              <CardDescription className="text-lg">
                Get personalized travel recommendations and itineraries powered
                by AI
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Link to="/chat">
                <Button className="w-full" size="lg">
                  Start Planning
                </Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-all duration-300 cursor-pointer">
            <CardHeader className="text-center">
              <div className="mx-auto mb-4 p-4 bg-blue-100 rounded-full w-16 h-16 flex items-center justify-center">
                <MessageCircle className="w-8 h-8 text-blue-600" />
              </div>
              <CardTitle className="text-2xl">Talk to Community</CardTitle>
              <CardDescription className="text-lg">
                Connect with fellow travelers and share experiences
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Link to="/community">
                <Button className="w-full" variant="outline" size="lg">
                  Join Community
                </Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
