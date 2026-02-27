import { useCallback, useContext } from "react"
import { Badge } from "@/components/ui/badge"
import { TelemetryContext } from "@/providers/telemetry-provider"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"
import { Zap, Signal, Satellite, ArrowUp, Activity } from "lucide-react"
import type { TelemetryMessage } from "@/types/telemetry"

interface TopBarData {
  armed: boolean | undefined
  voltage: number | undefined
  sats: number | undefined
  alt: number | undefined
  rssi: number | undefined
  fps: number | undefined
}

export function TopBar() {
  const { status, dataStatus } = useContext(TelemetryContext)

  const data = useTelemetryValue<TopBarData>(
    useCallback(
      (msg: TelemetryMessage) => ({
        armed: msg.status?.armed,
        voltage: msg.status?.vbat,
        sats: msg.gps?.sats,
        alt: msg.gps?.altitude,
        rssi: msg.status?.rssi,
        fps: msg.stats?.fps,
      }),
      [],
    ),
  )

  const serverDot =
    status === "connected"
      ? "bg-green-500"
      : status === "connecting"
        ? "bg-yellow-500"
        : "bg-red-500"

  const fcDot =
    dataStatus === "receiving"
      ? "bg-green-500"
      : dataStatus === "stale"
        ? "bg-yellow-500"
        : "bg-red-500"

  return (
    <header className="h-10 bg-card border-b border-border flex items-center px-4 gap-4 shrink-0">
      <div className="flex items-center gap-2">
        <div className="flex items-center gap-1">
          <div className={`w-2 h-2 rounded-full transition-colors duration-500 ${serverDot} ${status === "connected" ? "animate-pulse" : ""}`} title="Server" />
          <div className={`w-2 h-2 rounded-full transition-colors duration-500 ${fcDot} ${dataStatus === "receiving" ? "animate-pulse" : ""}`} title="FC data" />
        </div>
        <span className="text-xs font-semibold tracking-widest uppercase">
          FPV Ground Station
        </span>
      </div>

      <div className="flex-1" />

      {data?.armed !== undefined && (
        <Badge
          variant={data.armed ? "destructive" : "secondary"}
          className={`text-[10px] ${data.armed ? "animate-pulse" : ""}`}
        >
          {data.armed ? "ARMED" : "DISARMED"}
        </Badge>
      )}

      {data?.voltage !== undefined && (
        <span className="text-xs tabular-nums flex items-center gap-1">
          <Zap className="size-3 text-muted-foreground" />
          {data.voltage.toFixed(1)}V
        </span>
      )}

      {data?.rssi !== undefined && (
        <span className="text-xs tabular-nums flex items-center gap-1">
          <Signal className="size-3 text-muted-foreground" />
          {data.rssi}
        </span>
      )}

      {data?.sats !== undefined && (
        <span className="text-xs tabular-nums flex items-center gap-1">
          <Satellite className="size-3 text-muted-foreground" />
          {data.sats}
        </span>
      )}

      {data?.alt !== undefined && (
        <span className="text-xs tabular-nums flex items-center gap-1">
          <ArrowUp className="size-3 text-muted-foreground" />
          {data.alt.toFixed(1)}
          <span className="text-muted-foreground">m</span>
        </span>
      )}

      {data?.fps !== undefined && (
        <span className="text-xs tabular-nums text-muted-foreground flex items-center gap-1">
          <Activity className="size-3" />
          {data.fps.toFixed(1)}
        </span>
      )}
    </header>
  )
}
