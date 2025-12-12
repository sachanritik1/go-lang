import type { Metadata } from "next";
import { cookies } from "next/headers";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

import TopNav from "@/components/top-nav";
import { Toaster } from "@/components/ui/sonner";
import { goFetch } from "@/lib/go-api";
import type { User } from "@/lib/types";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Workouts",
  description: "Workout tracker frontend",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const cookieStore = await cookies();
  const isAuthed = Boolean(cookieStore.get("auth_token")?.value);
  const username = cookieStore.get("username")?.value;

  let selfUser: User | undefined;
  if (isAuthed) {
    const res = await goFetch("/users/self", { method: "GET" });
    const data = (await res.json().catch(() => ({}))) as {
      user?: User;
    };
    if (res.ok && data.user) {
      selfUser = data.user;
    }
  }

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <TopNav isAuthed={isAuthed} username={selfUser?.username ?? username} />
        {children}
        <Toaster richColors />
      </body>
    </html>
  );
}
