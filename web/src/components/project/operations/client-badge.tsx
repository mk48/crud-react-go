import { cn } from "@/lib/utils"
import { clientLabel } from "./labels"

interface props {
  client: string
  className?: string
}

// The app an operation came from, as a small tag - so "did I do this on my
// phone, or did a job do it?" is answered at a glance.
const ClientBadge: React.FC<props> = ({ client, className }) => {
  const automated = client === "system" || client.startsWith("batch:")

  return (
    <span
      className={cn(
        "inline-flex items-center rounded-md border px-1.5 py-0.5 text-xs font-medium whitespace-nowrap",
        automated
          ? "border-amber-300 bg-amber-50 text-amber-800"
          : "border-gray-200 bg-gray-50 text-gray-700",
        className
      )}
    >
      {clientLabel(client)}
    </span>
  )
}

export default ClientBadge
