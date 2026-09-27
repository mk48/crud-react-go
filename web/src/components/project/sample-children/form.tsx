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
import SamplesSelect from "@/components/project/samples/select"
import { useTranslation } from "react-i18next"
import { Loader2 } from "lucide-react"

const formSchema = z.object({
  // Required foreign key - the combobox yields "" until a sample is picked.
  sampleItemId: z.string().min(1),
  // Display-only: lets SamplesSelect show the picked sample's name without
  // re-fetching it. Never sent to the server.
  sampleItemName: z.string(),
  name: z.string().trim().min(1).max(200),
})
export type SampleChildrenFormSchema = z.infer<typeof formSchema>

// ==================================== Props ==================================
type props = {
  defaultValues: SampleChildrenFormSchema
  submitButtonText: string
  isBusy: boolean
  onSubmit: (data: SampleChildrenFormSchema) => void
}

const SampleChildrenForm: React.FC<props> = ({
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
      {/* ----------------- Sample item (parent) --------------------------- */}

      <form.Field name="sampleItemId">
        {(sampleItemIdField) => (
          <form.Field name="sampleItemName">
            {(sampleItemNameField) => (
              <Field data-invalid={!sampleItemIdField.state.meta.isValid}>
                <FieldLabel>{t("sampleChildren.sample-item")}</FieldLabel>

                <SamplesSelect
                  id={sampleItemIdField.state.value}
                  name={sampleItemNameField.state.value}
                  onSelect={(id, name) => {
                    sampleItemIdField.handleChange(id)
                    sampleItemNameField.handleChange(name)
                  }}
                />

                <FieldDescription>
                  {t("sampleChildren.sample-item-description")}
                </FieldDescription>
                {!sampleItemIdField.state.meta.isValid && (
                  <FieldError errors={sampleItemIdField.state.meta.errors} />
                )}
              </Field>
            )}
          </form.Field>
        )}
      </form.Field>

      {/* ----------------- Name --------------------------- */}

      <form.Field name="name">
        {(field) => (
          <Field data-invalid={!field.state.meta.isValid}>
            <FieldLabel htmlFor={field.name}>
              {t("sampleChildren.name")}
            </FieldLabel>

            <Input
              id={field.name}
              name={field.name}
              placeholder={t("sampleChildren.name-placeholder")}
              type="text"
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
            />

            <FieldDescription>
              {t("sampleChildren.name-description")}
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

export default SampleChildrenForm
