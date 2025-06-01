import { object } from "yup";
import { stringRequired } from "../../../utils/validation";
import { useCustomForm } from "../../../hooks/useCustomForm";
import { InputField } from "../../../components/fields/InputField";
import { Button } from "../../../components/buttons/Button";
import {
  SelectField,
  type SelectFieldOption,
} from "../../../components/fields/SelectField";
import { releaseService } from "../../../services/release";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useMemo, type FC } from "react";
import { useModal } from "../../../components/modal/hook";
import toast from "react-hot-toast";
import { getErrorMessage } from "../../../utils/errors";

const schema = object({
  name: stringRequired(), // Name of the new PocketBase instance
  instanceSource: stringRequired(), // Source for the instance (template, version, etc.)
  restartPolicy: stringRequired(), // Restart policy: "no" or "on-failure"
});

type Props = {
  onSaveRecord?: () => void;
};

export const ServiceForm: FC<Props> = ({ onSaveRecord }) => {
  const { closeModal } = useModal();
  const form = useCustomForm(schema, {
    defaultValues: { restartPolicy: "on-failure" },
  });

  const releasesQuery = useQuery({
    queryKey: ["releases"],
    queryFn: releaseService.fetchAll,
  });

  const releaseOptions = useMemo<SelectFieldOption[]>(() => {
    return (
      releasesQuery.data?.map(r => ({
        label: `${r.repositoryName} v${r.version}`,
        value: r.id,
      })) ?? []
    );
  }, [releasesQuery.data]);

  const createInstanceMutation = useMutation({
    mutationFn: releaseService.createInstance,
    onSuccess: () => {
      toast.success("Service created successfully");
      onSaveRecord?.();
      closeModal();
    },
    onError: error => toast.error(getErrorMessage(error)),
  });

  const handleFormSubmit = form.handleSubmit(
    ({ instanceSource, name, restartPolicy }) => {
      createInstanceMutation.mutate({
        name,
        release: instanceSource,
        restart_policy: restartPolicy,
      });
    },
  );

  return (
    <div className="w-[360px]">
      <form onSubmit={handleFormSubmit} className="space-y-4">
        <InputField
          label="Instance Name"
          registration={form.register("name")}
          autoComplete="off"
          error={form.formState.errors.name}
        />

        <SelectField
          label="Source/Version"
          options={releaseOptions}
          isLoading={releasesQuery.isLoading}
          onReload={releasesQuery.refetch}
          registration={form.register("instanceSource")}
          autoComplete="off"
          error={form.formState.errors.instanceSource}
        />

        <SelectField
          label="Restart Policy"
          options={[
            { label: "No", value: "no" },
            { label: "On Failure", value: "on-failure" },
          ]}
          registration={form.register("restartPolicy")}
          autoComplete="off"
          error={form.formState.errors.restartPolicy}
        />

        <div className="form-control mt-6">
          <Button type="submit" label="Guardar" loading={false} />
        </div>
      </form>
    </div>
  );
};
