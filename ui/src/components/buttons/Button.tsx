import type { ButtonHTMLAttributes } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  loading?: boolean;
  label: string;
}

export const Button = ({
  loading = false,
  label,
  disabled,
  ...props
}: ButtonProps) => {
  return (
    <button
      {...props}
      disabled={disabled || loading}
      className="btn btn-primary w-full flex items-center justify-center gap-2"
    >
      {loading && <span className="loading loading-spinner loading-sm" />}
      <span>{label}</span>
    </button>
  );
};
