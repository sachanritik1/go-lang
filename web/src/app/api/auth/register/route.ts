import { NextResponse } from "next/server";

import { getGoApiBaseUrl } from "@/lib/go-api";

export async function POST(req: Request) {
  const body = await req.json().catch(() => null);
  if (
    !body ||
    typeof body.username !== "string" ||
    typeof body.email !== "string" ||
    typeof body.password !== "string"
  ) {
    return NextResponse.json(
      { error: "invalid request payload" },
      { status: 400 }
    );
  }

  const res = await fetch(`${getGoApiBaseUrl()}/users`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify({
      username: body.username,
      email: body.email,
      password: body.password,
      bio: typeof body.bio === "string" ? body.bio : "",
    }),
    cache: "no-store",
  });

  const data = await res.json().catch(() => ({}));
  return NextResponse.json(data, { status: res.status });
}
