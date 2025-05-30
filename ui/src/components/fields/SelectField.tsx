import type { SelectHTMLAttributes } from "react";
import type { FieldError, UseFormRegisterReturn } from "react-hook-form";

interface Option {
  label: string;
  value: string;
}

interface SelectFieldProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  options: Option[];
  error?: FieldError;
  registration: UseFormRegisterReturn;
}

export const SelectField = ({ label, options, error, registration, ...props }: SelectFieldProps) => {
  return (
    <div className="flex flex-col space-y-1">
      <label className="font-medium">{label}</label>
      <select {...registration} {...props} className="select select-bordered w-full">
        <option value="">Selecciona una opción</option>
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      {error && <p className="text-xs text-error mt-1">{error.message}</p>}
    </div>
  );
};
