import Link from "next/link";
import { cookies } from "next/headers";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import WorkoutEditor from "@/app/workouts/components/workout-editor";
import { goFetch } from "@/lib/go-api";
import type { Workout } from "@/lib/types";

export default async function WorkoutDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  const res = await goFetch(`/workouts/${id}`, { method: "GET" });
  const data = (await res.json().catch(() => ({}))) as {
    workout?: Workout;
    error?: string;
  };

  if (!res.ok || !data.workout) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-10">
        <Card>
          <CardHeader>
            <CardTitle>Workout</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              {data.error ?? "Workout not found"}
            </p>
            <div className="mt-4">
              <Link className="underline" href="/workouts">
                Back to workouts
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  const cookieStore = await cookies();
  const canEdit = Boolean(cookieStore.get("auth_token")?.value);

  return (
    <div className="mx-auto max-w-3xl px-4 py-10">
      <div className="mb-6">
        <Link className="underline" href="/workouts">
          Back to workouts
        </Link>
      </div>
      <WorkoutEditor workout={data.workout} canEdit={canEdit} />
    </div>
  );
}
