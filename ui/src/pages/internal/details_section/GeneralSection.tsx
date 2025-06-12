import { useQuery } from "@tanstack/react-query";
import type { FC } from "react";
import { releaseService } from "../../../services/release";
import { ServiceForm } from "../forms/ServiceForm";
import { Navigate } from "react-router-dom";
import classNames from "classnames";

type Props = {
  service_id: string;
};

export const GeneralSection: FC<Props> = ({ service_id }) => {
  const serviceQuery = useQuery({
    queryKey: ["services", service_id],
    queryFn: () => releaseService.fetchServiceByID(service_id),
    refetchOnMount: true,
  });

  if (serviceQuery.isLoading) {
    return <div className="p-4">Loading...</div>;
  }

  if (serviceQuery.isError) {
    return (
      <div className="p-4">
        <p className="text-error">Failed to load service.</p>
        <button
          className="btn btn-sm btn-outline mt-2"
          onClick={() => serviceQuery.refetch()}
        >
          Retry
        </button>
      </div>
    );
  }
  if (serviceQuery.data == null) return <Navigate to={"/"} />;
  const service = serviceQuery.data;
  const status =
    service.status.charAt(0).toUpperCase() + service.status.slice(1);
  return (
    <div className="relative pt-4">
      <ServiceForm
        record={service}
        onSaveRecord={() => setTimeout(() => serviceQuery.refetch())}
      />
      <p
        className={classNames("badge badge-sm absolute -top-2 right-2", {
          "badge-success": service.status === "running",
          "badge-warning":
            service.status === "pending" || service.status === "idle",
          "badge-error": service.status === "failure",
          "badge-neutral": !["running", "pending", "idle", "failure"].includes(
            service.status,
          ),
        })}
      >
        {status}
      </p>
    </div>
  );
};
