import "./index.css";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { ConfirmModalProvider } from "./hooks/confirm-modal/provider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";

import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { AppRoutes } from "./routes/AppRoutes";
import { ModalProvider } from "./components/modal/ModalProvider";

const queryClient = new QueryClient({
  defaultOptions: { queries: { refetchOnWindowFocus: false } },
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ModalProvider>
      <Toaster reverseOrder={false} />
      <QueryClientProvider client={queryClient}>
        <ConfirmModalProvider>
          <AppRoutes />
        </ConfirmModalProvider>
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </ModalProvider>
  </StrictMode>,
);
