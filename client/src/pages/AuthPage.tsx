import { useRouter, Link } from "@tanstack/react-router";
import { setTokens } from "@/lib/tanstack-query";
import { useState } from "react";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

// ---------------- SCHEMA ----------------
const formSchema = z.object({
  emailid: z.string().email("Enter a valid email"),
  password: z.string().min(6, "Password must be at least 6 characters"),
  fullname: z.string().optional(), // only for signup
});

const API_URL: string = (import.meta as any).env?.VITE_API_URL ?? "";
console.log(API_URL);
// ---------------- COMPONENT ----------------
export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true);
  const router = useRouter();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      emailid: "",
      password: "",
      fullname: "",
    },
  });

  // ---------------- MUTATION ----------------
  const authMutation = useMutation({
    mutationFn: async (data: z.infer<typeof formSchema>) => {
      const endpoint = isLogin ? "login" : "signup";
      const response = await fetch(`${API_URL}/auth/${endpoint}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          email: data.emailid,
          password: data.password,
          fullname: !isLogin ? data.fullname : undefined,
        }),
      });
      console.log("Response status:", response.status);
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `${endpoint} failed`);
      }

      return response.json();
    },
    onSuccess: async (data) => {
      setTokens(data);
      await router.navigate({ to: "/chat" });
    },
  });

  const onSubmit = (values: z.infer<typeof formSchema>) => {
    authMutation.mutate(values);
  };

  // ---------------- UI ----------------
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-50 to-blue-50 flex flex-col justify-center p-4">
      <Link to="/" className="mb-4 self-start">
        <Button variant="ghost">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Home
        </Button>
      </Link>

      <Card className="w-full max-w-md mx-auto">
        <CardHeader className="text-center">
          <CardTitle className="text-3xl font-bold">Kaiyo AI</CardTitle>
          <CardDescription>
            {isLogin ? "Welcome back! Please log in." : "Create your account"}
          </CardDescription>
        </CardHeader>

        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              {/* Full name field (signup only) */}
              {!isLogin && (
                <FormField
                  control={form.control}
                  name="fullname"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Full Name</FormLabel>
                      <FormControl>
                        <Input placeholder="Enter your full name" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              {/* Email */}
              <FormField
                control={form.control}
                name="emailid"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter your email" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Password */}
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Enter your password"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Submit */}
              <Button
                type="submit"
                className="w-full"
                size="lg"
                disabled={authMutation.isPending}
              >
                {authMutation.isPending
                  ? isLogin
                    ? "Logging in..."
                    : "Signing up..."
                  : isLogin
                    ? "Login"
                    : "Sign Up"}
              </Button>

              {authMutation.isError && (
                <p className="text-red-600 text-sm text-center">
                  {(authMutation.error as Error).message}
                </p>
              )}
            </form>
          </Form>

          <div className="mt-4 text-center">
            <button
              type="button"
              onClick={() => setIsLogin(!isLogin)}
              className="text-sm text-purple-600 hover:underline"
            >
              {isLogin
                ? "Don't have an account? Sign up"
                : "Already have an account? Sign in"}
            </button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
