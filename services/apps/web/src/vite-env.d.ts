/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string
  readonly VITE_TURNSTILE_SITE_KEY?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

interface Window {
  turnstile?: {
    render: (
      container: HTMLElement,
      options: {
        sitekey: string
        callback?: (token: string) => void
        "expired-callback"?: () => void
        "error-callback"?: () => void
        theme?: "light" | "dark" | "auto"
      },
    ) => string
    reset: (widgetId?: string) => void
  }
}
