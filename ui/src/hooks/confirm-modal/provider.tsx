import { ConfirmModalContext } from "./context";
import { useState, type ReactNode } from "react";

export const ConfirmModalProvider = ({ children }: { children: ReactNode }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [resolve, setResolve] = useState<(value: boolean) => void>(() => {});
  const [title, setTitle] = useState<string>("Confirmation");
  const [message, setMessage] = useState<string>(
    "Are you sure you want to do this?",
  );

  const openModal = (title: string, message: string): Promise<boolean> => {
    setIsOpen(true);
    setTitle(title);
    setMessage(message);
    return new Promise(res => setResolve(() => res));
  };

  const closeModal = () => {
    setIsOpen(false);
  };

  const handleConfirm = () => {
    resolve(true);
    closeModal();
  };

  const handleCancel = () => {
    resolve(false);
    closeModal();
  };

  const Modal = () => {
    if (!isOpen) return null;

    return (
      <div className="fixed inset-0 flex items-center justify-center bg-black/50 z-[9999]">
        <div className="bg-base-100 p-6 rounded-lg shadow-lg w-full max-w-sm">
          <h3 className="text-lg font-bold mb-2">{title}</h3>
          <p className="text-base-content/80 mb-6">{message}</p>
          <div className="flex justify-end gap-4">
            <button className="btn btn-ghost" onClick={handleCancel}>
              Cancel
            </button>
            <button className="btn btn-primary" onClick={handleConfirm}>
              Confirm
            </button>
          </div>
        </div>
      </div>
    );
  };

  return (
    <ConfirmModalContext.Provider value={{ openModal }}>
      <>
        <Modal />
        {children}
      </>
    </ConfirmModalContext.Provider>
  );
};
