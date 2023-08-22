import { ISecret } from "@/entities/entites";
import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

interface GetSecretRequest {
  secretKey: string;
}

interface GetSecretResponse {
  secret: ISecret;
}

interface CreateSecretRequest {
  message: string;
}

interface CreateSecretResponse {
  secretKey: string;
}

export const secretsApi = createApi({
  reducerPath: "secretsApi",
  baseQuery: fetchBaseQuery({
    baseUrl: `http://${import.meta.env.VITE_API_URL}/api`,
  }),
  endpoints: (builder) => ({
    getSecret: builder.query<GetSecretResponse, GetSecretRequest>({
      query: ({ secretKey }) => ({
        url: `/secrets/${secretKey}`,
      }),
    }),

    createSecret: builder.mutation<CreateSecretResponse, CreateSecretRequest>({
      query: ({ message }) => ({
        url: `/secrets`,
        method: "POST",
        body: { message },
      }),
    }),
  }),
});

export const { useGetSecretQuery, useCreateSecretMutation } = secretsApi;
