export const joinUrls = (...parts: string[]): string => {
  const leadingSlash = /^\/+/;
  const trailingSlash = /\/+$/;

  const cleanedParts = parts.map((part, i) => {
    const isLast = i === parts.length - 1;
    const cleaned = part.replace(leadingSlash, "");
    return isLast ? cleaned : cleaned.replace(trailingSlash, "");
  });

  return cleanedParts.filter(Boolean).join("/");
};

export const formatUrl = (
  protocol?: string,
  hostname?: string,
  port?: string,
): string => {
  if (!protocol || !hostname) return "";

  const normalizedProtocol = protocol.endsWith(":") ? protocol : `${protocol}:`;

  const isDefaultPort =
    (normalizedProtocol === "http:" && port === "80") ||
    (normalizedProtocol === "https:" && port === "443");

  const portSegment = port && !isDefaultPort ? `:${port}` : "";

  return `${normalizedProtocol}//${hostname}${portSegment}`;
};

export const extractParts = (url: URL) => {
  const { protocol, hostname, port } = url;
  const isDefaultPort =
    (protocol === "http:" && port === "80") ||
    (protocol === "https:" && port === "443");
  const portPart = port && !isDefaultPort ? `:${port}` : "";
  return { protocol, hostname, portPart };
};
