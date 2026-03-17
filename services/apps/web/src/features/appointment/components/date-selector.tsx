import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react"
import {
  addLocalDays,
  formatLongDate,
  formatShortDate,
  startOfLocalDay,
} from "../helpers/datetime"

interface DateSelectorProps {
  date: Date
  onChange: (date: Date) => void
}

const DateSelector = ({ date, onChange }: DateSelectorProps) => {
  const today = startOfLocalDay(new Date())
  const currentDay = startOfLocalDay(date)
  const previousDay = addLocalDays(currentDay, -1)
  const nextDay = addLocalDays(currentDay, 1)
  const canGoPrevious = previousDay >= today

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-xl">Select a day</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4 pt-6">
        <div className="rounded-xl border bg-muted p-5">
          <p className="text-[11px] tracking-[0.24em] text-muted-foreground uppercase">
            Current date
          </p>
          <p className="mt-3 text-2xl font-semibold tracking-tight">
            {formatLongDate(currentDay)}
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-1 xl:grid-cols-2">
          {canGoPrevious ? (
            <Button
              className="h-auto justify-between px-4 py-4 text-left"
              onClick={() => onChange(previousDay)}
            >
              <span className="flex items-center gap-3">
                <ChevronLeftIcon className="size-4" />
                <span>
                  <span className="block text-[11px] tracking-[0.2em] text-muted-foreground uppercase">
                    Previous day
                  </span>
                  <span className="block text-sm font-medium">
                    {formatShortDate(previousDay)}
                  </span>
                </span>
              </span>
            </Button>
          ) : null}

          <Button
            className="h-auto justify-between px-4 py-4 text-left"
            onClick={() => onChange(nextDay)}
          >
            <span>
              <span className="block text-[11px] tracking-[0.2em] text-muted-foreground uppercase">
                Next day
              </span>
              <span className="block text-sm font-medium">
                {formatShortDate(nextDay)}
              </span>
            </span>
            <ChevronRightIcon className="size-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

export default DateSelector
