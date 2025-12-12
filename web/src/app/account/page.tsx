import Link from "next/link";
import { cookies } from "next/headers";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { goFetch } from "@/lib/go-api";
import type { User } from "@/lib/types";

export default async function AccountPage() {
  const cookieStore = await cookies();
  const isAuthed = Boolean(cookieStore.get("auth_token")?.value);

  if (!isAuthed) {
    return (
      <div className="mx-auto max-w-xl px-4 py-10">
        <Card>
          <CardHeader>
            <CardTitle>Account</CardTitle>
            <CardDescription>You are not logged in.</CardDescription>
          </CardHeader>
          <CardContent className="flex gap-3">
            <Button asChild>
              <Link href="/login">Login</Link>
            </Button>
            <Button asChild variant="secondary">
              <Link href="/register">Register</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const res = await goFetch("/users/self", { method: "GET" });
  const data = (await res.json().catch(() => ({}))) as {
    user?: User;
    error?: string;
  };

  if (!res.ok || !data.user) {
    return (
      <div className="mx-auto max-w-xl px-4 py-10">
        <Card>
          <CardHeader>
            <CardTitle>Account</CardTitle>
            <CardDescription>Could not load your account.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            <p className="text-sm text-muted-foreground">
              {data.error ?? "Please try logging in again."}
            </p>
            <div className="flex gap-3">
              <Button asChild variant="secondary">
                <Link href="/login">Login</Link>
              </Button>
              <form action="/api/auth/logout" method="post">
                <Button type="submit">Clear session</Button>
              </form>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-xl px-4 py-10">
      <Card>
        <CardHeader>
          <CardTitle>Account</CardTitle>
          <CardDescription>Logged in user details.</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="grid gap-2 text-sm">
            <div>
              <span className="text-muted-foreground">Username:</span>{" "}
              <span className="font-medium">{data.user.username}</span>
            </div>
            <div>
              <span className="text-muted-foreground">Email:</span>{" "}
              <span className="font-medium">{data.user.email}</span>
            </div>
            <div>
              <span className="text-muted-foreground">Bio:</span>{" "}
              <span className="font-medium">{data.user.bio || "-"}</span>
            </div>
          </div>

          <div className="flex gap-3">
            <Button asChild variant="secondary">
              <Link href="/workouts">Go to workouts</Link>
            </Button>
            <form action="/api/auth/logout" method="post">
              <Button type="submit">Logout</Button>
            </form>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
