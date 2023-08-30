import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { ClipboardCheck, ClipboardEdit } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { copyTextToClipboard } from "@/lib/clipboard";
import { useEffect, useState } from "react";

interface Props {
  open: boolean;
  setOpen: (open: boolean) => void;
  secretKey: string;
  withSecretPhrase: boolean;
}

export default function DialogSecretCreated({
  open,
  setOpen,
  secretKey,
  withSecretPhrase,
}: Props) {
  const [copied, setCopied] = useState(false);
  useEffect(() => {
    setCopied(false);
  }, [open]);

  let linkToSecret = `${window.location.origin}/secrets/${secretKey.trim()}`;
  if (withSecretPhrase) {
    linkToSecret += `?withSecretPhrase=true`;
  }

  const handleClick = async () => {
    await copyTextToClipboard(linkToSecret);
    setCopied(true);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
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
        <h2>Secret Created</h2>

        <Badge
          variant="outline"
          className="text-start text-sm w-max overflow-x-auto"
        >
          {linkToSecret}
          <Button
            size="icon"
            variant="outline"
            className="m-2"
            onClick={handleClick}
          >
            {copied ? (
              <ClipboardCheck scale={"sm"} />
            ) : (
              <ClipboardEdit scale={"sm"} />
            )}
          </Button>
        </Badge>
      </DialogContent>
    </Dialog>
  );
}
