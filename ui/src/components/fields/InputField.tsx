import type { InputHTMLAttributes, HTMLInputAutoCompleteAttribute } from "react";
import type { FieldError, UseFormRegisterReturn } from "react-hook-form";

interface InputFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  placeholder?: string;
  error?: FieldError;
  registration: UseFormRegisterReturn;
  autoComplete?: HTMLInputAutoCompleteAttribute;
}

export const InputField = ({ label, error, registration, placeholder, autoComplete, ...props }: InputFieldProps) => {
  return (
    <div className="form-control w-full">
      <label className="label">
        <span className="label-text">{label}</span>
      </label>
      <input
        {...registration}
        {...props}
        autoComplete={autoComplete ?? "off"}
        placeholder={placeholder}
        className={`input input-bordered w-full transition-all duration-150 focus:outline-none focus:ring-1 focus:ring-primary ${
          error ? "input-error" : ""
        }`}
      />
      {error && (
        <label className="label">
          <span className="label-text-alt text-error">{error.message}</span>
        </label>
      )}
    </div>
  );
};
