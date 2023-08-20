"use client";

import ButtonBackward from "@/components/button-backward";
import ViewSecret from "@/components/view-secret";

interface Props {
  params: { secretKey: string };
}

export default function Page({}: Props) {
  return (
    <div className="flex flex-row items-start justify-between">
      <div className="basis-1/3">
        <ButtonBackward />
      </div>

      <div className="basis-1/3">
        <ViewSecret />
      </div>

      <div className="basis-1/3"></div>
    </div>
  );
}
