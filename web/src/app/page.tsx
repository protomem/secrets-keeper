"use client";

import { useRouter } from "next/navigation";

import NewSecret from "@/components/new-secret";

export default function Home() {
  const route = useRouter();
  const key = "secretKey";

  return (
    <main className="flex flex-col items-center">
      <NewSecret
        onSubmit={() => {
          route.push(`/secrets/${key}`);
        }}
      />
    </main>
  );
}
