import ButtonBackwards from "@/components/button-backwards/ButtonBackwards";
import ViewSecretCard from "@/components/view-secret-card/ViewSecretCard";
import { ISecret } from "@/entities/entites";
import { useGetSecretQuery } from "@/feature/secrets/secrets.api";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function ViewSecretPage() {
  const params = useParams();
  const [secret, setSecret] = useState<ISecret | null>(null);
  const { data, isSuccess, isError, isLoading } = useGetSecretQuery({
    secretKey: params.secretKey || "",
  });

  useEffect(() => {
    if (!data || secret !== null) return;

    setSecret(data.secret);
  }, [data, secret]);

  return (
    <div className="flex flex-row items-start justify-between">
      <div className="basis-1/3">
        <ButtonBackwards />
      </div>

      <div className="basis-1/3">
        {isLoading && (
          <div className="m-8 text-center">
            <h2 className="text-3xl">Loading...</h2>
          </div>
        )}
        {isSuccess && secret && <ViewSecretCard secret={secret} />}
        {isError && (
          <div className="m-8 text-center">
            <h2 className="text-3xl">Secret Not Found</h2>
          </div>
        )}
      </div>

      <div className="basis-1/3" />
    </div>
  );
}
