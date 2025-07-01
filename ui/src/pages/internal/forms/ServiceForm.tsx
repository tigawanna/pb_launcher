import { object } from "yup";
import { stringRequired } from "../../../utils/validation";
import { useCustomForm } from "../../../hooks/useCustomForm";
import { InputField } from "../../../components/fields/InputField";
import { Button } from "../../../components/buttons/Button";
import {
  SelectField,
  type SelectFieldOption,
} from "../../../components/fields/SelectField";
import { serviceService, type ServiceDto } from "../../../services/services";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useMemo, type FC } from "react";
import { useModal } from "../../../components/modal/hook";
import toast from "react-hot-toast";
import { getErrorMessage } from "../../../utils/errors";
import classNames from "classnames";
import { releaseService } from "../../../services/release";

const schema = object({
  name: stringRequired(), // Name of the new PocketBase instance
  instanceSource: stringRequired(), // Source for the instance (template, version, etc.)
  restartPolicy: stringRequired(), // Restart policy: "no" or "on-failure"
});

type Props = {
  record?: ServiceDto;
  onSaveRecord?: () => void;
  width?: number;
};

export const ServiceForm: FC<Props> = ({ onSaveRecord, record, width }) => {
  const { closeModal } = useModal();
  const form = useCustomForm(schema, {
    defaultValues: {
      name: record?.name,
      instanceSource: record?.release_id,
      restartPolicy: record?.restart_policy ?? "on-failure",
    },
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
    mutationFn: serviceService.createServiceInstance,
    onSuccess: () => {
      toast.success("Service created successfully");
      closeModal();
      onSaveRecord?.();
    },
    onError: error => toast.error(getErrorMessage(error)),
  });

  const updateInstanceMutation = useMutation({
    mutationFn: serviceService.updateServiceInstance,
    onSuccess: () => {
      toast.success("Service updated successfully");
      closeModal();
      onSaveRecord?.();
    },
    onError: error => toast.error(getErrorMessage(error)),
  });

  const handleFormSubmit = form.handleSubmit(
    ({ instanceSource, name, restartPolicy }) => {
      if (record == null)
        createInstanceMutation.mutate({
          name,
          release: instanceSource,
          restart_policy: restartPolicy,
        });
      else
        updateInstanceMutation.mutate({
          id: record.id,
          name,
          release: instanceSource,
          restart_policy: restartPolicy,
        });
    },
  );
  return (
    <div style={{ width: width }}>
      <form onSubmit={handleFormSubmit} className="space-y-5">
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
          disabled={record != null}
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
        <div
          className={classNames("mt-8", {
            "flex justify-end": width == null || width > 400,
          })}
        >
          <div
            className={classNames("form-control", {
              "w-[200px]": width == null || width > 400,
            })}
          >
            <Button type="submit" label="Guardar" loading={false} />
          </div>
        </div>
      </form>
    </div>
  );
};
