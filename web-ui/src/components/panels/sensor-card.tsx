import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Cpu } from "lucide-react"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Stat } from "./stat"
import type { ExtraData } from "@/types/telemetry"

export function SensorCard() {
  const data = useTelemetryValue<ExtraData | undefined>(
    useCallback((msg) => msg.extra, []),
  )

  const hwColor =
    data === undefined
      ? "bg-muted"
      : data.hw_status === 0
        ? "bg-green-500"
        : "bg-red-500"

  const hwLabel =
    data === undefined ? "—" : data.hw_status === 0 ? "OK" : "FAIL"

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Cpu className="size-3" />
          Hardware
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <div className="flex items-center gap-1.5">
          <div className={`w-2 h-2 rounded-full ${hwColor}`} />
          <span className="text-[10px] text-muted-foreground">
            HW Status: {hwLabel}
          </span>
        </div>
        <Stat label="Packets" value={data?.x_counter} />
        <Stat label="Disarm Reason" value={data?.disarm_reason} />
      </CardContent>
    </Card>
  )
}
