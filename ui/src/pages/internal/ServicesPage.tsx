import { useEffect, useState } from "react";
import {
  MoreVertical,
  Plus,
  Trash2,
  RefreshCcw,
  Pencil,
  Power,
} from "lucide-react";

type Service = {
  id: string;
  name: string;
  version: string;
  running: boolean;
  restartPolicy: string;
  url: string;
  startedAt: string;
};

const STATUS_FILTER_KEY = "pb-dashboard-status-filter";
type TStatus = "all" | "running" | "stopped";
export const ServicesPage = () => {
  const [services, setServices] = useState<Service[]>([]);
  const [query, setQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<TStatus>(
    () => (localStorage.getItem(STATUS_FILTER_KEY) as TStatus) || "all",
  );

  const refreshServices = () => {
    // Aquí simularías una recarga desde la API
    setServices([
      {
        id: "1",
        name: "Auth Service",
        version: "0.18.5",
        running: true,
        restartPolicy: "always",
        url: "http://localhost:8090",
        startedAt: "2025-05-29T18:22:12.000Z",
      },
      {
        id: "2",
        name: "Email Queue",
        version: "0.18.5",
        running: false,
        restartPolicy: "on-failure",
        url: "http://localhost:8091",
        startedAt: "2025-05-29T15:00:00.000Z",
      },
      {
        id: "3",
        name: "File Storage",
        version: "0.18.5",
        running: true,
        restartPolicy: "always",
        url: "http://localhost:8092",
        startedAt: "2025-05-28T10:12:00.000Z",
      },
      {
        id: "4",
        name: "Webhook Dispatcher",
        version: "0.18.5",
        running: false,
        restartPolicy: "never",
        url: "http://localhost:8093",
        startedAt: "2025-05-29T08:45:00.000Z",
      },
      {
        id: "5",
        name: "Realtime Listener",
        version: "0.18.5",
        running: true,
        restartPolicy: "always",
        url: "http://localhost:8094",
        startedAt: "2025-05-29T06:30:00.000Z",
      },
      {
        id: "6",
        name: "PDF Generator",
        version: "0.18.5",
        running: false,
        restartPolicy: "on-failure",
        url: "http://localhost:8095",
        startedAt: "2025-05-27T20:00:00.000Z",
      },
    ]);
  };

  useEffect(() => {
    refreshServices();
  }, []);

  useEffect(() => {
    localStorage.setItem(STATUS_FILTER_KEY, statusFilter);
  }, [statusFilter]);

  const filtered = services
    .filter(s => s.name.toLowerCase().includes(query.toLowerCase()))
    .filter(s => {
      if (statusFilter === "all") return true;
      return statusFilter === "running" ? s.running : !s.running;
    });

  const handleDelete = (id: string) => {
    setServices(prev => prev.filter(s => s.id !== id));
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
          <button onClick={refreshServices} className="btn btn-sm btn-ghost">
            <RefreshCcw className="w-4 h-4" />
          </button>
        </div>

        <div className="flex flex-col md:flex-row gap-4 select-none">
          <select
            className="select select-sm select-bordered w-full sm:w-60"
            value={statusFilter}
            onChange={e => setStatusFilter(e.target.value as TStatus)}
          >
            <option value="all">All</option>
            <option value="running">Running</option>
            <option value="stopped">Stopped</option>
          </select>
          <button className="btn btn-sm btn-primary gap-2 w-full sm:w-auto">
            <Plus className="w-4 h-4" />
            New instance
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {filtered.map(service => (
          <div
            key={service.id}
            className="card bg-base-100 shadow border border-base-300"
          >
            <div className="card-body space-y-3">
              <div className="flex justify-between items-start">
                <h2 className="card-title text-base-content">{service.name}</h2>
                <div className="dropdown dropdown-end">
                  <label
                    tabIndex={0}
                    className="btn btn-sm btn-ghost btn-circle text-base-content"
                  >
                    <MoreVertical className="w-4 h-4" />
                  </label>
                  <ul
                    tabIndex={0}
                    className="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-40 z-[1] border border-base-300 space-y-1"
                  >
                    <li>
                      <button className="text-primary">
                        <Pencil className="w-4 h-4" />
                        Edit
                      </button>
                    </li>
                    <li>
                      <details>
                        <summary className="text-base-content">
                          <Power className="w-4 h-4" />
                          Power
                        </summary>
                        <ul className="p-1">
                          <li>
                            <button className="text-success">Start</button>
                          </li>
                          <li>
                            <button className="text-warning">Restart</button>
                          </li>
                          <li>
                            <button className="text-error">Stop</button>
                          </li>
                        </ul>
                      </details>
                    </li>
                    <li>
                      <button
                        onClick={() => handleDelete(service.id)}
                        className="text-error"
                      >
                        <Trash2 className="w-4 h-4" />
                        Delete
                      </button>
                    </li>
                  </ul>
                </div>
              </div>

              <div className="text-sm space-y-1 text-base-content/80">
                <div className="flex justify-between">
                  <span className="font-medium">Version:</span>
                  <span>{service.version}</span>
                </div>
                <div className="flex justify-between">
                  <span className="font-medium">Status:</span>
                  <span
                    className={`badge badge-sm ${service.running ? "badge-success" : "badge-error"}`}
                  >
                    {service.running ? "Running" : "Stopped"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="font-medium">Policy:</span>
                  <span className="capitalize">{service.restartPolicy}</span>
                </div>
                <div className="flex justify-between items-center gap-2">
                  <span className="font-medium">URL:</span>
                  <a
                    href={service.url}
                    target="_blank"
                    rel="noreferrer"
                    className="link link-primary truncate"
                  >
                    {service.url}
                  </a>
                </div>
                <div className="flex justify-between">
                  <span className="font-medium">Started:</span>
                  <span>{new Date(service.startedAt).toLocaleString()}</span>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
