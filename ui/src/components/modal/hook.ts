import { useContext } from "react";
import { ModalContext, type ModalContextType } from "./context";

export const useModal = (): ModalContextType => {
  const context = useContext(ModalContext);
  if (!context) {
    throw new Error("useModal debe usarse dentro de ModalProvider");
  }
  return context;
};
