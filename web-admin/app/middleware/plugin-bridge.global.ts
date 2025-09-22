const pluginId = "com.powerx.plugins.base";
const pluginAdminPrefix = `/_p/${pluginId}/admin`;

// 与 i18n 配置保持一致，默认语言不带前缀，其它语言带前缀
const supportedLocales = ["zh", "en"];

const ensureLeadingSlash = (value: string) =>
  value.startsWith("/") ? value : `/${value}`;

const normalizePath = (value: string) => {
  const withSlash = ensureLeadingSlash(value);
  if (withSlash.length > 1 && withSlash.endsWith("/")) {
    return withSlash.replace(/\/+$/, "");
  }
  return withSlash;
};

// 动态规则列表，返回命中的内部目标路由；若返回 null 表示不处理
const resolveAlias = (internalPath: string): string | null => {
  // 首页 / Dashboard
  if (internalPath === "/" || internalPath === "") {
    return "/";
  }
  if (internalPath === "/plugins/base" || internalPath === "/plugins/base/dashboard") {
    return "/";
  }

  // 笔记模块：/plugins/base/notes/... -> /notes/...
  const notesMatch = internalPath.match(/^\/plugins\/base\/notes(?:\/(.*))?$/);
  if (notesMatch) {
    const suffix = notesMatch[1] ? `/${notesMatch[1]}` : "";
    return `/notes${suffix}`;
  }

  // 设置模块：/plugins/base/settings -> /settings/project，子路由保持相对路径
  if (internalPath === "/plugins/base/settings") {
    return "/settings/project";
  }
  const settingsMatch = internalPath.match(/^\/plugins\/base\/settings\/(.*)$/);
  if (settingsMatch) {
    const subPath = settingsMatch[1];
    if (subPath === "preferences" || subPath === "project") {
      return "/settings/project";
    }
    return `/settings/${subPath}`;
  }

  // 团队模块：/plugins/base/team/... -> /team/...
  const teamMatch = internalPath.match(/^\/plugins\/base\/team(?:\/(.*))?$/);
  if (teamMatch) {
    const suffix = teamMatch[1] ? `/${teamMatch[1]}` : "/management";
    return `/team${suffix}`;
  }

  return null;
};

export default defineNuxtRouteMiddleware((to) => {
  const router = useRouter();
  const path = to.path;

  // 识别是否带有多语言前缀
  let localePrefix = "";
  let remainingPath = path;
  const segments = path.split("/").filter(Boolean);
  const potentialLocale = segments[0];

  if (potentialLocale && supportedLocales.includes(potentialLocale)) {
    localePrefix = `/${potentialLocale}`;
    remainingPath = path.slice(localePrefix.length) || "/";
  }

  if (!remainingPath.startsWith(pluginAdminPrefix)) {
    return;
  }

  // 去掉 PowerX 转发前缀
  const internalPath = remainingPath.slice(pluginAdminPrefix.length) || "/";
  const normalizedInternalPath = normalizePath(internalPath);

  const aliasTarget = resolveAlias(normalizedInternalPath);
  const fallbackTarget = normalizedInternalPath;
  const targetPath = (aliasTarget ?? fallbackTarget) || "/";
  const normalizedTargetPath = normalizePath(targetPath);
  const finalPath =
    localePrefix + (normalizedTargetPath === "/" ? "" : normalizedTargetPath);

  if (finalPath === to.path) {
    return;
  }

  // 如果目标路由不存在，则沿用原始路径，避免跳入 404
  const resolved = router.resolve({ path: finalPath || "/", query: to.query });
  if (!resolved.matched.length) {
    return;
  }

  return navigateTo(
    {
      path: finalPath || "/",
      query: to.query,
      hash: to.hash,
    },
    { replace: true }
  );
});
