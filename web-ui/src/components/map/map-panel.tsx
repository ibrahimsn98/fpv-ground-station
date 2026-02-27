import { useCallback, useContext, useEffect, useRef, useState } from "react"
import { MapContainer, TileLayer, Polyline, useMap } from "react-leaflet"
import L from "leaflet"
import "leaflet/dist/leaflet.css"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { TelemetryContext } from "@/providers/telemetry-provider"
import type { TelemetryMessage } from "@/types/telemetry"

function createUAVIcon(heading: number) {
  const svg = `<svg width="28" height="28" viewBox="0 0 28 28" xmlns="http://www.w3.org/2000/svg">
    <g transform="rotate(${heading}, 14, 14)">
      <polygon points="14,2 8,24 14,19 20,24" fill="#fbbf24" stroke="#000" stroke-width="1.5"/>
    </g>
  </svg>`
  return L.divIcon({
    html: svg,
    className: "",
    iconSize: [28, 28],
    iconAnchor: [14, 14],
  })
}

function createHomeIcon() {
  const svg = `<svg width="20" height="20" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
    <circle cx="10" cy="10" r="8" fill="#22c55e" stroke="#000" stroke-width="1.5" opacity="0.8"/>
    <text x="10" y="14" text-anchor="middle" font-size="10" fill="#000" font-weight="bold">H</text>
  </svg>`
  return L.divIcon({
    html: svg,
    className: "",
    iconSize: [20, 20],
    iconAnchor: [10, 10],
  })
}

function MapUpdater({
  track,
  setTrack,
}: {
  track: [number, number][]
  setTrack: React.Dispatch<React.SetStateAction<[number, number][]>>
}) {
  const map = useMap()
  const { subscribe } = useContext(TelemetryContext)
  const markerRef = useRef<L.Marker | null>(null)
  const homeMarkerRef = useRef<L.Marker | null>(null)
  const hasCenter = useRef(false)

  useEffect(() => {
    return subscribe((msg: TelemetryMessage) => {
      // Update UAV position from GPS frame
      if (msg.gps && msg.gps.lat !== 0) {
        const lat = msg.gps.lat
        const lon = msg.gps.lon
        const heading = msg.attitude?.heading ?? 0

        const pos: L.LatLngExpression = [lat, lon]

        if (!markerRef.current) {
          markerRef.current = L.marker(pos, {
            icon: createUAVIcon(heading),
          }).addTo(map)
        } else {
          markerRef.current.setLatLng(pos)
          markerRef.current.setIcon(createUAVIcon(heading))
        }

        if (!hasCenter.current) {
          map.setView(pos, 16)
          hasCenter.current = true
        }

        setTrack((prev) => [...prev, [lat, lon] as [number, number]])
      }

      // Update home position from Origin frame
      if (msg.origin && msg.origin.lat !== 0 && msg.origin.fix > 0) {
        const homePos: L.LatLngExpression = [msg.origin.lat, msg.origin.lon]
        if (!homeMarkerRef.current) {
          homeMarkerRef.current = L.marker(homePos, {
            icon: createHomeIcon(),
          }).addTo(map)
        } else {
          homeMarkerRef.current.setLatLng(homePos)
        }
      }
    })
  }, [subscribe, map, setTrack])

  return track.length > 1 ? (
    <Polyline positions={track} color="#fbbf24" weight={2} opacity={0.7} />
  ) : null
}

export function MapPanel() {
  const { subscribe } = useContext(TelemetryContext)
  const [sats, setSats] = useState(0)
  const [fix, setFix] = useState(0)
  const [track, setTrack] = useState<[number, number][]>([])

  useEffect(() => {
    fetch("/api/track")
      .then((r) => r.json())
      .then((points: [number, number][]) => {
        if (points?.length) setTrack(points)
      })
      .catch(() => {})
  }, [])

  useEffect(() => {
    return subscribe((msg) => {
      if (msg.gps) {
        setSats(msg.gps.sats)
        setFix(msg.gps.fix)
      }
    })
  }, [subscribe])

  const clearRoute = useCallback(() => {
    fetch("/api/track", { method: "DELETE" })
      .then(() => setTrack([]))
      .catch(() => {})
  }, [])

  const fixLabel = fix >= 2 ? (fix === 3 ? "3D" : "2D") : "No Fix"
  const fixVariant = fix >= 2 ? "default" : "destructive"

  return (
    <Card className="flex flex-col p-0 overflow-hidden">
      <CardContent className="p-0 relative">
        <MapContainer
          center={[0, 0]}
          zoom={3}
          minZoom={3}
          maxBounds={[[-85, -180], [85, 180]]}
          maxBoundsViscosity={1.0}
          className="w-full aspect-[16/9]"
          zoomControl={false}
          attributionControl={false}
        >
          <TileLayer
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            className="map-tiles"
          />
          <MapUpdater track={track} setTrack={setTrack} />
        </MapContainer>
        <div className="absolute top-2 left-2 z-[1000] flex items-center gap-2">
          <span className="text-[10px] font-medium uppercase tracking-wider text-white/80 drop-shadow-[0_1px_2px_rgba(0,0,0,0.8)]">
            Map
          </span>
          <Badge variant={fixVariant as "default" | "destructive"} className="text-[10px]">
            {fixLabel}
          </Badge>
          <span className="text-[10px] text-white/80 drop-shadow-[0_1px_2px_rgba(0,0,0,0.8)]">{sats} sats</span>
        </div>
        <button
          onClick={clearRoute}
          className="absolute top-2 right-2 z-[1000] px-2 py-0.5 text-[10px] font-medium rounded bg-black/50 text-white/80 hover:bg-black/70 hover:text-white transition-colors backdrop-blur-sm"
        >
          Clear Route
        </button>
      </CardContent>
    </Card>
  )
}
