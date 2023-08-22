import { Dialog, DialogContent } from "@/components/ui/dialog";

interface Props {
  open: boolean;
  setOpen: (open: boolean) => void;
  secretKey: string;
}

export default function DialogSecretCreated({
  open,
  setOpen,
  secretKey,
}: Props) {
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent
        onCloseAutoFocus={(event) => {
          event.preventDefault();
        }}
      >
        <h2>Secret Created</h2>
        {secretKey}
      </DialogContent>
    </Dialog>
  );
}
