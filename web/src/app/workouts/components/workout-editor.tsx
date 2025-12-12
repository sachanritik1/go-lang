"use client";

import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";

import type { Workout, WorkoutEntry } from "@/lib/types";

type Props = {
  workout: Workout;
  canEdit: boolean;
};

function normalizeNumber(value: string): number | null {
  const trimmed = value.trim();
  if (!trimmed) return null;
  const n = Number(trimmed);
  return Number.isFinite(n) ? n : null;
}

export default function WorkoutEditor({ workout, canEdit }: Props) {
  const router = useRouter();

  const [title, setTitle] = useState(workout.title ?? "");
  const [description, setDescription] = useState(workout.description ?? "");
  const [durationMinutes, setDurationMinutes] = useState<number>(
    workout.duration_minutes ?? 0
  );
  const [caloriesBurned, setCaloriesBurned] = useState<number>(
    workout.calories_burned ?? 0
  );

  const [entries, setEntries] = useState<WorkoutEntry[]>(workout.entries ?? []);

  const nextOrderIndex = useMemo(() => {
    const max = entries.reduce(
      (acc, e) => Math.max(acc, e.order_index ?? 0),
      0
    );
    return max + 1;
  }, [entries]);

  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isEntryDialogOpen, setIsEntryDialogOpen] = useState(false);

  // dialog form state
  const [exerciseName, setExerciseName] = useState("");
  const [sets, setSets] = useState<number>(3);
  const [reps, setReps] = useState<string>("");
  const [durationSeconds, setDurationSeconds] = useState<string>("");
  const [weight, setWeight] = useState<string>("");
  const [notes, setNotes] = useState("");

  function buildEntry(): WorkoutEntry | null {
    if (!exerciseName.trim()) {
      toast.error("Exercise name is required");
      return null;
    }

    const newEntry: WorkoutEntry = {
      id: 0,
      workout_id: workout.id,
      exercise_name: exerciseName.trim(),
      sets,
      reps: normalizeNumber(reps),
      duration_seconds: normalizeNumber(durationSeconds),
      weight: normalizeNumber(weight),
      notes,
      order_index: nextOrderIndex,
    };

    return newEntry;
  }

  function resetEntryForm() {
    setExerciseName("");
    setSets(3);
    setReps("");
    setDurationSeconds("");
    setWeight("");
    setNotes("");
  }

  async function save(entriesOverride?: WorkoutEntry[]) {
    setIsSaving(true);
    try {
      const entriesToSave = entriesOverride ?? entries;

      const res = await fetch(`/api/workouts/${workout.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          title,
          description,
          duration_minutes: durationMinutes,
          calories_burned: caloriesBurned,
          // NOTE: backend expects workout_entries (not entries) for update.
          workout_entries: entriesToSave.map((e, idx) => ({
            ...e,
            order_index: idx + 1,
            workout_id: workout.id,
          })),
        }),
      });

      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        toast.error(
          (data as { error?: string }).error ?? "Failed to save workout"
        );
        return;
      }

      toast.success("Saved");
      router.refresh();
    } finally {
      setIsSaving(false);
    }
  }

  async function addEntryAndSave() {
    const newEntry = buildEntry();
    if (!newEntry) return;

    const nextEntries = [...entries, newEntry];
    setEntries(nextEntries);

    resetEntryForm();
    setIsEntryDialogOpen(false);

    await save(nextEntries);
  }

  async function deleteWorkout() {
    if (!confirm("Delete this workout?")) return;

    setIsDeleting(true);
    try {
      const res = await fetch(`/api/workouts/${workout.id}`, {
        method: "DELETE",
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        toast.error(
          (data as { error?: string }).error ?? "Failed to delete workout"
        );
        return;
      }
      toast.success("Deleted");
      router.push("/workouts");
      router.refresh();
    } finally {
      setIsDeleting(false);
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Workout</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="grid gap-2">
            <Label htmlFor="title">Title</Label>
            <Input
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={!canEdit}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={!canEdit}
            />
          </div>
          <div className="grid gap-2 sm:grid-cols-2">
            <div className="grid gap-2">
              <Label htmlFor="duration">Duration (minutes)</Label>
              <Input
                id="duration"
                type="number"
                value={durationMinutes}
                onChange={(e) => setDurationMinutes(Number(e.target.value))}
                disabled={!canEdit}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="calories">Calories</Label>
              <Input
                id="calories"
                type="number"
                value={caloriesBurned}
                onChange={(e) => setCaloriesBurned(Number(e.target.value))}
                disabled={!canEdit}
              />
            </div>
          </div>

          {canEdit ? (
            <div className="flex gap-3">
              <Button onClick={() => void save()} disabled={isSaving}>
                {isSaving ? "Saving..." : "Save"}
              </Button>
              <Button
                variant="destructive"
                onClick={deleteWorkout}
                disabled={isDeleting}
              >
                {isDeleting ? "Deleting..." : "Delete"}
              </Button>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">
              Login to edit this workout.
            </p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Entries</CardTitle>
          {canEdit ? (
            <Dialog
              open={isEntryDialogOpen}
              onOpenChange={setIsEntryDialogOpen}
            >
              <DialogTrigger asChild>
                <Button>Add entry</Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Add entry</DialogTitle>
                </DialogHeader>
                <div className="grid gap-4">
                  <div className="grid gap-2">
                    <Label htmlFor="exercise">Exercise</Label>
                    <Input
                      id="exercise"
                      value={exerciseName}
                      onChange={(e) => setExerciseName(e.target.value)}
                    />
                  </div>
                  <div className="grid gap-2 sm:grid-cols-2">
                    <div className="grid gap-2">
                      <Label htmlFor="sets">Sets</Label>
                      <Input
                        id="sets"
                        type="number"
                        min={1}
                        value={sets}
                        onChange={(e) => setSets(Number(e.target.value))}
                      />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor="reps">Reps (optional)</Label>
                      <Input
                        id="reps"
                        inputMode="numeric"
                        value={reps}
                        onChange={(e) => setReps(e.target.value)}
                      />
                    </div>
                  </div>
                  <div className="grid gap-2 sm:grid-cols-2">
                    <div className="grid gap-2">
                      <Label htmlFor="durationSeconds">
                        Duration seconds (optional)
                      </Label>
                      <Input
                        id="durationSeconds"
                        inputMode="numeric"
                        value={durationSeconds}
                        onChange={(e) => setDurationSeconds(e.target.value)}
                      />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor="weight">Weight (optional)</Label>
                      <Input
                        id="weight"
                        inputMode="decimal"
                        value={weight}
                        onChange={(e) => setWeight(e.target.value)}
                      />
                    </div>
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="notes">Notes</Label>
                    <Textarea
                      id="notes"
                      value={notes}
                      onChange={(e) => setNotes(e.target.value)}
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button onClick={addEntryAndSave} disabled={isSaving}>
                    {isSaving ? "Adding..." : "Add"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          ) : null}
        </CardHeader>
        <CardContent>
          {entries.length === 0 ? (
            <p className="text-sm text-muted-foreground">No entries yet.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>#</TableHead>
                  <TableHead>Exercise</TableHead>
                  <TableHead>Sets</TableHead>
                  <TableHead>Reps</TableHead>
                  <TableHead>Weight</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {entries
                  .slice()
                  .sort((a, b) => (a.order_index ?? 0) - (b.order_index ?? 0))
                  .map((e) => (
                    <TableRow key={`${e.order_index}-${e.exercise_name}`}>
                      <TableCell>{e.order_index}</TableCell>
                      <TableCell>{e.exercise_name}</TableCell>
                      <TableCell>{e.sets}</TableCell>
                      <TableCell>{e.reps ?? "-"}</TableCell>
                      <TableCell>{e.weight ?? "-"}</TableCell>
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
