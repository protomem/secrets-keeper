import { configureStore } from "@reduxjs/toolkit";
import { secretsApi } from "./secrets/secrets.api";

export const store = configureStore({
  reducer: {
    [secretsApi.reducerPath]: secretsApi.reducer,
  },
  devTools: true,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(secretsApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
