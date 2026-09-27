import { format } from "date-fns"
import DisplayTimeRelativeToNow from "./display-time-relative-now"

interface props {
  time: string
}

const DisplayTime: React.FC<props> = ({ time }) => {
  return (
    <>
      <span className="text-xs text-muted-foreground">
        <DisplayTimeRelativeToNow time={time} />
      </span>
      <span className="text-muted-foreground"> | </span>
      <span className="text-xs text-muted-foreground">
        {format(time, "yyyy-MM-dd HH:mm:ss")}
      </span>
    </>
  )
}

export default DisplayTime
