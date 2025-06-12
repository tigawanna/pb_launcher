import { useQuery } from "@tanstack/react-query";
import { configService } from "../services/config";

export const useProxyConfigs = () => {
  const { data } = useQuery({
    queryKey: ["proxy_configs"],
    queryFn: ({ signal }) => configService.fetchProxyConfigs(signal),
  });
  return data ?? {};
};
