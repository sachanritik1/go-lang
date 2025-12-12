import { NextResponse } from "next/server";

import { getGoApiBaseUrl } from "@/lib/go-api";

export async function POST(req: Request) {
  const body = await req.json().catch(() => null);
  if (
    !body ||
    typeof body.username !== "string" ||
    typeof body.password !== "string"
  ) {
    return NextResponse.json(
      { error: "invalid request payload" },
      { status: 400 }
    );
  }

  const res = await fetch(`${getGoApiBaseUrl()}/tokens/authentication`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify({ username: body.username, password: body.password }),
    cache: "no-store",
  });

  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    return NextResponse.json(data, { status: res.status });
  }

  const authTokenRaw = (data as { auth_token?: unknown }).auth_token;
  const token =
    typeof authTokenRaw === "string"
      ? authTokenRaw
      : typeof authTokenRaw === "object" &&
        authTokenRaw !== null &&
        "token" in authTokenRaw &&
        typeof (authTokenRaw as { token?: unknown }).token === "string"
      ? ((authTokenRaw as { token: string }).token as string)
      : undefined;

  if (typeof token !== "string" || token.length === 0) {
    return NextResponse.json(
      { error: "missing auth_token in response" },
      { status: 502 }
    );
  }

  const response = NextResponse.json({ ok: true });
  response.cookies.set("auth_token", token, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24,
  });

  // Used only for UI; still kept HttpOnly so it can be read server-side.
  response.cookies.set("username", body.username, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24,
  });

  return response;
}
