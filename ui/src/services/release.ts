import { pb } from "./client/pb";

export const RELEASES_COLLECTION = "releases";
export const SERVICES_COLLECTION = "services";

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
            version: r.version
        }));
    },
    createInstance: async (data: { name: string, release: string, restart_policy: string }) => {
        const services = pb.collection(SERVICES_COLLECTION);
        services.create({
            name: data.name,
            release: data.release,
            restart_policy: data.restart_policy
        })
    }
};