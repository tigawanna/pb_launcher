import * as yup from "yup";

// Default error messages
const DEFAULT_REQUIRED_MESSAGE = "This field is required";
const DEFAULT_TYPE_BOOLEAN_ERROR = "The value must be a boolean";
const DEFAULT_TYPE_NUMBER_ERROR = "The value must be a number";
const DEFAULT_POSITIVE_NUMBER_ERROR = "The value must be a positive number";
const DEFAULT_EMAIL_ERROR = "Please enter a valid email address";
const DEFAULT_REQUIRED_EMAIL_ERROR = "Email is required";
const DEFAULT_REGEX_ERROR = "Invalid value format";
const DEFAULT_PHONE_ERROR = "Please enter a valid phone number";

// String sanitizer
export const sanitizeString = (value: unknown) => {
  if (value == null) return null;
  if (typeof value !== "string") return value;
  const trimed = value.trim();
  if (trimed === "") return null;
  return trimed;
};

// Required string field
export const stringRequired = (message = DEFAULT_REQUIRED_MESSAGE) =>
  yup.string().transform(sanitizeString).required(message);

// Required boolean field
export const booleanRequired = (message = DEFAULT_REQUIRED_MESSAGE) =>
  yup
    .boolean()
    .transform(sanitizeString)
    .typeError(DEFAULT_TYPE_BOOLEAN_ERROR)
    .required(message);

// Numeric string validation
type NumericStringOptions = {
  minLength?: number;
  maxLength?: number;
};

const numericString = (options?: NumericStringOptions) => {
  let schema = yup
    .string()
    .transform(sanitizeString)
    .test("numeric-only", "Only numeric characters are allowed", value => {
      if (value == null || value === "") return true;
      return /^[0-9]+$/.test(value);
    });

  if (options?.minLength != null) {
    schema = schema.min(
      options.minLength,
      `Minimum length is ${options.minLength} characters`,
    );
  }
  if (options?.maxLength != null) {
    schema = schema.max(
      options.maxLength,
      `Maximum length is ${options.maxLength} characters`,
    );
  }
  return schema;
};

export const numericStringRequired = (options?: NumericStringOptions) =>
  numericString(options).required(DEFAULT_REQUIRED_MESSAGE);

export const numericStringNullable = (options?: NumericStringOptions) =>
  numericString(options).nullable().default(null);

// Nullable string
export const stringNullable = () =>
  yup.string().transform(sanitizeString).nullable().default(null);

// Positive number validations
export const positiveNumberNullable = () =>
  yup
    .number()
    .transform(value => (isNaN(value) ? 0 : value))
    .typeError(DEFAULT_TYPE_NUMBER_ERROR)
    .min(0, DEFAULT_POSITIVE_NUMBER_ERROR)
    .nullable()
    .default(null);

export const positiveNumberRequired = (message = DEFAULT_REQUIRED_MESSAGE) =>
  yup
    .number()
    .typeError(DEFAULT_TYPE_NUMBER_ERROR)
    .positive(DEFAULT_POSITIVE_NUMBER_ERROR)
    .required(message);

// Email validations
export const emailRequired = (message = DEFAULT_REQUIRED_EMAIL_ERROR) =>
  yup.string().email(DEFAULT_EMAIL_ERROR).required(message);

export const emailNullable = (message = DEFAULT_EMAIL_ERROR) =>
  yup
    .string()
    .transform(sanitizeString)
    .test("valid-email", message, value => {
      if (value == null || value === "") return true;
      return yup.string().email().isValidSync(value);
    })
    .nullable()
    .default(null);

// Regex-matching strings
export const stringNullableMatching = (
  pattern: RegExp,
  message = DEFAULT_REGEX_ERROR,
) =>
  yup
    .string()
    .transform(sanitizeString)
    .test("match-pattern", message, value => {
      if (value == null || value === "") return true;
      return pattern.test(value);
    })
    .nullable()
    .default(null);

export const stringRequiredMatching = (
  pattern: RegExp,
  message = DEFAULT_REGEX_ERROR,
) =>
  yup
    .string()
    .transform(sanitizeString)
    .trim()
    .test("match-pattern", message, value => {
      if (value == null || value === "") return true;
      return pattern.test(value);
    })
    .required(DEFAULT_REQUIRED_MESSAGE);

// Phone number validations (Peru format)
const PERU_PHONE_REGEX = /^(?:\+51)\s?[0-9]{3}\s?[0-9]{3}\s?[0-9]{3}$/;

const transformPhoneNumber = (value: unknown) => {
  if (value == null) return null;
  if (typeof value !== "string") return value;
  const trimmed = value.trim().replace(/\s+/g, "");
  if (trimmed === "") return null;
  return trimmed.startsWith("+") ? trimmed : `+51 ${trimmed}`;
};

export const phoneNullable = (message = DEFAULT_PHONE_ERROR) =>
  yup
    .string()
    .transform(transformPhoneNumber)
    .test("valid-phone", message, value => {
      if (value == null || value === "") return true;
      return PERU_PHONE_REGEX.test(value);
    })
    .nullable()
    .default(null);

export const phoneRequired = (message = DEFAULT_PHONE_ERROR) =>
  yup
    .string()
    .transform(transformPhoneNumber)
    .test("valid-phone", message, value => {
      if (value == null || value === "") return true;
      return PERU_PHONE_REGEX.test(value);
    })
    .required(DEFAULT_REQUIRED_MESSAGE);

// Password validation
export const passwordRequired = (message = "Password is required") =>
  yup
    .string()
    .transform(sanitizeString)
    .min(8, "Password must be at least 8 characters")
    .matches(/[A-Z]/, "Password must contain at least one uppercase letter")
    .matches(/[a-z]/, "Password must contain at least one lowercase letter")
    .matches(/[0-9]/, "Password must contain at least one number")
    .required(message);

// Enum validation (string)
export const stringEnumRequired = (
  values: string[],
  message = DEFAULT_REQUIRED_MESSAGE,
) =>
  yup
    .string()
    .oneOf(values, `Value must be one of: ${values.join(", ")}`)
    .required(message);

// Date validation
export const dateRequired = (message = "Date is required") =>
  yup.date().typeError("Invalid date format").required(message);

export const dateNullable = () =>
  yup.date().typeError("Invalid date format").nullable().default(null);

// Convert numeric string to number
export const numericStringToNumber = () =>
  yup
    .string()
    .transform(value =>
      value == null || value === "" ? null : parseFloat(value),
    )
    .typeError("Must be a valid number")
    .nullable()
    .default(null);

// Example: login schema using utilities
export const loginSchema = yup.object({
  email: emailRequired(),
  password: passwordRequired(),
});

const DEFAULT_URL_ERROR = "Invalid URL";
const REQUIRED_URL_ERROR = "URL is required";

export const urlNullable = (message = DEFAULT_URL_ERROR) =>
  yup
    .string()
    .transform(sanitizeString)
    .test("valid-url", message, value => {
      if (value == null || value === "") return true;
      try {
        const url = new URL(value);
        return url.protocol === "http:" || url.protocol === "https:";
      } catch {
        return false;
      }
    })
    .nullable()
    .default(null);

export const urlRequired = (
  message = DEFAULT_URL_ERROR,
  required = REQUIRED_URL_ERROR,
) =>
  yup
    .string()
    .transform(sanitizeString)
    .test("valid-url", message, value => {
      if (value == null || value === "") return true;
      try {
        const url = new URL(value);
        return url.protocol === "http:" || url.protocol === "https:";
      } catch {
        return false;
      }
    })
    .required(required);
