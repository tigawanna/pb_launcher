import { useEffect, useState } from "react";
import { pb } from "../services/client/pb";
import type { AuthRecord } from "pocketbase";

export function useSession() {
  const [user, setUser] = useState<AuthRecord | null>(() => {
    return pb.authStore.record;
  });

  useEffect(() => {
    const unsubscribe = pb.authStore.onChange(() =>
      setUser(pb.authStore.record),
    );
    setUser(pb.authStore.record);
    return () => unsubscribe();
  }, []);

  return { user };
}
