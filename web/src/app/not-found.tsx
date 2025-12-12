import Link from "next/link";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function NotFound() {
  return (
    <div className="mx-auto max-w-xl px-4 py-10">
      <Card>
        <CardHeader>
          <CardTitle>Not found</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            The page you requested does not exist.
          </p>
          <div className="mt-4">
            <Link className="underline" href="/workouts">
              Go to workouts
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
