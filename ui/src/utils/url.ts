export const joinUrls = (...parts: string[]): string => {
  const leadingSlash = new RegExp("^/+");
  const trailingSlash = new RegExp("/+$");

  return parts
    .map((p) => p.replace(leadingSlash, "").replace(trailingSlash, ""))
    .filter(Boolean)
    .join("/");
};
