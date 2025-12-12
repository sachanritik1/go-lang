export type ApiEnvelope<T extends Record<string, unknown>> = T;

export type WorkoutEntry = {
  id: number;
  workout_id: number;
  exercise_name: string;
  sets: number;
  reps: number | null;
  duration_seconds: number | null;
  weight: number | null;
  notes: string;
  order_index: number;
};

export type Workout = {
  id: number;
  title: string;
  description: string;
  duration_minutes: number;
  calories_burned: number;
  entries: WorkoutEntry[];
  user_id: number;
};

export type User = {
  id: number;
  username: string;
  email: string;
  bio: string;
  created_at: string;
  updated_at: string;
};
