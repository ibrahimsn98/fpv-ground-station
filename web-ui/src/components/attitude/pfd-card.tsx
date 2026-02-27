import { useCallback } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Gauge } from "lucide-react"
import { AttitudeIndicator } from "./attitude-indicator"
import { HeadingCompass } from "./heading-compass"
import { useTelemetryValue } from "@/hooks/use-telemetry-value"

interface PFDOverlayData {
  speed?: number
  altitude?: number
}

export function PFDCard() {
  const overlay = useTelemetryValue<PFDOverlayData>(
    useCallback(
      (msg) => ({
        speed: msg.gps?.ground_speed,
        altitude: msg.gps?.altitude,
      }),
      [],
    ),
  )

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
          <Gauge className="size-3" />
          Primary Flight Display
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col items-center gap-2 px-3">
        {/* Attitude indicator with flanking speed/altitude boxes */}
        <div className="relative w-full max-w-[200px] aspect-square">
          <AttitudeIndicator />

          {/* Speed — left side */}
          <div
            className="absolute left-0 top-1/2 -translate-y-1/2 bg-black/60 rounded px-1 py-0.5 border border-white/10"
            style={{ fontFamily: "'JetBrains Mono Variable', monospace", fontVariantNumeric: "tabular-nums" }}
          >
            <div className="text-[8px] text-white/60 leading-none">SPD</div>
            <div className="text-[11px] text-white leading-tight">
              {overlay?.speed != null ? overlay.speed.toFixed(1) : "---"}
            </div>
            <div className="text-[7px] text-white/50 leading-none">m/s</div>
          </div>

          {/* Altitude — right side */}
          <div
            className="absolute right-0 top-1/2 -translate-y-1/2 bg-black/60 rounded px-1 py-0.5 border border-white/10"
            style={{ fontFamily: "'JetBrains Mono Variable', monospace", fontVariantNumeric: "tabular-nums" }}
          >
            <div className="text-[8px] text-white/60 leading-none">ALT</div>
            <div className="text-[11px] text-white leading-tight">
              {overlay?.altitude != null ? overlay.altitude.toFixed(1) : "---"}
            </div>
            <div className="text-[7px] text-white/50 leading-none">m</div>
          </div>
        </div>

        <div className="w-full max-w-[200px] h-12">
          <HeadingCompass />
        </div>
      </CardContent>
    </Card>
  )
}
