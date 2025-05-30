import { HttpError } from "../services/client/errors";

export function getErrorMessage(error: unknown): string {
  if (typeof error === "string") return error;
  if (error instanceof HttpError) return error.message;
  if (error instanceof DOMException && error.name === "AbortError")
    return "Request was cancelled.";
  if (error instanceof Error) return error.message;
  return "Failed to connect to server. Please check your connection.";
}
