import { useCopyToClipboard } from "@uidotdev/usehooks";
import { Check, Copy } from "lucide-react";
import type { FC } from "react";
import { useState } from "react";

type Props = {
  username: string;
  password: string;
  onResetCredentials?: () => void;
};

export const DefaultCredentialsCard: FC<Props> = ({ username, password }) => {
  const [, copyToClipboard] = useCopyToClipboard();
  const [copiedField, setCopiedField] = useState<
    "username" | "password" | null
  >(null);

  const handleCopy = (value: string, field: "username" | "password") => {
    copyToClipboard(value);
    setCopiedField(field);
    setTimeout(() => setCopiedField(null), 1200);
  };

  return (
    <div className="card w-[350px] max-w-sm bg-base-100 shadow-xl">
      <div className="card-body space-y-4">
        <h2 className="card-title">Default Credentials</h2>

        <p className="text-sm text-warning">
          These credentials were generated automatically. You must change them
          after accessing the platform.
        </p>

        <div className="space-y-2">
          <div className="flex items-center justify-between gap-4">
            <div>
              <span className="font-semibold">Username:</span>
              <div className="truncate">{username}</div>
            </div>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => handleCopy(username, "username")}
            >
              {copiedField === "username" ? (
                <Check size={18} />
              ) : (
                <Copy size={18} />
              )}
            </button>
          </div>

          <div className="flex items-center justify-between gap-4">
            <div>
              <span className="font-semibold">Password:</span>
              <div className="truncate">{password}</div>
            </div>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => handleCopy(password, "password")}
            >
              {copiedField === "password" ? (
                <Check size={18} />
              ) : (
                <Copy size={18} />
              )}
            </button>
          </div>
        </div>

        {/* <div className="card-actions justify-end">
          <button
            className="btn btn-outline btn-error btn-sm"
            onClick={onResetCredentials}
          >
            Reset Credentials
          </button>
        </div> */}
      </div>
    </div>
  );
};
