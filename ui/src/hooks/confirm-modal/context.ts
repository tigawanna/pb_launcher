import { createContext } from "react";

interface ConfirmModalContextType {
  openModal: (title: string, message: string) => Promise<boolean>;
}

export const ConfirmModalContext = createContext<
  ConfirmModalContextType | undefined
>(undefined);
