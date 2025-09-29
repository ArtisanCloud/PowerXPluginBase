// ~/lib/bridge/powerx-bridge-client.ts
export type PowerXIncomingMessage =
  | { source: 'powerx'; type: 'sync';   locale: string; theme: string; hostOrigin?: string; pluginId?: string; instanceId?: string }
  | { source: 'powerx'; type: 'locale'; locale: string }
  | { source: 'powerx'; type: 'theme';  theme: string };

export type PowerXOutgoingMessage =
  | { source: 'plugin'; type: 'ready'; pluginId?: string; instanceId?: string }
  | { source: 'plugin'; type: 'request-sync' }
  | { source: 'plugin'; type: 'ping'; ts: number };

export type BridgeOptions = {
  allowedOrigins?: string[];
  debug?: boolean;
  pluginId?: string;
  instanceId?: string;
  onLocale?: (locale: string) => void;
  onTheme?: (theme: string) => void;
  onSync?: (p: { locale: string; theme: string; pluginId?: string; instanceId?: string }) => void;
};

function guessHostOrigin(): string | null {
  try { if (document.referrer) return new URL(document.referrer).origin; } catch {}
  try {
    const u = new URL(window.location.href);
    const px = u.searchParams.get('px_host');
    if (px) return new URL(px).origin;
  } catch {}
  return null;
}

export class PowerXBridgeClient {
  private allowed = new Set<string>();
  private debug: boolean;
  private onLocale?: (locale: string) => void;
  private onTheme?: (theme: string) => void;
  private onSync?: (p: { locale: string; theme: string; pluginId?: string; instanceId?: string }) => void;
  private boundHandler: (e: MessageEvent) => void;
  private pluginId?: string;
  private instanceId?: string;

  constructor(opts: BridgeOptions = {}) {
    this.debug = !!opts.debug;
    this.onLocale = opts.onLocale;
    this.onTheme = opts.onTheme;
    this.onSync = opts.onSync;
    this.pluginId = opts.pluginId;
    this.instanceId = opts.instanceId;

    const initial = (opts.allowedOrigins && opts.allowedOrigins.length > 0)
      ? opts.allowedOrigins
      : (guessHostOrigin() ? [guessHostOrigin() as string] : []);
    initial.forEach(o => this.allowed.add(o));

    this.boundHandler = (e: MessageEvent) => this.handleMessage(e);
  }

  start() {
    window.addEventListener('message', this.boundHandler, false);
    this.post({ source: 'plugin', type: 'ready', pluginId: this.pluginId, instanceId: this.instanceId });
    this.requestSync();
    if (this.debug) console.log('[PowerXBridge] started. allowed=', Array.from(this.allowed));
  }

  stop() {
    window.removeEventListener('message', this.boundHandler, false);
    if (this.debug) console.log('[PowerXBridge] stopped');
  }

  requestSync() {
    this.post({ source: 'plugin', type: 'request-sync' });
    if (this.debug) console.log('[PowerXBridge] request-sync sent');
  }

  ping() {
    this.post({ source: 'plugin', type: 'ping', ts: Date.now() });
  }

  addAllowedOrigin(origin: string) {
    if (!origin) return;
    this.allowed.add(origin);
    if (this.debug) console.log('[PowerXBridge] allow origin:', origin);
  }

  private post(msg: PowerXOutgoingMessage) {
    try { window.parent?.postMessage(msg, '*'); } catch {}
  }

  private handleMessage(e: MessageEvent) {
    const origin = e.origin;
    const data = e.data as PowerXIncomingMessage | undefined;
    if (!data || data.source !== 'powerx') return;

    if (data.type === 'sync' && data.hostOrigin) {
      this.allowed = new Set([data.hostOrigin]);
      if (this.debug) console.log('[PowerXBridge] synced host origin ->', data.hostOrigin);
    }

    if (this.allowed.size && !this.allowed.has(origin)) {
      if (this.debug) console.warn('[PowerXBridge] blocked message from', origin, data);
      return;
    }

    if (this.debug) console.log('[PowerXBridge] <=', origin, data);

    switch (data.type) {
      case 'sync':
        this.onSync?.({ locale: data.locale, theme: data.theme, pluginId: data.pluginId, instanceId: data.instanceId });
        this.onLocale?.(data.locale);
        this.onTheme?.(data.theme);
        break;
      case 'locale':
        this.onLocale?.(data.locale);
        break;
      case 'theme':
        this.onTheme?.(data.theme);
        break;
    }
  }
}

export function initPowerXBridge(opts: BridgeOptions = {}) {
  const client = new PowerXBridgeClient(opts);
  // @ts-expect-error: expose for console debug
  (window as any).PowerXBridge = client;
  client.start();
  return client;
}
