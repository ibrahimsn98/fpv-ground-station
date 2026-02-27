import { TelemetryProvider } from "@/providers/telemetry-provider"
import { Dashboard } from "@/components/dashboard"

export function App() {
  return (
    <div className="dark">
      <TelemetryProvider>
        <Dashboard />
      </TelemetryProvider>
    </div>
  )
}

export default App
