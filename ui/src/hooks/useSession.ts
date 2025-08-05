import { useEffect, useRef, useState } from "react";
import { pb } from "../services/client/pb";
import type { AuthRecord } from "pocketbase";
import { authService } from "../services/auth";

export function useSession() {
  const [user, setUser] = useState<AuthRecord | null>(() => {
    return pb.authStore.isValid ? pb.authStore.record : null;
  });
  const ensure_once = useRef(true);
  useEffect(() => {
    const unsubscribe = pb.authStore.onChange(() => {
      setTimeout(() => setUser(pb.authStore.record), 0);
    });
    if (ensure_once.current) {
      if (pb.authStore.isValid) authService.refresh();
      ensure_once.current = false;
    }
    if (pb.authStore.isValid) setUser(pb.authStore.record);
    return () => unsubscribe();
  }, []);

  return { user };
}
