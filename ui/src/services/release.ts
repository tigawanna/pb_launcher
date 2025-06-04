import { pb } from "./client/pb";
export const base_url = import.meta.env.VITE_BASE_URL ?? "/";

export const RELEASES_COLLECTION = "releases";
export const SERVICES_COLLECTION = "services";
export const COMANDS_COLLECTION = "comands";

interface ReleaseDto {
  id: string;
  version: string;
  expand: {
    repository: {
      id: string;
      name: string;
    };
  };
}

interface _Service {
  id: string;
  name: string;
  status: "idle" | "pending" | "running" | "stopped" | "failure";

  url?: string;

  boot_user_email: string;
  boot_user_password: string;
  last_started: string;

  restart_policy: string;
  error_message: string;

  created: string;

  repository: string;
  release_id: string;
  release_version: string;

  // expand
  expand: {
    release: {
      id: string;
      version: string;
      expand: {
        repository: {
          name: string;
        };
      };
    };
  };
}

export type ServiceDto = Omit<_Service, "expand">;

export const releaseService = {
  fetchAll: async () => {
    const releases = pb.collection(RELEASES_COLLECTION);
    const records = await releases.getFullList<ReleaseDto>({
      expand: "repository",
      fields: "id,version,expand.repository.id,expand.repository.name",
      sort: "repository,-version",
    });
    return records.map(r => ({
      id: r.id,
      repositoryId: r.expand.repository.id,
      repositoryName: r.expand.repository.name,
      version: r.version,
    }));
  },
  createServiceInstance: async (data: {
    name: string;
    release: string;
    restart_policy: string;
  }) => {
    const services = pb.collection(SERVICES_COLLECTION);
    services.create({
      name: data.name,
      release: data.release,
      restart_policy: data.restart_policy,
    });
  },
  updateServiceInstance: async (data: {
    id: string;
    name: string;
    release: string;
    restart_policy: string;
  }) => {
    const services = pb.collection(SERVICES_COLLECTION);
    await services.update(data.id, {
      name: data.name,
      release: data.release,
      restart_policy: data.restart_policy,
    });
  },

  buildServiceUrl(id: string): string {
    const url = base_url ? new URL(base_url) : new URL(window.location.href);
    const { protocol, hostname, port } = url;
    const isDefaultPort =
      (protocol === "http:" && port === "80") ||
      (protocol === "https:" && port === "443");
    const portPart = port && !isDefaultPort ? `:${port}` : "";
    return `${protocol}//${id}.${hostname}${portPart}`;
  },

  fetchAllServices: async (): Promise<ServiceDto[]> => {
    const [services, commands] = await Promise.all([
      pb
        .collection(SERVICES_COLLECTION)
        .getFullList<
          Omit<_Service, "repository" | "release_id" | "release_version">
        >({
          filter: `deleted=""`,
          expand: "release.repository",
          fields: [
            "id",
            "name",
            "status",
            "boot_user_email",
            "boot_user_password",
            "last_started",
            "restart_policy",
            "error_message",
            "created",
            "release",
            "expand.release.id",
            "expand.release.version",
            "expand.release.expand.repository.name",
          ].join(","),
        }),
      pb.collection(COMANDS_COLLECTION).getFullList<{ service: string }>({
        fields: "service",
        filter: `status="pending"`,
      }),
    ]);
    const pendingServices = new Set(commands.map(c => c.service));
    return services.map(
      (s): ServiceDto => ({
        id: s.id,
        name: s.name,
        status: pendingServices.has(s.id) ? "pending" : s.status,
        url: releaseService.buildServiceUrl(s.id),
        boot_user_email: s.boot_user_email,
        boot_user_password: s.boot_user_password,
        last_started: s.last_started,
        restart_policy: s.restart_policy,
        error_message: s.error_message,
        created: s.created,
        repository: s.expand.release.expand.repository.name,
        release_id: s.expand.release.id,
        release_version: s.expand.release.version,
      }),
    );
  },
  deleteServiceInstance: async (id: string) => {
    const services = pb.collection(SERVICES_COLLECTION);
    await services.update(id, { deleted: new Date().toJSON() });
  },
  executeServiceCommand: async (data: {
    service_id: string;
    action: "stop" | "start" | "restart";
  }) => {
    const comands = pb.collection(COMANDS_COLLECTION);
    await comands.create({ service: data.service_id, action: data.action });
  },
};
