import Link from "next/link";
import { Button } from "./ui/button";

import { ChevronLeft } from "lucide-react";

export default function ButtonBackward() {
  return (
    <Button
      asChild
      variant={"outline"}
      className="m-8"
      onClick={() => {
        localStorage.removeItem("secret");
      }}
    >
      <Link href={"/"}>
        <ChevronLeft />
      </Link>
    </Button>
  );
}
