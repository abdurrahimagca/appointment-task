import { useMutation } from "@tanstack/react-query"
import { useFieldArray, useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useTurnstile, turnstileSiteKey } from "@/components/turnstile"
import { apiFetchClient } from "@/lib/api/client"
import type { components } from "@/lib/api/schema"

interface ClientFormProps {
  allowMultipleParticipants: boolean
  slotId: string
}

type ErrorModel = components["schemas"]["ErrorModel"]

const clientSchema = z.object({
  email: z.email("Enter a valid email address"),
  fullName: z.string().trim().min(1, "Full name is required"),
})

const schema = z.object({
  clients: z.array(clientSchema).min(1, "At least one participant is required"),
})

type ClientFormValues = z.infer<typeof schema>

const isErrorModel = (error: unknown): error is ErrorModel => {
  if (!error || typeof error !== "object") {
    return false
  }

  return "type" in error
}

const getErrorMessage = (error: ErrorModel | null | undefined) => {
  if (!error) {
    return "Unable to create appointment."
  }

  if (error.errors && error.errors.length > 0) {
    return error.errors
      .map((detail) => detail.message)
      .filter(Boolean)
      .join(", ")
  }

  return error.detail ?? error.title ?? "Unable to create appointment."
}

const defaultClient = {
  email: "",
  fullName: "",
}

const ClientForm = ({
  allowMultipleParticipants,
  slotId,
}: ClientFormProps) => {
  const form = useForm<ClientFormValues>({
    defaultValues: {
      clients: [defaultClient],
    },
    resolver: zodResolver(schema),
  })
  const { append, fields, remove } = useFieldArray({
    control: form.control,
    name: "clients",
  })

  const { containerRef: turnstileRef, hasToken: hasTurnstileToken, getToken: getTurnstileToken, reset: resetTurnstile } = useTurnstile()

  const { data, error, isPending, mutate } = useMutation({
    mutationFn: async (values: ClientFormValues) => {
      const response = await apiFetchClient.POST("/appointments", {
        body: {
          client: values.clients[0],
          clients: values.clients,
          slotId,
          turnstileToken: getTurnstileToken() || undefined,
        },
      })

      if (response.error) {
        throw response.error
      }

      return response.data
    },
    onSuccess: () => {
      form.reset({ clients: [defaultClient] })
      resetTurnstile()
    },
  })

  const onSubmit = form.handleSubmit((values) => {
    mutate(values)
  })

  return (
    <form className="space-y-4" onSubmit={onSubmit}>
      {fields.map((field, index) => (
        <div className="space-y-4 rounded-xl border p-4" key={field.id}>
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium">
              Participant {index + 1}
            </p>
            {allowMultipleParticipants && fields.length > 1 ? (
              <Button
                className="rounded-md border px-3 py-1.5 text-sm"
                onClick={() => remove(index)}
                type="button"
              >
                Remove
              </Button>
            ) : null}
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium" htmlFor={`fullName-${index}`}>
              Full name
            </label>
            <Input
              aria-invalid={Boolean(form.formState.errors.clients?.[index]?.fullName)}
              id={`fullName-${index}`}
              placeholder="Jane Doe"
              {...form.register(`clients.${index}.fullName`)}
            />
            {form.formState.errors.clients?.[index]?.fullName ? (
              <p className="text-sm text-destructive">
                {form.formState.errors.clients[index]?.fullName?.message}
              </p>
            ) : null}
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium" htmlFor={`email-${index}`}>
              Email
            </label>
            <Input
              aria-invalid={Boolean(form.formState.errors.clients?.[index]?.email)}
              id={`email-${index}`}
              placeholder="jane@example.com"
              type="email"
              {...form.register(`clients.${index}.email`)}
            />
            {form.formState.errors.clients?.[index]?.email ? (
              <p className="text-sm text-destructive">
                {form.formState.errors.clients[index]?.email?.message}
              </p>
            ) : null}
          </div>
        </div>
      ))}

      {allowMultipleParticipants ? (
        <Button
          className="rounded-md border px-4 py-2 text-sm"
          onClick={() => append(defaultClient)}
          type="button"
        >
          Add another participant
        </Button>
      ) : null}

      {isErrorModel(error) ? (
        <div
          role="alert"
          className="rounded-xl border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive"
        >
          {getErrorMessage(error)}
        </div>
      ) : null}

      {turnstileSiteKey ? (
        <div className="min-h-16" ref={turnstileRef} />
      ) : null}

      {data ? (
        <div
          role="status"
          className="rounded-xl border border-green-500/20 bg-green-500/10 px-4 py-3 text-sm text-green-400"
        >
          Appointment created successfully.
        </div>
      ) : null}

      <Button
        className="w-full"
        disabled={isPending || Boolean(turnstileSiteKey && !hasTurnstileToken)}
        type="submit"
      >
        {isPending ? "Booking..." : "Book appointment"}
      </Button>
    </form>
  )
}

export default ClientForm
