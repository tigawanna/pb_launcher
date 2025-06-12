import type { PropsWithChildren } from "react";
import { NavLink, Outlet } from "react-router-dom";
import { Server, User, LogOut, Settings } from "lucide-react";
import { useConfirmModal } from "../hooks/useConfirmModal";
import { authService } from "../services/auth";
import { useViewportHeight } from "../hooks/useViewportHeight";

export const DASHBOARD_LAYOUT_APP_BAR_HEIGHT = 56;

export const DashboardLayout = ({ children }: PropsWithChildren) => {
  const height = useViewportHeight();

  const confirm = useConfirmModal();

  const logout = async () => {
    const confirmed = await confirm(
      "Sign out",
      "Are you sure you want to sign out?",
    );
    if (!confirmed) return;
    await authService.logout();
  };

  return (
    <div style={{ height }} className="bg-base-200 flex flex-col items-center">
      <header
        style={{ height: DASHBOARD_LAYOUT_APP_BAR_HEIGHT }}
        className="w-full bg-base-100 shadow-sm"
      >
        <div className="mx-auto w-full px-4 py-3 flex items-center justify-between">
          <NavLink
            to="/"
            className="btn btn-sm btn-ghost gap-2 text-base-content"
          >
            <Server className="w-4 h-4" />
            Services
          </NavLink>

          <div className="dropdown dropdown-end">
            <label tabIndex={0} className="btn btn-sm btn-ghost gap-2">
              <User className="w-4 h-4" />
              Account
            </label>
            <ul
              tabIndex={0}
              className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-48 mt-2 z-[1]"
            >
              <li>
                <NavLink to="/account/settings">
                  <Settings className="w-4 h-4" />
                  Settings
                </NavLink>
              </li>
              <li>
                <button onClick={logout}>
                  <LogOut className="w-4 h-4" />
                  Sign out
                </button>
              </li>
            </ul>
          </div>
        </div>
      </header>

      <main
        style={{ height: height - DASHBOARD_LAYOUT_APP_BAR_HEIGHT }}
        className="w-full flex-1 px-4 py-6"
      >
        {children || <Outlet />}
      </main>
    </div>
  );
};
