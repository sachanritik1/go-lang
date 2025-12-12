import { NextResponse } from "next/server";

import { getAuthTokenFromCookies, getGoApiBaseUrl } from "@/lib/go-api";

export async function GET() {
  const token = await getAuthTokenFromCookies();
  if (!token) {
    return NextResponse.json(
      { error: "authentication required" },
      { status: 401 }
    );
  }

  const res = await fetch(`${getGoApiBaseUrl()}/users/self`, {
    method: "GET",
    headers: {
      Accept: "application/json",
      Authorization: `Bearer ${token}`,
    },
    cache: "no-store",
  });

  const data = await res.json().catch(() => ({}));
  return NextResponse.json(data, { status: res.status });
}
