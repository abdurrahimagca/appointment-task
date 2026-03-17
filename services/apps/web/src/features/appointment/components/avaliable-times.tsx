import { useState } from "react"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import type { components } from "@/lib/api/schema"
import ClientForm from "./client-form"
import {
  formatLocalTime,
  formatTimeInTimezone,
  getLocalTimezone,
} from "../helpers/datetime"

type AppointmentSlot = components["schemas"]["AppointmentSlot"]

interface AvaliableTimesProps {
  errorMessage?: string
  hasNextPage?: boolean
  isFetchingNextPage: boolean
  isLoading: boolean
  onLoadMore: () => Promise<unknown>
  providerTimezone: string
  slots: AppointmentSlot[]
}

const AvaliableTimes = ({
  errorMessage,
  hasNextPage,
  isFetchingNextPage,
  isLoading,
  onLoadMore,
  providerTimezone,
  slots,
}: AvaliableTimesProps) => {
  const [selectedSlotId, setSelectedSlotId] = useState<string | null>(null)
  const selectedSlot = slots.find((slot) => slot.id === selectedSlotId) ?? null

  const hasSlots = slots.length > 0

  return (
    <Card>
      <CardHeader>
        <CardDescription>Availability for the selected date</CardDescription>
        <CardTitle className="text-xl">Available times</CardTitle>
      </CardHeader>
      <CardContent className="space-y-5 pt-6">
        <div className="flex items-center justify-between rounded-xl border bg-muted px-4 py-3 text-sm text-muted-foreground">
          <span>{getLocalTimezone()}</span>
          <span>{slots.length} slots loaded</span>
        </div>

        {isLoading ? (
          <p className="text-sm text-muted-foreground">Loading time slots...</p>
        ) : null}

        {!isLoading && errorMessage ? (
          <p className="rounded-xl border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
            {errorMessage}
          </p>
        ) : null}

        {!isLoading && !errorMessage && !hasSlots ? (
          <p className="rounded-xl border bg-muted px-4 py-6 text-sm text-muted-foreground">
            No available times were returned for this date.
          </p>
        ) : null}

        {slots.length > 0 ? (
          <div className="grid gap-3 md:grid-cols-2">
            {slots.map((slot) => {
              const isSelected = selectedSlotId === slot.id

              return (
                <Button
                  className={`h-auto justify-start p-4 whitespace-normal ${
                    isSelected ? "border-primary bg-accent" : ""
                  }`}
                  onClick={() => setSelectedSlotId(slot.id)}
                  key={slot.id}
                  type="button"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="text-base font-semibold">
                        {formatLocalTime(slot.startTimeUtc)} -{" "}
                        {formatLocalTime(slot.endTimeUtc)}
                      </p>
                      <p className="mt-1 text-sm text-muted-foreground">
                        Provider time:{" "}
                        {formatTimeInTimezone(
                          slot.startTimeUtc,
                          providerTimezone
                        )}{" "}
                        -{" "}
                        {formatTimeInTimezone(
                          slot.endTimeUtc,
                          providerTimezone
                        )}
                      </p>
                    </div>

                    {slot.acceptMultipleParticipants ? (
                      <span className="rounded-full border px-2.5 py-1 text-[11px] font-medium tracking-[0.16em] uppercase">
                        Multi-client
                      </span>
                    ) : null}
                  </div>
                </Button>
              )
            })}
          </div>
        ) : null}

        {selectedSlot ? (
          <div className="rounded-xl border bg-muted p-5">
            <div className="mb-4 space-y-1">
              <p className="text-sm font-medium">
                Book {formatLocalTime(selectedSlot.startTimeUtc)} -{" "}
                {formatLocalTime(selectedSlot.endTimeUtc)}
              </p>
              <p className="text-sm text-muted-foreground">
                Provider time:{" "}
                {formatTimeInTimezone(
                  selectedSlot.startTimeUtc,
                  providerTimezone
                )}{" "}
                -{" "}
                {formatTimeInTimezone(
                  selectedSlot.endTimeUtc,
                  providerTimezone
                )}
              </p>
            </div>
            <ClientForm
              allowMultipleParticipants={
                selectedSlot.acceptMultipleParticipants
              }
              slotId={selectedSlot.id}
            />
          </div>
        ) : null}

        {hasNextPage ? (
          <Button disabled={isFetchingNextPage} onClick={onLoadMore}>
            {isFetchingNextPage ? "Loading more times..." : "Load more times"}
          </Button>
        ) : null}
      </CardContent>
    </Card>
  )
}

export default AvaliableTimes
