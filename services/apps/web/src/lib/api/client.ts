import createClient from "openapi-react-query"
import createFetchClient from "openapi-fetch"

import type { paths } from "@/lib/api/schema"

const baseUrl = import.meta.env.VITE_API_BASE_URL?.trim() || "/api"

export const apiFetchClient = createFetchClient<paths>({
  baseUrl,
})

export const $api = createClient(apiFetchClient)
