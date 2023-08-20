import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function ViewSecret() {
  return (
    <Card className="w-[50rem] m-8">
      <CardHeader>
        <CardTitle className="text-center">Secret</CardTitle>
      </CardHeader>
      <CardContent>{localStorage.getItem("secret")}</CardContent>
    </Card>
  );
}
