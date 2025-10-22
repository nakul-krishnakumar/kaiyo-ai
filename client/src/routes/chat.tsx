import { createFileRoute, redirect } from "@tanstack/react-router";
import ChatPage from "@/pages/ChatPage";
import { isAuthenticated } from "@/lib/tanstack-query";

export const Route = createFileRoute("/chat")({
  component: ChatPage,
  beforeLoad: () => {
    if (!isAuthenticated()) {
      throw redirect({
        to: "/login",
      });
    }
  },
});
