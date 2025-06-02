import type { FC } from "react";
import type { ServiceDto } from "../../../services/release";
import { Copy, MoreVertical, Pencil, Power, Trash2 } from "lucide-react";
import toast from "react-hot-toast";
import classNames from "classnames";
import { useCopyToClipboard } from "@uidotdev/usehooks";

type Props = {
  service: ServiceDto;
  onDelete: () => void;
  onEdit: () => void;
};

export const ServiceCard: FC<Props> = ({ service, onDelete, onEdit }) => {
  const [, copyToClipboard] = useCopyToClipboard();
  return (
    <div
      key={service.id}
      className="card bg-base-100 shadow border border-base-300"
    >
      <div className="card-body">
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
              className="dropdown-content menu p-2 shadow-lg bg-base-100 text-base-content rounded-box min-w-[10rem] z-[1] space-y-1"
            >
              <li>
                <button
                  className="flex items-center gap-2 w-full justify-start hover:bg-base-200 text-primary"
                  onClick={onEdit}
                >
                  <Pencil className="w-4 h-4" />
                  Edit
                </button>
              </li>
              <li>
                <details className="group">
                  <summary className="flex items-center gap-2 text-base-content cursor-pointer select-none py-1 px-2 hover:bg-base-200 rounded-md">
                    <Power className="w-4 h-4" />
                    <span className="font-medium">Power</span>
                  </summary>
                  <ul className="mt-1 space-y-1">
                    <li>
                      <button
                        disabled={
                          service.status === "running" ||
                          service.status === "pending"
                        }
                        className={classNames(
                          "flex items-center gap-2 w-full px-2 py-1 rounded-md text-left",
                          service.status === "running" ||
                            service.status === "pending"
                            ? "text-base-content/60 cursor-not-allowed"
                            : "text-success hover:bg-success/10 hover:text-success",
                        )}
                      >
                        Start
                      </button>
                    </li>
                    <li>
                      <button
                        disabled={service.status !== "running"}
                        className={classNames(
                          "flex items-center gap-2 w-full px-2 py-1 rounded-md text-left",
                          service.status !== "running"
                            ? "text-base-content/60 cursor-not-allowed"
                            : "text-warning hover:bg-warning/10 hover:text-warning",
                        )}
                      >
                        Restart
                      </button>
                    </li>
                    <li>
                      <button
                        disabled={service.status !== "running"}
                        className={classNames(
                          "flex items-center gap-2 w-full px-2 py-1 rounded-md text-left",
                          service.status !== "running"
                            ? "text-base-content/60 cursor-not-allowed"
                            : "text-error hover:bg-error/10 hover:text-error",
                        )}
                      >
                        Stop
                      </button>
                    </li>
                  </ul>
                </details>
              </li>
              <li>
                <button
                  onClick={onDelete}
                  className="flex items-center gap-2 w-full justify-start text-error hover:bg-base-200"
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
            <span>{`${service.repository} v${service.release_version}`}</span>
          </div>
          <div className="flex justify-between">
            <span className="font-medium">Status:</span>
            <span
              className={`badge badge-sm ${
                service.status === "running"
                  ? "badge-success"
                  : service.status === "pending" || service.status === "idle"
                    ? "badge-warning"
                    : service.status === "failure"
                      ? "badge-error"
                      : "badge-neutral"
              }`}
            >
              {service.status.charAt(0).toUpperCase() + service.status.slice(1)}
            </span>
          </div>

          <div className="flex justify-between">
            <span className="font-medium">Started:</span>
            <span>{new Date(service.created).toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-xs">
            <span className="font-medium">Restart Policy:</span>
            <span className="capitalize">{service.restart_policy}</span>
          </div>
        </div>
        <div className="flex">
          <a
            href={service.url}
            target="_blank"
            rel="noreferrer"
            className="link link-primary truncate text-xs flex-1"
          >
            {service.url}
          </a>
          <Copy
            className="w-4 h-4 select-none active:translate-[0.2px] cursor-pointer"
            onClick={() => {
              copyToClipboard(service.url ?? "");
              toast.success("URL copied to clipboard");
            }}
          />
        </div>
      </div>
    </div>
  );
};
