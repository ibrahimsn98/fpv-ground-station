import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Battery } from "lucide-react"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Stat } from "./stat"
import type { StatusData } from "@/types/telemetry"

export function BatteryCard() {
  const data = useTelemetryValue<StatusData | undefined>(
    useCallback((msg) => msg.status, []),
  )

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Battery className="size-3" />
          Battery
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <Stat label="Voltage" value={data?.vbat.toFixed(2)} unit="V" />
        <Stat label="Consumed" value={data?.mah_drawn} unit="mAh" />
        <Stat label="RSSI" value={data?.rssi} />
        <Stat label="Airspeed" value={data?.airspeed} unit="m/s" />
      </CardContent>
    </Card>
  )
}
