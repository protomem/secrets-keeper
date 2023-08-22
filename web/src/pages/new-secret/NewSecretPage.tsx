import { useState } from "react";

import DialogSecretCreated from "@/components/dialog-secret-created/DialogSecretCreated";
import NewSecretCard from "@/components/new-secret-card/NewSecretCard";

export default function NewSecretPage() {
  const [openDialog, setOpenDialog] = useState(false);
  const [secretKey, setSecretKey] = useState("");

  const onSubmit = (secretKey: string) => {
    setOpenDialog(true);
    setSecretKey(secretKey);
  };

  return (
    <div className="flex flex-col items-center justify-center">
      <NewSecretCard onSubmit={onSubmit} />
      <DialogSecretCreated
        open={openDialog}
        setOpen={setOpenDialog}
        secretKey={secretKey}
      />
    </div>
  );
}
