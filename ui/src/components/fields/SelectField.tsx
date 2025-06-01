import type { SelectHTMLAttributes } from "react";
import type { FieldError, UseFormRegisterReturn } from "react-hook-form";
import { RotateCw } from "lucide-react";

export interface SelectFieldOption {
  label: string;
  value: string;
}

interface SelectFieldProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  options: SelectFieldOption[];
  error?: FieldError;
  placeholderOptionLabel?: string;
  registration: UseFormRegisterReturn;
  isLoading?: boolean;
  onReload?: () => void;
}

export const SelectField = ({
  label,
  options,
  error,
  registration,
  placeholderOptionLabel = "Select an option",
  isLoading = false,
  onReload,
  ...props
}: SelectFieldProps) => {
  return (
    <div className="flex flex-col space-y-1">
      <div className="flex items-center justify-between">
        <label className="font-medium">{label}</label>
        {onReload && (
          <button
            type="button"
            onClick={onReload}
            className="btn btn-xs btn-ghost text-primary"
            aria-label="Reload options"
          >
            <RotateCw className="w-4 h-4" />
          </button>
        )}
      </div>

      <select
        {...registration}
        {...props}
        className="select select-md select-bordered w-full focus:outline-none focus:ring-1 focus:ring-primary"
        disabled={isLoading || props.disabled}
      >
        {placeholderOptionLabel && (
          <option value="">{placeholderOptionLabel}</option>
        )}
        {options.map(option => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>

      {error && (
        <p className="text-xs text-error mt-1 text-wrap">{error.message}</p>
      )}
    </div>
  );
};
