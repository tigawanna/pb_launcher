import { createContext, type ReactNode, type CSSProperties } from "react";

export type ModalComponentContent = ReactNode | ((index: number) => ReactNode);

export type ModalComponent = {
  title?: string;
  width?: CSSProperties["width"];
  height?: CSSProperties["height"];
  zIndex?: CSSProperties["zIndex"];
  closeOnBackdropClick?: boolean;
  disableCloseButton?: boolean;
  //
  content?: ModalComponentContent;
};

export interface ModalContextType {
  openModal: (
    modalContent: ModalComponentContent,
    props?: Omit<ModalComponent, "content">,
  ) => void;
  closeModal: () => void;
  closeAllModals: () => void;
}

export const ModalContext = createContext<ModalContextType | undefined>(
  undefined,
);
