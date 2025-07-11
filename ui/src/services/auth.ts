import { joinUrls } from "../utils/url";
import { HttpError } from "./client/errors";
import { pb } from "./client/pb";

export const USERS_COLLECTION = "users";

export const authService = {
  isInitialSetupDone: async (signal: AbortSignal): Promise<boolean> => {
    const url = joinUrls(pb.baseURL, "/x-api/setup/admin-exists");
    const response = await fetch(url, { signal });
    const json = await response.json();

    if (!response.ok) {
      throw new HttpError(
        response.status,
        json?.message || "Unexpected error",
        json,
      );
    }

    if (json.message === "yes") return true;
    if (json.message === "no") return false;

    throw new HttpError(500, "Invalid response from server", json);
  },

  setup: async (data: { email: string; password: string }): Promise<void> => {
    const url = joinUrls(pb.baseURL, "/x-api/setup/admin");
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const json = await response.json();
      throw new HttpError(
        response.status,
        json?.message || "Unexpected error",
        json,
      );
    }
  },

  login: async (credentials: { email: string; password: string }) => {
    const users = pb.collection(USERS_COLLECTION);
    await users.authWithPassword(credentials.email, credentials.password);
  },
  refresh: async () => {
    const users = pb.collection(USERS_COLLECTION);
    await users.authRefresh();
  },
  logout: async () => pb.authStore.clear(),
};
