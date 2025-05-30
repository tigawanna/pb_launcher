import PocketBase from "pocketbase";
const api_url = import.meta.env.VITE_API_URL ?? "/";
export const pb = new PocketBase(api_url);
