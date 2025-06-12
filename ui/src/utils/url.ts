export const joinUrls = (...parts: string[]): string => {
  const leadingSlash = new RegExp("^/+");
  const trailingSlash = new RegExp("/+$");

  return parts
    .map(p => p.replace(leadingSlash, "").replace(trailingSlash, ""))
    .filter(Boolean)
    .join("/");
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
