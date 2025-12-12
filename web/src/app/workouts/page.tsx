import Link from "next/link";
import { cookies } from "next/headers";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { goFetch } from "@/lib/go-api";
import type { Workout } from "@/lib/types";

export default async function WorkoutsPage() {
  const res = await goFetch("/workouts", { method: "GET" });
  const data = (await res.json().catch(() => ({}))) as {
    workouts?: Workout[];
    error?: string;
  };

  const cookieStore = await cookies();
  const isAuthed = Boolean(cookieStore.get("auth_token")?.value);
  const workouts = data.workouts ?? [];

  return (
    <div className="mx-auto max-w-4xl px-4 py-10">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Workouts</h1>
        {isAuthed ? (
          <Button asChild>
            <Link href="/workouts/new">New workout</Link>
          </Button>
        ) : (
          <Button asChild variant="secondary">
            <Link href="/login">Login to create</Link>
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All workouts</CardTitle>
        </CardHeader>
        <CardContent>
          {!res.ok ? (
            <p className="text-sm text-muted-foreground">
              {data.error ?? "Failed to load workouts"}
            </p>
          ) : workouts.length === 0 ? (
            <p className="text-sm text-muted-foreground">No workouts yet.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Title</TableHead>
                  <TableHead>Duration</TableHead>
                  <TableHead>Calories</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {workouts.map((w) => (
                  <TableRow key={w.id}>
                    <TableCell className="font-medium">{w.title}</TableCell>
                    <TableCell>{w.duration_minutes} min</TableCell>
                    <TableCell>{w.calories_burned}</TableCell>
                    <TableCell className="text-right">
                      <Button asChild size="sm" variant="secondary">
                        <Link href={`/workouts/${w.id}`}>View</Link>
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
