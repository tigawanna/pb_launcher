import classNames from "classnames";
import type {
  InputHTMLAttributes,
  HTMLInputAutoCompleteAttribute,
} from "react";
import type { FieldError, UseFormRegisterReturn } from "react-hook-form";

interface InputFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  placeholder?: string;
  error?: FieldError;
  registration: UseFormRegisterReturn;
  autoComplete?: HTMLInputAutoCompleteAttribute;
}

export const InputField = ({
  label,
  error,
  registration,
  placeholder,
  autoComplete,
  type,
  ...props
}: InputFieldProps) => {
  const isCheckbox = type === "checkbox";

  return (
    <div
      className={classNames("form-control", {
        "w-full": !isCheckbox,
        "flex-row items-center gap-2": isCheckbox,
      })}
    >
      <label className="label cursor-pointer">
        {isCheckbox && (
          <input
            type="checkbox"
            {...registration}
            {...props}
            className={classNames("checkbox", { "checkbox-error": error })}
          />
        )}
        <span className={`label-text ${isCheckbox ? "" : "mb-1"}`}>
          {label}
        </span>
      </label>
      {!isCheckbox && (
        <input
          {...registration}
          {...props}
          type={type}
          autoComplete={autoComplete ?? "off"}
          placeholder={placeholder}
          className={classNames(
            "input input-md",
            "input-bordered",
            "w-full",
            "transition-all duration-150",
            "focus:outline-none focus:ring-1 focus:ring-primary",
            { "input-error": error },
          )}
        />
      )}
      {error && (
        <span className="label-text-alt text-xs text-error">
          {error.message}
        </span>
      )}
    </div>
  );
};
