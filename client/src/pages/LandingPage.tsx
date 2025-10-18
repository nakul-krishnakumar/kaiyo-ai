import { motion } from "framer-motion";
import { useRouter } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Compass, TrendingUp } from "lucide-react";
import { isAuthenticated } from "@/lib/tanstack-query";

const LandingPage = () => {
  const router = useRouter();

  const handlePlanTrip = () => {
    if (isAuthenticated()) {
      router.navigate({ to: "/chat" });
    } else {
      router.navigate({ to: "/login" });
    }
  };

  const handleTrends = () => {
    if (isAuthenticated()) {
      router.navigate({ to: "/community" });
    } else {
      router.navigate({ to: "/login" });
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center px-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="text-center max-w-4xl"
      >
        <motion.h1
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.2, duration: 0.5 }}
          className="text-6xl md:text-7xl font-bold mb-6 bg-gradient-hero bg-clip-text text-transparent"
        >
          Kaiyo AI
        </motion.h1>

        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4, duration: 0.5 }}
          className="text-xl md:text-2xl text-muted-foreground mb-12"
        >
          Your Personal AI Travel Planner
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6, duration: 0.5 }}
          className="flex flex-col sm:flex-row gap-4 justify-center items-center"
        >
          <Button
            variant="hero"
            size="lg"
            onClick={handlePlanTrip}
            className="w-full sm:w-auto min-w-[240px]"
          >
            <Compass className="mr-2" />
            Plan Your Next Trip
          </Button>

          <Button
            variant="default"
            size="lg"
            onClick={handleTrends}
            className="w-full sm:w-auto min-w-[240px]"
          >
            <TrendingUp className="mr-2" />
            Know the Trend
          </Button>
        </motion.div>
      </motion.div>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.8, duration: 0.5 }}
        className="mt-20 text-center"
      >
        <p className="text-sm text-muted-foreground">
          Discover personalized travel itineraries powered by AI
        </p>
      </motion.div>
    </div>
  );
};

export default LandingPage;
