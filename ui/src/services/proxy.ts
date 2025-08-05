import { pb } from "./client/pb";
import { domainsService, type DomainDto } from "./services_domain";

const PROXY_ENTRIES = "proxy_entries";

export interface ProxyEntryDto {
  id: string;
  name: string;
  target_url: string;
  enabled: "yes" | "no";
  domains?: DomainDto[];
}

export const proxyEntryService = {
  fetchById: async (id: string) => {
    const proxyEntries = pb.collection(PROXY_ENTRIES);
    const record = await proxyEntries.getOne<ProxyEntryDto>(id, {
      filter: `deleted=""`,
    });
    return record;
  },

  fetchAll: async () => {
    const [records, domains] = await Promise.all([
      pb.collection(PROXY_ENTRIES).getFullList<ProxyEntryDto>({
        filter: `deleted=""`,
      }),
      domainsService.fetchFullList(),
    ]);
    return records.map(
      (entry): ProxyEntryDto => ({
        ...entry,
        domains: domains.filter(
          d => d.proxy_entry === entry.id && d.x_has_valid_ssl_cert,
        ),
      }),
    );
  },
  delete: async (id: string) => {
    const proxyEntries = pb.collection(PROXY_ENTRIES);
    await proxyEntries.update(id, { deleted: new Date().toJSON() });
  },
  create: async (data: { name: string; target_url: string }) => {
    const proxyEntries = pb.collection(PROXY_ENTRIES);
    await proxyEntries.create({ ...data, enabled: "yes" });
  },
  update: async (data: {
    id: string;
    name: string;
    target_url: string;
    enabled: string;
  }) => {
    const proxyEntries = pb.collection(PROXY_ENTRIES);
    await proxyEntries.update(data.id, data);
  },
};
