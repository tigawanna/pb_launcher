import { useQuery } from "@tanstack/react-query";
import type { FC } from "react";
import { Navigate } from "react-router-dom";
import { ErrorFallback } from "../../../components/helpers/ErrorFallback";
import { proxyEntryService } from "../../../services/proxy";
import { ProxyEntryForm } from "../forms/ProxyEntryForm";

type Props = {
  proxy_id: string;
  onChange: () => void;
};

export const ProxyGeneralSection: FC<Props> = ({ proxy_id, onChange }) => {
  const proxyMapQuery = useQuery({
    queryKey: ["proxy_entires", proxy_id],
    queryFn: () => proxyEntryService.fetchById(proxy_id),
    refetchOnMount: true,
  });

  if (proxyMapQuery.isFetching) {
    return <div className="p-4">Loading...</div>;
  }

  if (proxyMapQuery.isError)
    return (
      <ErrorFallback
        error={proxyMapQuery.error}
        onRetry={() => setTimeout(proxyMapQuery.refetch)}
      />
    );

  if (proxyMapQuery.data == null) return <Navigate to={"/proxy"} />;

  return (
    <div className="relative pt-4">
      <ProxyEntryForm
        record={proxyMapQuery.data}
        onSaveRecord={() => {
          setTimeout(() => proxyMapQuery.refetch());
          onChange();
        }}
      />
    </div>
  );
};
