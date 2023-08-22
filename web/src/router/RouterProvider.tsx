import { BrowserRouter, Route, Routes } from "react-router-dom";

import NewSecretPage from "@/pages/new-secret/NewSecretPage";
import NotFoundPage from "@/pages/not-found/NotFoundPage";
import ViewSecretPage from "@/pages/view-secret/ViewSecretPage";

export default function RouterProvider() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<NewSecretPage />} index />
        <Route path="/secrets/:secretKey" element={<ViewSecretPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  );
}
