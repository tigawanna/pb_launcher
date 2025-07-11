import { useEffect, useRef, useState } from "react";
import { pb } from "../services/client/pb";
import type { AuthRecord } from "pocketbase";
import { authService } from "../services/auth";

export function useSession() {
  const [user, setUser] = useState<AuthRecord | null>(() => {
    return pb.authStore.record;
  });
  const ensure_once = useRef(true);
  useEffect(() => {
    if (ensure_once.current) {
      if (pb.authStore.isValid) authService.refresh();
      ensure_once.current = false;
    }
    const unsubscribe = pb.authStore.onChange(() => {
      setTimeout(() => setUser(pb.authStore.record), 0);
    });
    if (pb.authStore.isValid) setUser(pb.authStore.record);
    return () => unsubscribe();
  }, []);

  return { user };
}
