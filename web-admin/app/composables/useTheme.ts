export const useTheme = () => {
  const colorMode = useColorMode();

  // 主题状态，直接使用 colorMode 的值
  const theme = computed(() => colorMode.value);

  // 可用主题列表 - 只有 light 和 dark
  const themes = [
    { value: "light", label: "浅色主题" },
    { value: "dark", label: "深色主题" },
  ];

  // 设置主题
  const setTheme = (newTheme: string) => {
    console.log("设置主题:", newTheme);
    colorMode.preference = newTheme;
  };

  // 初始化主题
  const initTheme = () => {
    if (import.meta.client) {
      // 优先从环境变量读取
      const runtimeConfig = useRuntimeConfig();
      const envTheme = runtimeConfig.public.defaultTheme;
      // 其次从 localStorage 读取
      const savedTheme = localStorage.getItem("theme");
      // 最后使用默认主题
      const initialTheme = envTheme || savedTheme || "light";

      // 同步当前主题状态
      theme.value = initialTheme;
      setTheme(initialTheme);
    }
  };

  // 监听外部 iframe 消息
  const listenToParentMessages = () => {
    if (import.meta.client) {
      window.addEventListener("message", (event) => {
        // 验证消息来源（可根据需要调整）
        if (event.data && event.data.type === "THEME_CHANGE") {
          setTheme(event.data.theme);
        }
        // 监听语言切换消息
        if (event.data && event.data.type === "LANGUAGE_CHANGE") {
          const { setLocale } = useI18n();
          setLocale(event.data.locale);
        }
      });
    }
  };

  return {
    theme: readonly(theme),
    themes,
    setTheme,
    initTheme,
    listenToParentMessages,
  };
};
