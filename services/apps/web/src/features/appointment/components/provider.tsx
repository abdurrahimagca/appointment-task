import { useState } from "react"
import { useInfiniteQuery, useQuery } from "@tanstack/react-query"
import type { InfiniteData } from "@tanstack/react-query"
import { CalendarDaysIcon, Clock3Icon, MapPinIcon } from "lucide-react"
import ReactMarkdown from "react-markdown"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import type { components } from "@/lib/api/schema"
import { apiFetchClient, $api } from "@/lib/api/client"
import AvaliableTimes from "./avaliable-times"
import DateSelector from "./date-selector"
import { formatDateParam, formatMediumDate } from "../helpers/datetime"

interface ProviderProps {
  username: string
}

type AvailabilityResponse = components["schemas"]["AvailabilityResponse"]
type ErrorModel = components["schemas"]["ErrorModel"]
type AvailabilityCursor = {
  cursorId: string
  cursorStartTime: string
}

const SLOT_PAGE_SIZE = 9

const getErrorMessage = (error: ErrorModel | undefined, fallback: string) => {
  if (!error) {
    return fallback
  }

  if (error.errors && error.errors.length > 0) {
    return error.errors
      .map((detail) => detail.message)
      .filter(Boolean)
      .join(", ")
  }

  return error.detail ?? fallback
}

const isErrorModel = (error: unknown): error is ErrorModel => {
  if (!error || typeof error !== "object") {
    return false
  }

  return "type" in error
}

const Provider = ({ username }: ProviderProps) => {
  const [selectedDate, setSelectedDate] = useState(() => new Date())
  const dateParam = formatDateParam(selectedDate)

  const {
    data: providerData,
    isLoading: isLoadingProvider,
    error: providerError,
  } = useQuery(
    $api.queryOptions(
      "get",
      "/providers/{username}",
      {
        params: {
          path: {
            username,
          },
        },
      },
      {
        enabled: username.length > 0,
        retry: false,
      }
    )
  )

  const {
    data: availabilityData,
    error: availabilityError,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading: isLoadingAvailability,
  } = useInfiniteQuery<
    AvailabilityResponse,
    ErrorModel,
    InfiniteData<AvailabilityResponse>,
    string[],
    AvailabilityCursor | undefined
  >({
    enabled: username.length > 0,
    getNextPageParam: (lastPage) => {
      if (!lastPage.nextCursor) {
        return undefined
      }

      return {
        cursorId: lastPage.nextCursor.id,
        cursorStartTime: lastPage.nextCursor.startTime,
      }
    },
    initialPageParam: undefined as AvailabilityCursor | undefined,
    queryFn: async ({ pageParam }) => {
      const { data, error } = await apiFetchClient.GET(
        "/providers/{username}/availability",
        {
          params: {
            path: {
              username,
            },
            query: {
              cursorId: pageParam?.cursorId,
              cursorStartTime: pageParam?.cursorStartTime,
              date: dateParam,
              limit: SLOT_PAGE_SIZE,
            },
          },
        }
      )

      if (error) {
        throw error
      }

      return data as AvailabilityResponse
    },
    queryKey: ["provider-availability", username, dateParam],
    retry: false,
  })

  if (isLoadingProvider) {
    return (
      <div className="py-12 text-sm text-muted-foreground">
        Loading provider...
      </div>
    )
  }

  if (providerError || !providerData) {
    return (
      <div className="rounded-xl border border-destructive/50 bg-destructive/10 px-5 py-4 text-sm text-destructive">
        {getErrorMessage(
          providerError ?? undefined,
          "Unable to load provider."
        )}
      </div>
    )
  }

  const provider = providerData.provider
  const timezone = availabilityData?.pages[0]?.timezone ?? provider.timezone
  const slots =
    availabilityData?.pages.flatMap((page) => page.slots ?? []) ?? []

  return (
    <div className="space-y-6">
      <section className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_360px] lg:items-start">
        <Card>
          <CardHeader>
            <CardTitle className="text-3xl tracking-tight sm:text-4xl">
              {provider.fullName}
            </CardTitle>
            <div className="flex flex-wrap gap-3 text-sm text-muted-foreground">
              <div className="inline-flex items-center gap-2 rounded-full border px-3 py-2">
                <MapPinIcon className="size-4" />
                {provider.timezone}
              </div>
              <div className="inline-flex items-center gap-2 rounded-full border px-3 py-2">
                <CalendarDaysIcon className="size-4" />
                {formatMediumDate(selectedDate)}
              </div>
              <div className="inline-flex items-center gap-2 rounded-full border px-3 py-2">
                <Clock3Icon className="size-4" />
                Provider timezone
              </div>
            </div>
          </CardHeader>
          <CardContent className="pt-6">
            <div className="text-sm leading-7 text-muted-foreground sm:text-base [&_a]:text-foreground [&_a]:underline [&_a]:underline-offset-4 [&_li]:ml-5 [&_li]:list-disc [&_p+p]:mt-4 [&_strong]:text-foreground">
              {provider.bio ? (
                <ReactMarkdown>{provider.bio}</ReactMarkdown>
              ) : null}
            </div>
          </CardContent>
        </Card>

        <DateSelector date={selectedDate} onChange={setSelectedDate} />
      </section>

      <div>
        <AvaliableTimes
          key={`${username}-${dateParam}`}
          errorMessage={
            isErrorModel(availabilityError)
              ? getErrorMessage(
                  availabilityError,
                  "Unable to load available times."
                )
              : undefined
          }
          hasNextPage={hasNextPage}
          isFetchingNextPage={isFetchingNextPage}
          isLoading={isLoadingAvailability}
          onLoadMore={fetchNextPage}
          providerTimezone={timezone}
          slots={slots}
        />
      </div>
    </div>
  )
}

export { Provider }
