import { type PropsWithChildren } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { LoginPage } from "../pages/auth/LoginPage";
import { RegisterPage } from "../pages/auth/RegisterPage";
import { useSession } from "../hooks/useSession";
import { useQuery } from "@tanstack/react-query";
import { authService } from "../services/auth";

import { QueryErrorView } from "../components/QueryErrorView";
import { DashboardLayout } from "../layouts/DashboardLayout";
import { ServicesPage } from "../pages/internal/ServicesPage";

const PrivateRoute = ({
  children,
  redirectTo,
}: PropsWithChildren<{ redirectTo: string }>) => {
  const { user } = useSession();
  return user ? children : <Navigate to={redirectTo} replace />;
};

const PublicRoute = ({ children }: PropsWithChildren) => {
  const { user } = useSession();
  return !user ? children : <Navigate to="/" replace />;
};

export const AppRoutes = () => {
  const adminExistsQuery = useQuery<boolean>({
    queryKey: ["admin-created"],
    queryFn: async ({ signal }) => {
      const setupDone = await authService.isInitialSetupDone(signal);
      if (!setupDone) await authService.logout();
      return setupDone;
    },
  });

  if (adminExistsQuery.isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <span className="text-gray-500 text-lg animate-pulse">Loading...</span>
      </div>
    );
  }

  if (adminExistsQuery.isError)
    return (
      <QueryErrorView
        error={adminExistsQuery.error}
        onRetry={adminExistsQuery.refetch}
      />
    );

  if (adminExistsQuery.data == null) {
    return (
      <QueryErrorView
        error={new Error("Unexpected empty response from server.")}
        onRetry={adminExistsQuery.refetch}
      />
    );
  }
  const isSetupDone = adminExistsQuery.data;
  return (
    <BrowserRouter>
      <Routes>
        {isSetupDone ? (
          <Route
            path="/login"
            element={
              <PublicRoute>
                <LoginPage />
              </PublicRoute>
            }
          />
        ) : (
          <Route
            path="/register"
            element={
              <PublicRoute>
                <RegisterPage refresh={adminExistsQuery.refetch} />
              </PublicRoute>
            }
          />
        )}
        <Route
          path="*"
          element={
            <PrivateRoute redirectTo={isSetupDone ? "/login" : "/register"}>
              <DashboardLayout />
            </PrivateRoute>
          }
        >
          <Route index element={<ServicesPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};
