import { useCallback, useContext } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { TelemetryContext } from "@/providers/telemetry-provider"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Radio } from "lucide-react"
import { Stat } from "./stat"
import type { StatsPayload } from "@/types/telemetry"

export function ConnectionCard() {
  const { status, dataStatus } = useContext(TelemetryContext)

  const stats = useTelemetryValue<StatsPayload | undefined>(
    useCallback((msg) => msg.stats, []),
  )

  const statusColor =
    status === "connected"
      ? "default"
      : status === "connecting"
        ? "secondary"
        : "destructive"

  function formatUptime(sec: number): string {
    const h = Math.floor(sec / 3600)
    const m = Math.floor((sec % 3600) / 60)
    const s = Math.floor(sec % 60)
    if (h > 0) return `${h}h${m}m`
    if (m > 0) return `${m}m${s}s`
    return `${s}s`
  }

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Radio className="size-3" />
          Connection
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 px-3">
        <div className="flex items-center gap-2 flex-wrap">
          <Badge variant={statusColor as "default" | "secondary" | "destructive"}>
            {status.toUpperCase()}
          </Badge>
          <Badge
            variant={
              dataStatus === "receiving"
                ? "default"
                : dataStatus === "stale"
                  ? "secondary"
                  : "destructive"
            }
          >
            {dataStatus === "receiving"
              ? "FC RECEIVING"
              : dataStatus === "stale"
                ? "FC STALE"
                : "NO FC DATA"}
          </Badge>
        </div>

        <Stat label="Uptime" value={stats ? formatUptime(stats.uptime_sec) : undefined} />
        <Stat label="FPS" value={stats?.fps.toFixed(1)} />
        <Stat label="Total" value={stats?.total} />
        <Stat label="CRC Err" value={stats?.crc_errors} />
        <Stat label="Decode Err" value={stats?.decode_errors} />
      </CardContent>
    </Card>
  )
}
