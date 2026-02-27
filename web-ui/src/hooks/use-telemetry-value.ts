import { useContext, useEffect, useState } from "react"
import { TelemetryContext } from "@/providers/telemetry-provider"
import type { TelemetryMessage } from "@/types/telemetry"

/**
 * Reactive selector for slower-updating panels.
 * Re-renders the component when the selected value changes.
 */
export function useTelemetryValue<T>(
  selector: (msg: TelemetryMessage) => T,
): T | undefined {
  const { subscribe } = useContext(TelemetryContext)
  const [value, setValue] = useState<T | undefined>(undefined)

  useEffect(() => {
    return subscribe((msg) => {
      setValue(selector(msg))
    })
  }, [subscribe, selector])

  return value
}
