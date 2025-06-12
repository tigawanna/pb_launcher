import { use100vh } from "react-div-100vh";

export const useViewportHeight = (): number => {
  const height = use100vh();
  return height ?? window.innerHeight;
};
