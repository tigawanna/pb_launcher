import classNames from "classnames";
import { ExternalLink, ShieldCheck, Trash2 } from "lucide-react";
import type { FC } from "react";
const statusStyles: Record<Domain["status"], string> = {
  active: "text-success",
  pending: "text-info",
  failure: "text-error",
  expired: "text-warning",
  renewal: "text-accent",
};

interface Domain {
  name: string;
  protocol: "http" | "https";
  status: "active" | "pending" | "failure" | "expired" | "renewal";
}

type Props = {
  domain: Domain;
  url?: string;
  port?: string;
  readonly?: boolean;
};
export const DomainCard: FC<Props> = ({ domain, port, readonly, url }) => {
  return (
    <div
      className={classNames(
        "rounded-xl border p-4 flex flex-col gap-2 hover:shadow-sm transition",
        "border-base-300 dark:border-base-100",
        "bg-base-100 dark:bg-base-200",
      )}
    >
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-1 truncate pl-2">
          <span className="text-sm font-medium text-base-content dark:text-neutral-content truncate">
            {domain.name}
          </span>
          <a
            href={url ?? `${domain.protocol}://${domain.name}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300"
          >
            <ExternalLink className="w-4 h-4" />
          </a>
        </div>
        <span
          className={classNames(
            "text-xs font-medium",
            statusStyles[domain.status],
          )}
        >
          {domain.status}
        </span>
      </div>

      <div className="flex justify-between items-center text-xs text-zinc-500 dark:text-zinc-400 mt-2">
        <div className="flex items-center gap-6">
          <div className="flex">
            <span className="badge badge-ghost badge-xs">
              {domain.protocol.toUpperCase()}
            </span>
            <span className="badge badge-ghost badge-xs">Port: {port}</span>
          </div>
          {!readonly && (
            <button
              className={classNames(
                "btn btn-xs gap-1 border",
                "text-zinc-700 dark:text-zinc-200",
                "bg-white dark:bg-zinc-800",
                "hover:bg-zinc-100 dark:hover:bg-zinc-700",
                "border-zinc-300 dark:border-zinc-700",
              )}
            >
              <ShieldCheck className="w-3 h-3 text-inherit" />
              Validar DNS
            </button>
          )}
        </div>
        {!readonly && (
          <button
            className={classNames(
              "btn-xs gap-1 cursor-pointer",
              "text-zinc-700 dark:text-zinc-200",
              "bg-white dark:bg-zinc-800",
              "hover:bg-zinc-100 dark:hover:bg-zinc-700",
              "border border-zinc-300 dark:border-zinc-700",
            )}
          >
            <Trash2 className="w-4 h-4" />
          </button>
        )}
      </div>
    </div>
  );
};
