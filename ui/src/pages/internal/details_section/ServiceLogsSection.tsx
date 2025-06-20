import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState, type FC } from "react";
import { serviceService, type ServiceLog } from "../../../services/services";
import { ErrorFallback } from "../../../components/helpers/ErrorFallback";
import { useViewportHeight } from "../../../hooks/useViewportHeight";

type Props = {
  service_id: string;
};

export const ServiceLogsSection: FC<Props> = ({ service_id }) => {
  const initLogsQuery = useQuery({
    queryKey: ["services", service_id],
    queryFn: ({ signal }) =>
      serviceService.fetchServiceLogs(signal, service_id, -1),
    refetchOnMount: true,
  });

  if (initLogsQuery.isFetching) {
    return <div className="p-4">Loading...</div>;
  }

  if (initLogsQuery.isError)
    return (
      <ErrorFallback
        error={initLogsQuery.error}
        onRetry={() => setTimeout(initLogsQuery.refetch)}
      />
    );

  return (
    <LogsView initLogs={initLogsQuery.data ?? []} service_id={service_id} />
  );
};

type LogsViewProps = {
  service_id: string;
  initLogs: ServiceLog[];
};

const LogsView: FC<LogsViewProps> = ({ initLogs, service_id }) => {
  const viewHeight = useViewportHeight();
  const [logs, setLogs] = useState<ServiceLog[]>(initLogs);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const controller = new AbortController();
    let isActive = true;
    let timeoutId: NodeJS.Timeout;

    const fetchLoop = async () => {
      try {
        const newLogs = await serviceService.fetchServiceLogs(
          controller.signal,
          service_id,
          10,
        );
        if (!isActive) return;
        let isScrolledToBottom = false;
        if (containerRef.current) {
          isScrolledToBottom =
            containerRef.current.scrollTop +
              containerRef.current.clientHeight >=
            containerRef.current.scrollHeight - 5;
        }
        setLogs(prev => mergeLogsUnique(prev, newLogs));
        if (isScrolledToBottom) {
          setTimeout(() => {
            if (containerRef.current) {
              containerRef.current.scrollTop =
                containerRef.current.scrollHeight;
            }
          }, 10);
        }
      } catch (err) {
        console.log(err);
        if (!isActive) return;
      }

      if (isActive) {
        timeoutId = setTimeout(fetchLoop, 1000);
      }
    };
    timeoutId = setTimeout(fetchLoop, 2000);
    return () => {
      isActive = false;
      controller.abort();
      clearTimeout(timeoutId);
    };
  }, [service_id]);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, []);

  return (
    <div
      style={{ height: viewHeight - 270 }}
      className="text-base-content overflow-y-auto font-mono text-sm"
      ref={containerRef}
    >
      <div className="whitespace-pre-wrap space-y-1">
        {(Array.isArray(logs) ? logs : []).map(log => (
          <div
            key={log.id}
            className={log.stream === "stderr" ? "text-error" : "text-success"}
          >
            {log.message}
          </div>
        ))}
      </div>
    </div>
  );
};

const mergeLogsUnique = (
  current: ServiceLog[],
  incoming: ServiceLog[],
): ServiceLog[] => {
  const existingIds = new Set(current.map(log => log.id));
  const newLogs = incoming.filter(log => !existingIds.has(log.id));
  return [...current, ...newLogs];
};
