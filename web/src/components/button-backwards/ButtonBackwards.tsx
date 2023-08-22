import { Button } from "@/components/ui/button";
import { ChevronLeft } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function ButtonBackwards() {
  const nav = useNavigate();

  return (
    <Button
      asChild
      variant={"outline"}
      className="m-8"
      size={"icon"}
      onClick={(event) => {
        event.preventDefault();
        nav("/", { replace: true });
      }}
    >
      <ChevronLeft className="w-12 h-9" />
    </Button>
  );
}
