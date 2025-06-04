import { useRef, useState, type FC } from "react";
import type { ServiceDto } from "../../../services/release";
import {
  Check,
  Copy,
  MoreVertical,
  Pencil,
  Power,
  ShieldAlert,
  Trash2,
} from "lucide-react";
import classNames from "classnames";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { useModal } from "../../../components/modal/hook";
import { DefaultCredentialsCard } from "./DefaultCredentialsCard";

type Props = {
  service: ServiceDto;
  onDelete: () => void;
  onEdit: () => void;
  onStart: () => void;
  onStop: () => void;
  onRestart: () => void;
};

export const ServiceCard: FC<Props> = ({
  service,
  onDelete,
  onEdit,
  onRestart,
  onStart,
  onStop,
}) => {
  const dropdownRef = useRef<HTMLDivElement>(null);

  const { openModal } = useModal();
  const [, copyToClipboard] = useCopyToClipboard();
  const [copiedField, setCopiedField] = useState<"url" | null>(null);
  const handleCopy = (value: string, field: "url") => {
    copyToClipboard(value);
    setCopiedField(field);
    setTimeout(() => setCopiedField(null), 1200);
  };

  const executeAfterBlur = (fn: () => void) => {
    (document.activeElement as HTMLElement)?.blur?.();
    fn();
  };

  const showDefaultCredentials = () => {
    openModal(
      <DefaultCredentialsCard
        username={service.boot_user_email}
        password={service.boot_user_password}
      />,
    );
  };

  return (
    <div
      key={service.id}
      className="card bg-base-100 shadow border border-base-300"
    >
      <div className="card-body">
        <div className="flex justify-between items-start">
          <h2 className="card-title text-base-content">{service.name}</h2>
          <div ref={dropdownRef} className="dropdown dropdown-end select-none">
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
                  onClick={() => executeAfterBlur(onEdit)}
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
                        onClick={() => executeAfterBlur(onStart)}
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
                        onClick={() => executeAfterBlur(onRestart)}
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
                        onClick={() => executeAfterBlur(onStop)}
                      >
                        Stop
                      </button>
                    </li>
                  </ul>
                </details>
              </li>
              <li>
                <button
                  onClick={() => executeAfterBlur(onDelete)}
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
        <div className="flex gap-8">
          <a
            href={service.url}
            target="_blank"
            rel="noreferrer"
            className="link link-primary truncate text-xs flex-1"
          >
            {service.url}
          </a>
          <div className="flex gap-4">
            {copiedField === "url" ? (
              <Check className="w-4 h-4 select-none active:translate-[0.5px] cursor-pointer" />
            ) : (
              <Copy
                className="w-4 h-4 select-none active:translate-[0.5px] cursor-pointer"
                onClick={() => handleCopy(service.url ?? "", "url")}
              />
            )}
            <ShieldAlert
              className="w-4 h-4 active:translate-[0.5px]"
              onClick={showDefaultCredentials}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
