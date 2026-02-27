import { TopBar } from "./top-bar"
import { PFDCard } from "./attitude/pfd-card"
import { MapPanel } from "./map/map-panel"
import { BatteryCard } from "./panels/battery-card"
import { SensorCard } from "./panels/sensor-card"
import { StatusCard } from "./panels/status-card"
import { GPSCard } from "./panels/gps-card"
import { ConnectionCard } from "./panels/connection-card"
import { NavCard } from "./panels/nav-card"
import { HomeCard } from "./panels/home-card"
import { AttitudeCard } from "./panels/attitude-card"

export function Dashboard() {
  return (
    <div className="h-screen flex flex-col bg-background text-foreground overflow-hidden">
      <TopBar />

      <div className="flex-1 grid grid-cols-12 gap-2 p-2 min-h-0 overflow-hidden">
        {/* Left column: PFD + Battery + Status */}
        <div className="col-span-4 xl:col-span-3 flex flex-col gap-2 overflow-y-auto min-h-0">
          <div className="animate-fade-up" style={{ animationDelay: "0ms" }}><PFDCard /></div>
          <div className="animate-fade-up" style={{ animationDelay: "75ms" }}><AttitudeCard /></div>
          <div className="animate-fade-up" style={{ animationDelay: "150ms" }}><BatteryCard /></div>
          <div className="animate-fade-up" style={{ animationDelay: "225ms" }}><StatusCard /></div>
        </div>

        {/* Center: Map + Connection */}
        <div className="col-span-8 xl:col-span-6 flex flex-col gap-2 min-h-0">
          <div className="min-h-0 animate-fade-up" style={{ animationDelay: "50ms" }}>
            <MapPanel />
          </div>
          <div className="animate-fade-up" style={{ animationDelay: "125ms" }}>
            <ConnectionCard />
          </div>
        </div>

        {/* Right column: data panels (visible on xl+) */}
        <div className="hidden xl:flex xl:col-span-3 flex-col gap-2 overflow-y-auto min-h-0">
          <div className="animate-fade-up" style={{ animationDelay: "0ms" }}><NavCard /></div>
          <div className="animate-fade-up" style={{ animationDelay: "75ms" }}><HomeCard /></div>
          <div className="animate-fade-up" style={{ animationDelay: "150ms" }}><SensorCard /></div>
          <div className="animate-fade-up" style={{ animationDelay: "225ms" }}><GPSCard /></div>
        </div>
      </div>
    </div>
  )
}
