import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Home } from "lucide-react"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Stat } from "./stat"
import type { OriginData } from "@/types/telemetry"

export function HomeCard() {
  const data = useTelemetryValue<OriginData | undefined>(
    useCallback((msg) => msg.origin, []),
  )

  const hasHome = data && data.fix > 0

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Home className="size-3" />
          Home
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <Badge variant={hasHome ? "default" : "secondary"} className="text-[10px]">
          {hasHome ? "HOME SET" : "NO HOME"}
        </Badge>
        <Stat label="Lat" value={data?.lat.toFixed(7)} unit={"\u00B0"} />
        <Stat label="Lon" value={data?.lon.toFixed(7)} unit={"\u00B0"} />
        <Stat label="Alt" value={data?.alt.toFixed(1)} unit="m" />
        <Stat label="OSD" value={data?.osd_on ? "On" : "Off"} />
      </CardContent>
    </Card>
  )
}
