import { useForm } from "@tanstack/react-form"
import * as z from "zod"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { useTranslation } from "react-i18next"
import { Loader2 } from "lucide-react"

const formSchema = z.object({
  name: z.string().trim(),
  isAdmin: z.boolean(),
})
export type UserFormSchema = z.infer<typeof formSchema>

// ==================================== Props ==================================
type props = {
  defaultValues: UserFormSchema
  submitButtonText: string
  isBusy: boolean
  // Editing yourself: the API refuses removing your own admin access.
  isAdminLocked?: boolean
  onSubmit: (data: UserFormSchema) => void
}

// Users have no Create form - they're provisioned automatically by the auth
// middleware on first sign-in. Only name and isAdmin can be edited here;
// sub/email are shown read-only in the view screen.
const UserForm: React.FC<props> = ({
  defaultValues,
  onSubmit,
  submitButtonText,
  isBusy,
  isAdminLocked = false,
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
            <FieldLabel htmlFor={field.name}>{t("user.name")}</FieldLabel>

            <Input
              id={field.name}
              name={field.name}
              placeholder={t("user.name-placeholder")}
              type="text"
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
            />

            <FieldDescription>{t("user.name-description")}</FieldDescription>
            {!field.state.meta.isValid && (
              <FieldError errors={field.state.meta.errors} />
            )}
          </Field>
        )}
      </form.Field>

      {/* ----------------- Is Admin --------------------------- */}

      <form.Field name="isAdmin">
        {(field) => (
          <Field
            data-invalid={!field.state.meta.isValid}
            orientation="horizontal"
          >
            <Checkbox
              id={field.name}
              checked={field.state.value}
              disabled={isAdminLocked}
              onCheckedChange={(checked) => field.handleChange(!!checked)}
            />
            <FieldLabel htmlFor={field.name}>{t("user.is-admin")}</FieldLabel>
            {isAdminLocked && (
              <FieldDescription>
                {t("user.is-admin-self-locked")}
              </FieldDescription>
            )}
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

export default UserForm
