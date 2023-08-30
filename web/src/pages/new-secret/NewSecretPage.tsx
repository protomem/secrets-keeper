import { useState } from "react";

import DialogSecretCreated from "@/components/dialog-secret-created/DialogSecretCreated";
import NewSecretCard from "@/components/new-secret-card/NewSecretCard";

export default function NewSecretPage() {
  const [openDialog, setOpenDialog] = useState(false);
  const [secretKey, setSecretKey] = useState("");
  const [withSecretPhrase, setWithSecretPhrase] = useState(false);

  const onSubmit = (secretKey: string, withSecretPhrase: boolean) => {
    setOpenDialog(true);
    setSecretKey(secretKey);
    setWithSecretPhrase(withSecretPhrase);
  };

  return (
    <div className="flex flex-col items-center justify-center">
      <NewSecretCard onSubmit={onSubmit} />
      <DialogSecretCreated
        open={openDialog}
        setOpen={setOpenDialog}
        secretKey={secretKey}
        withSecretPhrase={withSecretPhrase}
      />
    </div>
  );
}
