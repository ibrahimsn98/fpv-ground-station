import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Compass } from "lucide-react"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Stat } from "./stat"
import { GPS_MODE_NAMES, NAV_MODE_NAMES, NAV_ERROR_NAMES, NAV_ACTION_NAMES } from "@/types/flags"
import type { NavData } from "@/types/telemetry"

export function NavCard() {
  const data = useTelemetryValue<NavData | undefined>(
    useCallback((msg) => msg.nav, []),
  )

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Compass className="size-3" />
          Navigation
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <Stat label="GPS Mode" value={data ? (GPS_MODE_NAMES[data.gps_mode] ?? `${data.gps_mode}`) : undefined} />
        <Stat label="Nav Mode" value={data ? (NAV_MODE_NAMES[data.nav_mode] ?? `${data.nav_mode}`) : undefined} />
        <Stat label="Waypoint" value={data?.waypoint_num} />
        <Stat label="Action" value={data ? (NAV_ACTION_NAMES[data.nav_action] ?? `${data.nav_action}`) : undefined} />
        <Stat label="Nav Error" value={data ? (NAV_ERROR_NAMES[data.nav_error] ?? `${data.nav_error}`) : undefined} />
        <Stat label="Flags" value={data ? `0x${data.flags.toString(16).toUpperCase().padStart(2, "0")}` : undefined} />
      </CardContent>
    </Card>
  )
}
