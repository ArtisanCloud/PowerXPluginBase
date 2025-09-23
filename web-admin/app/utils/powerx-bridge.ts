export const PLUGIN_ID = "com.powerx.plugins.base";
export const PLUGIN_ADMIN_PREFIX = `/_p/${PLUGIN_ID}/admin`;
export const PLUGIN_ROUTE_BASE = `${PLUGIN_ADMIN_PREFIX}/plugins/base`;

const SUPPORTED_LOCALES = ["zh", "en"] as const;
type SupportedLocale = (typeof SUPPORTED_LOCALES)[number];

const stripEmptySegments = (segments: string[]) =>
  segments.filter((segment) => segment.trim().length > 0);

export const ensureLeadingSlash = (value: string) =>
  value.startsWith("/") ? value : `/${value}`;

export const normalizePath = (value: string): string => {
  if (!value) {
    return "/";
  }

  const withLeadingSlash = ensureLeadingSlash(value.trim());
  const segments = stripEmptySegments(withLeadingSlash.split("/"));

  if (!segments.length) {
    return "/";
  }

  return `/${segments.join("/")}`;
};

export const joinLocaleWithPath = (
  localePrefix: string,
  path: string
): string => {
  const normalized = normalizePath(path || "/");
  if (!localePrefix) {
    return normalized;
  }

  if (normalized === "/") {
    return localePrefix || "/";
  }

  return `${localePrefix}${normalized}`;
};

export interface LocaleStripResult {
  localePrefix: string;
  pathWithoutLocale: string;
}

export const stripLocalePrefix = (path: string): LocaleStripResult => {
  const normalizedPath = normalizePath(path || "/");

  const segments = stripEmptySegments(normalizedPath.split("/"));
  const potentialLocale = segments[0];

  if (
    potentialLocale &&
    SUPPORTED_LOCALES.includes(potentialLocale as SupportedLocale)
  ) {
    const remainingSegments = segments.slice(1);
    const remainingPath = remainingSegments.length
      ? `/${remainingSegments.join("/")}`
      : "/";

    return {
      localePrefix: `/${potentialLocale}`,
      pathWithoutLocale: remainingPath,
    };
  }

  return {
    localePrefix: "",
    pathWithoutLocale: normalizedPath,
  };
};

export interface ExtractedPluginRoute {
  localePrefix: string;
  internalPath: string;
  rawPath: string;
}

export const extractInternalRoute = (
  fullPath: string
): ExtractedPluginRoute | null => {
  const { localePrefix, pathWithoutLocale } = stripLocalePrefix(fullPath);
  const normalizedWithoutLocale = normalizePath(pathWithoutLocale);

  if (normalizedWithoutLocale.startsWith(PLUGIN_ROUTE_BASE)) {
    const internal = normalizedWithoutLocale.slice(PLUGIN_ROUTE_BASE.length);
    const internalPath = normalizePath(internal || "/");
    return {
      localePrefix,
      internalPath,
      rawPath: normalizedWithoutLocale,
    };
  }

  if (normalizedWithoutLocale.startsWith(PLUGIN_ADMIN_PREFIX)) {
    const internal = normalizedWithoutLocale.slice(PLUGIN_ADMIN_PREFIX.length);
    const internalPath = normalizePath(internal || "/");
    return {
      localePrefix,
      internalPath,
      rawPath: normalizedWithoutLocale,
    };
  }

  return null;
};

export const isPluginAdminPath = (fullPath: string) =>
  extractInternalRoute(fullPath) !== null;

export const resolveCanonicalInternalPath = (
  internalPath: string
): string => {
  const normalized = normalizePath(internalPath);
  const aliasTarget = resolveInternalAlias(normalized);
  if (!aliasTarget) {
    return normalized;
  }

  return normalizePath(aliasTarget);
};

const stripLegacyPluginPrefix = (path: string): string => {
  if (!path.startsWith("/plugins/base")) {
    return path;
  }

  const stripped = path.slice("/plugins/base".length);
  return normalizePath(stripped || "/");
};

export const resolveInternalAlias = (internalPath: string): string | null => {
  const normalized = normalizePath(internalPath);

  const sanitizedPath = stripLegacyPluginPrefix(normalized);

  if (sanitizedPath !== normalized) {
    return resolveInternalAlias(sanitizedPath);
  }

  if (sanitizedPath === "/" || sanitizedPath === "/dashboard") {
    return "/dashboard";
  }

  if (sanitizedPath.startsWith("/notes")) {
    return sanitizedPath;
  }

  if (sanitizedPath.startsWith("/reports")) {
    return sanitizedPath;
  }

  if (sanitizedPath === "/settings") {
    return "/settings/project";
  }

  if (sanitizedPath.startsWith("/settings/")) {
    const subPath = sanitizedPath.slice("/settings/".length);
    if (subPath === "preferences" || subPath === "project") {
      return "/settings/project";
    }
    return `/settings/${subPath}`;
  }

  if (sanitizedPath === "/team") {
    return "/team/management";
  }

  if (sanitizedPath.startsWith("/team/")) {
    return sanitizedPath;
  }

  return null;
};

export interface NavigationTarget {
  localePrefix: string;
  internalPath: string;
  targetPath: string;
  finalPath: string;
}

export const buildNavigationTarget = (
  fullPath: string
): NavigationTarget | null => {
  const extracted = extractInternalRoute(fullPath);
  if (!extracted) {
    return null;
  }

  const normalizedFullPath = normalizePath(fullPath);

  // 如果已经位于插件嵌入路径下，则无需再执行重定向
  if (extracted.rawPath.startsWith(PLUGIN_ROUTE_BASE)) {
    return null;
  }

  const aliasTarget = resolveInternalAlias(extracted.internalPath);
  if (!aliasTarget) {
    return null;
  }

  const normalizedTarget = normalizePath(aliasTarget);
  const finalPath = joinLocaleWithPath(
    extracted.localePrefix,
    normalizedTarget
  );

  const normalizedFinalPath = finalPath || "/";

  if (normalizedFinalPath === normalizedFullPath) {
    return null;
  }

  return {
    localePrefix: extracted.localePrefix,
    internalPath: extracted.internalPath,
    targetPath: normalizedTarget,
    finalPath: normalizedFinalPath,
  };
};
