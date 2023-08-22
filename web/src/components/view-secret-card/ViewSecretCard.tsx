import { ISecret } from "@/entities/entites";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Props {
  secret: ISecret;
}

export default function ViewSecretCard({ secret }: Props) {
  return (
    <Card className="w-[50rem] m-8">
      <CardHeader>
        <CardTitle>Secret</CardTitle>
      </CardHeader>
      <CardContent>{secret.message}</CardContent>
    </Card>
  );
}
