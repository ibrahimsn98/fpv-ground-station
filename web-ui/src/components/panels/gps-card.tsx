import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Satellite } from "lucide-react"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Stat } from "./stat"
import type { GPSData, ExtraData } from "@/types/telemetry"

interface GPSCombined {
  gps?: GPSData
  extra?: ExtraData
}

export function GPSCard() {
  const data = useTelemetryValue<GPSCombined>(
    useCallback(
      (msg) => ({
        gps: msg.gps,
        extra: msg.extra,
      }),
      [],
    ),
  )

  const gps = data?.gps
  const extra = data?.extra

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Satellite className="size-3" />
          GPS
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <Stat label="Lat" value={gps?.lat.toFixed(7)} unit={"\u00B0"} />
        <Stat label="Lon" value={gps?.lon.toFixed(7)} unit={"\u00B0"} />
        <Stat label="Alt" value={gps?.altitude.toFixed(1)} unit="m" />
        <Stat label="Speed" value={gps?.ground_speed} unit="m/s" />
        <Stat label="Sats" value={gps?.sats} />
        <Stat label="HDOP" value={extra?.hdop.toFixed(2)} />
        <Stat label="HW Status" value={extra?.hw_status === 0 ? "OK" : extra?.hw_status !== undefined ? "FAIL" : undefined} />
      </CardContent>
    </Card>
  )
}
