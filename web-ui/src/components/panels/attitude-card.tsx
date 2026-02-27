import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Compass } from "lucide-react"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Stat } from "./stat"
import type { AttitudeData } from "@/types/telemetry"

export function AttitudeCard() {
  const data = useTelemetryValue<AttitudeData | undefined>(
    useCallback((msg) => msg.attitude, []),
  )

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Compass className="size-3" />
          Attitude
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <Stat label="Pitch" value={data?.pitch.toFixed(1)} unit="deg" />
        <Stat label="Roll" value={data?.roll.toFixed(1)} unit="deg" />
        <Stat label="Heading" value={data?.heading.toFixed(1)} unit="deg" />
      </CardContent>
    </Card>
  )
}
