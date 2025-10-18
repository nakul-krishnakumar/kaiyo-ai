import { motion } from "framer-motion";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Link } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useRouter } from "@tanstack/react-router";
import { setTokens } from "@/lib/tanstack-query";

const API_URL: string = (import.meta as any).env?.VITE_API_URL ?? "";

const formSchema = z.object({
  emailid: z.string().min(1, "email is required"),
  password: z.string().min(1, "Password is required"),
});

type FormValues = z.infer<typeof formSchema>;

export const LoginPage = () => {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      emailid: "",
      password: "",
    },
  });
  const router = useRouter();
  const loginMutation = useMutation({
    mutationFn: async (data: FormValues) => {
      const response = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          email: data.emailid,
          password: data.password,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Login failed");
      }

      return response.json();
    },
    onSuccess: async (data) => {
      setTokens(data);
      await router.navigate({ to: "/chat" });
    },
  });

  const onSubmit = (values: FormValues) => {
    loginMutation.mutate(values);
  };

  return (
    <Card className="w-full  bg-[#FAFAFA]">
      <div className="min-h-screen flex items-center justify-center px-4">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="w-full max-w-md"
        >
          <Link to="/" className="mb-4">
            <Button variant="ghost" asChild>
              <span className="flex items-center">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Home
              </span>
            </Button>
          </Link>

          <Card className="shadow-large">
            <CardHeader className="space-y-1">
              <CardTitle className="text-2xl font-bold text-center">
                Welcome to Kaiyo AI
              </CardTitle>
              <CardDescription className="text-center">
                Enter your credentials to continue
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="space-y-4"
                >
                  <div className="space-y-2">
                    <FormField
                      control={form.control}
                      name="emailid"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Email ID</FormLabel>
                          <FormControl>
                            <Input
                              placeholder="Enter your email ID"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <div className="space-y-2">
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
                  </div>
                  <Button
                    type="submit"
                    variant="hero"
                    className="w-full"
                    size="lg"
                    disabled={loginMutation.isPending}
                  >
                    {loginMutation.isPending ? "Logging in..." : "Login"}
                  </Button>
                  {loginMutation.isError && (
                    <p className="text-red-600 text-sm">
                      {(loginMutation.error as Error).message}
                    </p>
                  )}
                </form>
              </Form>
              <div className="mt-4 text-center text-sm text-muted-foreground">
                Don't have an account?{" "}
                <a href="#" className="text-primary hover:underline">
                  Create account
                </a>
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    </Card>
  );
};

export default LoginPage;
