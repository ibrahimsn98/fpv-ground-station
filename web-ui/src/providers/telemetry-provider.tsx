import { createContext, type ReactNode } from "react"
import { useTelemetry, type TelemetryHandle } from "@/hooks/use-telemetry"
import type { TelemetryMessage } from "@/types/telemetry"

const noopHandle: TelemetryHandle = {
  status: "disconnected",
  dataStatus: "none",
  messageRef: { current: null },
  subscribe: () => () => {},
}

export const TelemetryContext = createContext<TelemetryHandle>(noopHandle)

export function TelemetryProvider({ children }: { children: ReactNode }) {
  const handle = useTelemetry()

  return (
    <TelemetryContext.Provider value={handle}>
      {children}
    </TelemetryContext.Provider>
  )
}

// Re-export for convenience
export type { TelemetryHandle, TelemetryMessage }
