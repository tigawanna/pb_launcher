import { pb } from "./client/pb";
const DOMAINS_COLLECTION = "services_domains";
const CERT_REQUESTS = "cert_requests";

export interface DomainDto {
  id: string;
  use_https: "yes" | "no";
  domain: string;
  x_cert_request_state?: "pending" | "approved" | "failed";
  x_has_valid_ssl_cert?: boolean;
  x_reached_max_attempt?: boolean;
  x_failed_error_message?: string;
}

export const domainsService = {
  fetchAll: async (service_id: string) => {
    const domains = pb.collection(DOMAINS_COLLECTION);
    const records = await domains.getFullList<DomainDto>({
      filter: `service="${service_id}"`,
    });
    return records;
  },

  createDomain: async (data: {
    use_https: boolean;
    domain: string;
    service: string;
  }) => {
    const services = pb.collection(DOMAINS_COLLECTION);
    await services.create({
      service: data.service,
      domain: data.domain,
      use_https: data.use_https ? "yes" : "no",
    });
  },

  updateDomain: async (data: { id: string; use_https: boolean }) => {
    const services = pb.collection(DOMAINS_COLLECTION);
    await services.update(data.id, {
      use_https: data.use_https ? "yes" : "no",
    });
  },

  deleteDomain: async (id: string) => {
    const services = pb.collection(DOMAINS_COLLECTION);
    await services.delete(id);
  },

  createSSLRequest: async (domain: string) => {
    const certRequest = pb.collection(CERT_REQUESTS);
    await certRequest.create({ domain });
  },
};
