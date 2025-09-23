import { buildNavigationTarget } from "~/utils/powerx-bridge";

export default defineNuxtRouteMiddleware((to) => {
  const navigation = buildNavigationTarget(to.path);
  if (!navigation) {
    return;
  }

  const finalPath = navigation.finalPath || "/";
  if (finalPath === to.path) {
    return;
  }

  const router = useRouter();
  const resolved = router.resolve({ path: finalPath, query: to.query });
  if (!resolved.matched.length) {
    return;
  }

  return navigateTo(
    {
      path: finalPath,
      query: to.query,
      hash: to.hash,
    },
    { replace: true }
  );
});
