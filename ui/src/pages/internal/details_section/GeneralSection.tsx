import { useQuery } from "@tanstack/react-query";
import type { FC } from "react";
import { releaseService } from "../../../services/release";

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

  return <div>{serviceQuery.data?.id}</div>;
};
