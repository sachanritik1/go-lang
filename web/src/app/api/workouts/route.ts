import { NextResponse } from "next/server";

import { getAuthTokenFromCookies, getGoApiBaseUrl } from "@/lib/go-api";

export async function GET() {
  const token = await getAuthTokenFromCookies();

  const res = await fetch(`${getGoApiBaseUrl()}/workouts`, {
    method: "GET",
    headers: {
      Accept: "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    cache: "no-store",
  });

  const data = await res.json().catch(() => ({}));
  return NextResponse.json(data, { status: res.status });
}

export async function POST(req: Request) {
  const token = await getAuthTokenFromCookies();
  if (!token) {
    return NextResponse.json(
      { error: "authentication required" },
      { status: 401 }
    );
  }

  const body = await req.json().catch(() => null);
  if (!body) {
    return NextResponse.json(
      { error: "invalid request payload" },
      { status: 400 }
    );
  }

  const res = await fetch(`${getGoApiBaseUrl()}/workouts`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
    cache: "no-store",
  });

  const data = await res.json().catch(() => ({}));
  return NextResponse.json(data, { status: res.status });
}
