import { useCopyToClipboard } from "@uidotdev/usehooks";
import { Check, Copy } from "lucide-react";
import { useState, type FC } from "react";

export const CopyableField: FC<{ value: string; isUrl?: boolean }> = ({
  value,
  isUrl = false,
}) => {
  const [, copyToClipboard] = useCopyToClipboard();
  const [copied, setCopied] = useState(false);

  const handleCopy = (val: string) => {
    copyToClipboard(val);
    setCopied(true);
    setTimeout(() => setCopied(false), 1200);
  };

  return (
    <div className="flex gap-8">
      {isUrl ? (
        <a
          href={value}
          target="_blank"
          rel="noreferrer"
          className="link link-primary truncate text-xs flex-1"
        >
          {value}
        </a>
      ) : (
        <span className="truncate text-xs flex-1">{value}</span>
      )}
      <div className="flex gap-4">
        {copied ? (
          <Check className="w-4 h-4 select-none active:translate-[0.5px] cursor-pointer" />
        ) : (
          <Copy
            className="w-4 h-4 select-none active:translate-[0.5px] cursor-pointer"
            onClick={() => handleCopy(value)}
          />
        )}
      </div>
    </div>
  );
};
