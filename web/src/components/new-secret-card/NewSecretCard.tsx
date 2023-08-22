import NewSecretForm from "../new-secret-form/NewSecretForm";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";

interface Props {
  onSubmit: (secretKey: string) => void;
}

export default function NewSecretCard({ onSubmit }: Props) {
  return (
    <Card className="w-[50rem] m-8">
      <CardHeader>
        <CardTitle className="text-center">New Secret</CardTitle>
      </CardHeader>

      <CardContent>
        <NewSecretForm onSubmit={onSubmit} />
      </CardContent>
    </Card>
  );
}
