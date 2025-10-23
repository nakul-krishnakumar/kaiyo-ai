import CommunityPage from "@/pages/CommunityPage";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { isAuthenticated } from "@/lib/tanstack-query";
export const Route = createFileRoute("/community")({
  component: CommunityPage,
  beforeLoad: () => {
    if (!isAuthenticated()) {
      throw redirect({
        to: "/login",
      });
    }
  },
});
