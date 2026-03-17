import { useCallback, useEffect, useRef, useState } from "react"

const TURNSTILE_SCRIPT_ID = "cf-turnstile-script"
const TURNSTILE_SCRIPT_URL =
  "https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit"

export const turnstileSiteKey = import.meta.env.VITE_TURNSTILE_SITE_KEY

/**
 * Hook that manages the Cloudflare Turnstile widget lifecycle.
 * Handles script loading, widget rendering, and token state.
 *
 * Returns no-ops when VITE_TURNSTILE_SITE_KEY is not set.
 */
export function useTurnstile() {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const widgetIdRef = useRef<string | null>(null)
  const tokenRef = useRef("")
  const [hasToken, setHasToken] = useState(false)

  const renderWidget = useCallback(() => {
    if (!containerRef.current || !window.turnstile || widgetIdRef.current) {
      return
    }

    widgetIdRef.current = window.turnstile.render(containerRef.current, {
      sitekey: turnstileSiteKey!,
      callback: (t) => {
        tokenRef.current = t
        setHasToken(true)
      },
      "expired-callback": () => {
        tokenRef.current = ""
        setHasToken(false)
      },
      "error-callback": () => {
        tokenRef.current = ""
        setHasToken(false)
      },
      theme: "auto",
    })
  }, [])

  useEffect(() => {
    if (!turnstileSiteKey) {
      return
    }

    if (window.turnstile) {
      renderWidget()
      return
    }

    const existingScript = document.getElementById(TURNSTILE_SCRIPT_ID)
    if (existingScript) {
      existingScript.addEventListener("load", renderWidget)
      return () => {
        existingScript.removeEventListener("load", renderWidget)
      }
    }

    const script = document.createElement("script")
    script.id = TURNSTILE_SCRIPT_ID
    script.src = TURNSTILE_SCRIPT_URL
    script.async = true
    script.defer = true
    script.addEventListener("load", renderWidget)
    document.head.appendChild(script)

    return () => {
      script.removeEventListener("load", renderWidget)
    }
  }, [renderWidget])

  const reset = useCallback(() => {
    tokenRef.current = ""
    setHasToken(false)
    if (widgetIdRef.current && window.turnstile) {
      window.turnstile.reset(widgetIdRef.current)
    }
  }, [])

  /** Read the current token value. Only call from event handlers / mutation fns, not render. */
  const getToken = useCallback(() => tokenRef.current, [])

  return { containerRef, hasToken, getToken, reset }
}
