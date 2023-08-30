import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";

interface Props {
  open: boolean;
  setOpen: (open: boolean) => void;
  setSecretPhrase: (secretPhrase: string) => void;
}

export default function DialogConfirmSecret({
  open,
  setOpen,
  setSecretPhrase: externalSetSecretPhrase,
}: Props) {
  const [secretPhrase, setSecretPhrase] = useState("");

  const handleOpenChange = (open: boolean) => {
    if (secretPhrase.length === 0) {
      return;
    }

    if (!open) {
      externalSetSecretPhrase(secretPhrase);
    }

    setOpen(open);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent
        onCloseAutoFocus={(event) => {
          event.preventDefault();
        }}
        className="flex flex-col items-center justify-center"
        style={{
          width: "auto",
          borderColor: "transparent",
        }}
      >
        <Label>Secret Phrase</Label>
        <Input
          onChange={(e) => setSecretPhrase(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && secretPhrase.length > 0) {
              handleOpenChange(false);
            }
          }}
        />
      </DialogContent>
    </Dialog>
  );
}
