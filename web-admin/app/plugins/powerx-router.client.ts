import { buildNavigationTarget } from "~/utils/powerx-bridge";

export default defineNuxtPlugin((nuxtApp) => {
  const router = nuxtApp.$router;

  router.beforeEach((to) => {
    const navigation = buildNavigationTarget(to.path);
    if (!navigation) {
      return;
    }

    const finalPath = navigation.finalPath || "/";
    if (finalPath === to.path) {
      return;
    }

    const resolved = router.resolve({ path: finalPath, query: to.query });
    if (!resolved.matched.length) {
      return;
    }

    return {
      path: finalPath,
      query: to.query,
      hash: to.hash,
    };
  });
});
