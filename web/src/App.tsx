import RouterProvider from "@/router/RouterProvider";
import ThemeProvider from "@/theme/ThemeProvider";
import StoreProvider from "@/feature/StoreProvider";

export default function App() {
  return (
    <>
      <ThemeProvider>
        <StoreProvider>
          <RouterProvider />
        </StoreProvider>
      </ThemeProvider>
    </>
  );
}
