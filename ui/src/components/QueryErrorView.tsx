import { getErrorMessage } from "../utils/errors";

interface QueryErrorViewProps {
  error: unknown;
  onRetry: () => void;
}

export const QueryErrorView: React.FC<QueryErrorViewProps> = ({
  error,
  onRetry,
}) => {
  return (
    <div className="flex flex-col items-center justify-center h-screen space-y-4 px-4 text-center">
      <div className="text-red-600 text-lg font-medium max-w-md">
        Error: {getErrorMessage(error)}
      </div>
      <button
        onClick={onRetry}
        className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
      >
        Retry
      </button>
    </div>
  );
};
