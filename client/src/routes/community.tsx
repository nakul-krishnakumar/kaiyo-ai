import CommunityPage from "@/pages/CommunityPage";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/community")({
  component: CommunityPage,
});
