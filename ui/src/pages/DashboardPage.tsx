import { useConfirmModal } from "../hooks/useConfirmModal";
import { authService } from "../services/auth";

export const DashboardPage = () => {
  const confirm = useConfirmModal();
  return (
    <div>
      <button
        className="btn"
        onClick={async () => {
          const ok = await confirm("Sms", "Seguro?");
          if (!ok) return;
          await authService.logout();
        }}
      >
        Close
      </button>
    </div>
  );
};
