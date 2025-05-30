import { useForm, type DefaultValues } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import type { ObjectSchema, InferType } from "yup";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const useCustomForm = <TSchema extends ObjectSchema<any>>(
  schema: TSchema,
  props?: {
    defaultValues?:
      | DefaultValues<InferType<TSchema>>
      | (() => Promise<DefaultValues<InferType<TSchema>>>);
  },
) => {
  const form = useForm<InferType<TSchema>>({
    resolver: yupResolver(schema),
    defaultValues: props?.defaultValues,
  });

  const wrappedHandleSubmit: typeof form.handleSubmit = (
    onValid,
    onInvalid,
  ) => {
    return form.handleSubmit(onValid, errors => {
      if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
        console.error("Validation Errors:", errors);
      }
      if (onInvalid) {
        onInvalid(errors);
      }
    });
  };
  return { ...form, handleSubmit: wrappedHandleSubmit };
};
