import toast from "react-hot-toast";
import { useMemo, useState } from "react";
import { Plus, RefreshCcw } from "lucide-react";
import { useModal } from "../../components/modal/hook";
import { ServiceForm } from "./forms/ServiceForm";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useLocalStorage } from "@uidotdev/usehooks";
import { releaseService, type ServiceDto } from "../../services/release";
import { ServiceCard } from "./components/ServiceCard";
import { useConfirmModal } from "../../hooks/useConfirmModal";
import { getErrorMessage } from "../../utils/errors";
import classNames from "classnames";

const STATUS_FILTER_KEY = "pb-dashboard-status-filter";
type TStatus = "all" | "running" | "stopped";
export const ServicesPage = () => {
  const { openModal } = useModal();
  const confirm = useConfirmModal();

  const servicesQuery = useQuery({
    queryKey: ["services"],
    queryFn: releaseService.fetchAllServices,
    refetchInterval: 5000,
  });

  const [query, setQuery] = useState("");
  const [statusFilter, setStatusFilter] = useLocalStorage<{ value: TStatus }>(
    STATUS_FILTER_KEY,
    { value: "all" },
  );

  const filtered = useMemo(() => {
    return (servicesQuery.data ?? [])
      .filter(s => s.name.toLowerCase().includes(query.toLowerCase()))
      .filter(s => {
        switch (statusFilter.value) {
          case "all":
            return true;
          case "running":
            return s.status === "running";
          case "stopped":
            return s.status === "stopped";
        }
      });
  }, [servicesQuery.data, query, statusFilter]);

  const deleteMutation = useMutation({
    mutationFn: releaseService.deleteServiceInstance,
    onSuccess: () => servicesQuery.refetch(),
    onError: error => toast.error(getErrorMessage(error)),
  });

  const handleDelete = async (id: string) => {
    const ok = await confirm(
      "Delete service",
      "Are you sure you want to delete this service?",
    );
    if (ok) {
      deleteMutation.mutate(id);
    }
  };

  const openCreateServiceModal = () => {
    openModal(
      <ServiceForm onSaveRecord={() => setTimeout(servicesQuery.refetch)} />,
      {
        title: "Create Service",
      },
    );
  };

  const openEditServiceModal = (service: ServiceDto) => {
    openModal(
      <ServiceForm
        record={service}
        onSaveRecord={() => setTimeout(servicesQuery.refetch)}
      />,
      {
        title: "Edit Service",
      },
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex gap-2 w-full sm:max-w-md">
          <input
            type="text"
            placeholder="Search service..."
            className="input input-sm input-bordered w-full"
            value={query}
            onChange={e => setQuery(e.target.value)}
          />
          <button
            onClick={() => servicesQuery.refetch()}
            className="btn btn-sm btn-ghost"
          >
            <RefreshCcw
              className={classNames("w-4 h-4", {
                "animate-spin": servicesQuery.isFetching,
              })}
            />
          </button>
        </div>

        <div className="flex flex-col md:flex-row gap-4 select-none">
          <select
            className="select select-sm select-bordered w-full sm:w-60"
            value={statusFilter.value}
            onChange={e =>
              setStatusFilter({ value: e.target.value as TStatus })
            }
          >
            <option value="all">All</option>
            <option value="running">Running</option>
            <option value="stopped">Stopped</option>
          </select>
          <button
            className="btn btn-sm btn-primary gap-2 w-full sm:w-auto"
            onClick={openCreateServiceModal}
          >
            <Plus className="w-4 h-4" />
            New instance
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {filtered.map(service => (
          <ServiceCard
            key={service.id}
            service={service}
            onEdit={() => openEditServiceModal(service)}
            onDelete={() => handleDelete(service.id)}
          />
        ))}
      </div>
    </div>
  );
};
