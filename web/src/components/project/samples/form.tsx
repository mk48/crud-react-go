import { useForm } from "@tanstack/react-form"
import * as z from "zod"
import { Button } from "@/components/ui/button"
import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { useTranslation } from "react-i18next"
import { Loader2 } from "lucide-react"

const formSchema = z.object({
  name: z.string().trim().min(1).max(200),
  // Optional on the server (nullable TEXT) - an empty string here is sent
  // as null by form-new/form-update.
  description: z.string(),
})
export type SamplesFormSchema = z.infer<typeof formSchema>

// ==================================== Props ==================================
type props = {
  defaultValues: SamplesFormSchema
  submitButtonText: string
  isBusy: boolean
  onSubmit: (data: SamplesFormSchema) => void
}

const SamplesForm: React.FC<props> = ({
  defaultValues,
  onSubmit,
  submitButtonText,
  isBusy,
}) => {
  const { t } = useTranslation()

  const form = useForm({
    defaultValues,
    validators: {
      onChange: formSchema,
    },
    onSubmit: ({ value }) => onSubmit(value),
  })

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        e.stopPropagation()
        form.handleSubmit()
      }}
      className="mx-auto min-w-96 space-y-8 py-10"
    >
      {/* ----------------- Name --------------------------- */}

      <form.Field name="name">
        {(field) => (
          <Field data-invalid={!field.state.meta.isValid}>
            <FieldLabel htmlFor={field.name}>{t("samples.name")}</FieldLabel>

            <Input
              id={field.name}
              name={field.name}
              placeholder={t("samples.name-placeholder")}
              type="text"
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
            />

            <FieldDescription>{t("samples.name-description")}</FieldDescription>
            {!field.state.meta.isValid && (
              <FieldError errors={field.state.meta.errors} />
            )}
          </Field>
        )}
      </form.Field>

      {/* ----------------- Description --------------------------- */}

      <form.Field name="description">
        {(field) => (
          <Field data-invalid={!field.state.meta.isValid}>
            <FieldLabel htmlFor={field.name}>
              {t("samples.description")}
            </FieldLabel>

            <Textarea
              id={field.name}
              name={field.name}
              placeholder={t("samples.description-placeholder")}
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
            />

            <FieldDescription>
              {t("samples.description-description")}
            </FieldDescription>
            {!field.state.meta.isValid && (
              <FieldError errors={field.state.meta.errors} />
            )}
          </Field>
        )}
      </form.Field>

      <Button type="submit" disabled={isBusy}>
        {isBusy && <Loader2 className="animate-spin" />}
        {submitButtonText}
      </Button>
    </form>
  )
}

export default SamplesForm
