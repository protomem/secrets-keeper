import { ISecret } from "@/entities/entites";
import { slug } from "@/lib/slug";
import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

interface GetSecretRequest {
  secretKey: string;
  secretPhrase?: string;
}

interface GetSecretResponse {
  secret: ISecret;
}

interface CreateSecretRequest {
  message: string;
  ttl: number;
  secretPhrase?: string;
}

interface CreateSecretResponse {
  secretKey: string;
  withSecretPhrase: boolean;
}

export const secretsApi = createApi({
  reducerPath: "secretsApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api`,
  }),
  endpoints: (builder) => ({
    getSecret: builder.query<GetSecretResponse, GetSecretRequest>({
      query: ({ secretKey, secretPhrase }) => ({
        url:
          secretPhrase !== undefined
            ? `/secrets/${secretKey}?secretPhrase=${slug(secretPhrase)}`
            : `/secrets/${secretKey}`,
      }),
    }),

    createSecret: builder.mutation<CreateSecretResponse, CreateSecretRequest>({
      query: ({ message, ttl, secretPhrase }) => ({
        url: `/secrets`,
        method: "POST",
        body: { message, ttl, secretPhrase: slug(secretPhrase) },
      }),
    }),
  }),
});

export const { useGetSecretQuery, useCreateSecretMutation } = secretsApi;
