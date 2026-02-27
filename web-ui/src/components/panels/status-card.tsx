import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ShieldCheck } from "lucide-react"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Stat } from "./stat"
import { getFlightModeName } from "@/types/flags"
import type { StatusData } from "@/types/telemetry"

export function StatusCard() {
  const data = useTelemetryValue<StatusData | undefined>(
    useCallback((msg) => msg.status, []),
  )

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <ShieldCheck className="size-3" />
          Status
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <div className="flex flex-wrap items-center gap-1.5">
          <Badge variant={data?.armed ? "destructive" : "default"} className={data?.armed ? "animate-pulse" : ""}>
            {data?.armed ? "ARMED" : "DISARMED"}
          </Badge>
          {data?.failsafe && (
            <Badge variant="destructive" className="text-[9px] px-1 py-0 animate-pulse">
              FAILSAFE
            </Badge>
          )}
        </div>
        <Stat label="Flight Mode" value={data ? getFlightModeName(data.flight_mode) : undefined} />
      </CardContent>
    </Card>
  )
}
