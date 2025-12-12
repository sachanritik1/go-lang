import "server-only";

import { cookies } from "next/headers";

export function getGoApiBaseUrl(): string {
  const url = process.env.GO_API_URL;
  const resolved = (
    url && url.length > 0 ? url : "http://localhost:8080"
  ).replace(/\/$/, "");
  return resolved;
}

export async function getAuthTokenFromCookies(): Promise<string | undefined> {
  const cookieStore = await cookies();
  return cookieStore.get("auth_token")?.value;
}

export async function goFetch(path: string, init: RequestInit = {}) {
  const base = getGoApiBaseUrl();
  const token = await getAuthTokenFromCookies();

  const headers = new Headers(init.headers);
  headers.set("Accept", "application/json");
  if (init.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  return fetch(`${base}${path.startsWith("/") ? path : `/${path}`}`, {
    ...init,
    headers,
    cache: "no-store",
  });
}
