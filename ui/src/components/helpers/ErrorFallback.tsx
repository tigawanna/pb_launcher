import type { FC } from "react";
import { getErrorMessage } from "../../utils/errors";

interface ErrorFallbackProps {
  error?: string | Error | null;
  onRetry: () => void;
}

export const ErrorFallback: FC<ErrorFallbackProps> = ({
  error = "Failed to load service.",
  onRetry,
}) => {
  return (
    <div className="p-4">
      <p className="text-error">{getErrorMessage(error)}</p>
      <button className="btn btn-sm btn-outline mt-2" onClick={onRetry}>
        Retry
      </button>
    </div>
  );
};
