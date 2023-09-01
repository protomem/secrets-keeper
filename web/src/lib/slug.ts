export function slug(text: string | undefined): string | undefined {
  if (!text) return undefined;

  return text
    .toLowerCase()
    .replace(/ /g, "-")
    .replace(/[^\w-]+/g, "");
}
