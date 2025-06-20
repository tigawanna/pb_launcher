import { useMemo, type FC } from "react";
import { DomainCard } from "../components/DomainCard";
import { useProxyConfigs } from "../../../hooks/useProxyConfigs";
import { formatUrl } from "../../../utils/url";
import { Plus } from "lucide-react";
import { useModal } from "../../../components/modal/hook";
import { DomainForm } from "../forms/DomainForm";
import {
  domainsService,
  type DomainDto,
} from "../../../services/services_domain";
import { useMutation, useQuery } from "@tanstack/react-query";
import { ErrorFallback } from "../../../components/helpers/ErrorFallback";
import { getErrorMessage } from "../../../utils/errors";
import toast from "react-hot-toast";
import { useConfirmModal } from "../../../hooks/useConfirmModal";

type Props = {
  service_id: string;
};

export const DomainsSection: FC<Props> = ({ service_id }) => {
  const confirm = useConfirmModal();
  const { openModal } = useModal();
  const proxy = useProxyConfigs();
  const domainsQuery = useQuery({
    queryKey: ["services", service_id, "domains"],
    queryFn: () => domainsService.fetchAll(service_id),
    refetchOnMount: true,
  });

  const proxyDomain = useMemo((): DomainDto => {
    return {
      id: "__",
      domain: proxy.base_domain ? `${service_id}.${proxy.base_domain}` : "--",
      status: "active",
      use_https: proxy.use_https ? "yes" : "no",
    };
  }, [proxy.base_domain, proxy.use_https, service_id]);

  const openCreateModal = () => {
    openModal(
      <DomainForm
        service_id={service_id}
        onSaveRecord={() => setTimeout(domainsQuery.refetch)}
        width={360}
      />,
      {
        title: "Create Domain",
      },
    );
  };

  const openEditModal = (record: DomainDto) => {
    openModal(
      <DomainForm
        service_id={service_id}
        width={360}
        record={record}
        onSaveRecord={() => setTimeout(domainsQuery.refetch)}
      />,
      {
        title: "Edit Domain",
      },
    );
  };

  const deleteMutation = useMutation({
    mutationFn: domainsService.deleteDomain,
    onSuccess: () => setTimeout(() => domainsQuery.refetch()),
    onError: error => toast.error(getErrorMessage(error)),
  });

  const handleDelete = async (id: string) => {
    const ok = await confirm(
      "Delete domain",
      "Are you sure you want to delete this domain?",
    );
    if (ok) {
      deleteMutation.mutate(id);
    }
  };

  if (domainsQuery.isFetching) {
    return <div className="p-4">Loading...</div>;
  }

  if (domainsQuery.isError)
    return (
      <ErrorFallback
        error={domainsQuery.error}
        onRetry={() => setTimeout(domainsQuery.refetch)}
      />
    );

  return (
    <div className="space-y-6">
      <div className="flex justify-end">
        <div className="flex gap-2">
          <button
            className="btn btn-sm btn-primary gap-2 w-full sm:w-auto"
            onClick={openCreateModal}
          >
            <Plus className="w-4 h-4" />
            New instance
          </button>
        </div>
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        <DomainCard
          readonly
          url={formatUrl(
            proxyDomain.use_https === "yes" ? "https" : "http",
            proxyDomain.domain,
            proxy.use_https ? proxy.https_port : proxy.http_port,
          )}
          port={proxy.use_https ? proxy.https_port : proxy.http_port}
          domain={proxyDomain}
        />
        {(domainsQuery.data ?? []).map(domain => (
          <DomainCard
            key={domain.id}
            domain={domain}
            onEdit={() => openEditModal(domain)}
            onDelete={() => handleDelete(domain.id)}
          />
        ))}
      </div>
    </div>
  );
};
