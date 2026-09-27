import { SECOND } from "@/lib/constants"
import { differenceInSeconds, formatDistanceToNowStrict } from "date-fns"
import type { JSX } from "react"
import { useState } from "react"
import { useInterval } from "usehooks-ts"

interface props {
  time: Date | string
}

const calculateDelay = (time: Date | string) => {
  const diff = differenceInSeconds(time, new Date())
  if (diff <= -60) {
    return 60
  } else {
    return 1
  }
}

const DisplayTimeRelativeToNow: React.FC<props> = ({ time }): JSX.Element => {
  // value not used, just to re-render, so the relative time will be calculated automatically in date-fns lib
  const [, setCount] = useState<number>(0)

  // Recomputed on every render (each tick re-renders), so the interval
  // slows from 1s to 60s once the time is over a minute old, and resets if
  // `time` changes.
  useInterval(() => setCount((c) => c + 1), calculateDelay(time) * SECOND)

  return <>{formatDistanceToNowStrict(time, { addSuffix: true })}</>
}

export default DisplayTimeRelativeToNow
