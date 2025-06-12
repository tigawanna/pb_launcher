import { useMemo, type FC } from "react";
import { DomainCard } from "../components/DomainCard";
import { useProxyConfigs } from "../../../hooks/useProxyConfigs";
import { formatUrl } from "../../../utils/url";

type Props = {
  service_id: string;
};

export const DomainsSection: FC<Props> = ({ service_id }) => {
  const proxy = useProxyConfigs();
  const proxyDomain = useMemo(() => {
    return {
      name: proxy.base_domain ? `${service_id}.${proxy.base_domain}` : "--",
      protocol: (proxy.use_https ? "https" : "http") as "https" | "http",
      status: "active" as "active" | "pending",
    };
  }, [proxy.base_domain, proxy.use_https, service_id]);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
      <DomainCard
        readonly
        url={formatUrl(
          proxyDomain.protocol,
          proxyDomain.name,
          proxy.use_https ? proxy.https_port : proxy.http_port,
        )}
        port={proxy.use_https ? proxy.https_port : proxy.http_port}
        domain={proxyDomain}
      />
      {/* {domains.map(d => (
        <DomainCard key={d.name} domain={d} />
      ))} */}
    </div>
  );
};
