import { useContext, useEffect, useRef } from "react"
import { TelemetryContext } from "@/providers/telemetry-provider"

const WIDTH = 200
const HEIGHT = 48
const CX = WIDTH / 2
const DEG_PX = 1.8 // pixels per degree
const FONT = "'JetBrains Mono Variable', monospace"
const CARDINALS: Record<number, string> = {
  0: "N",
  90: "E",
  180: "S",
  270: "W",
}

export function HeadingCompass() {
  const { subscribe } = useContext(TelemetryContext)
  const tapeRef = useRef<SVGGElement>(null)
  const headingTextRef = useRef<SVGTextElement>(null)

  const targetHeading = useRef(0)
  const currentHeading = useRef(0)

  useEffect(() => {
    return subscribe((msg) => {
      if (msg.attitude) {
        targetHeading.current = msg.attitude.heading
      }
    })
  }, [subscribe])

  useEffect(() => {
    let rafId: number

    function animate() {
      // Shortest-path interpolation
      let delta = targetHeading.current - currentHeading.current
      if (delta > 180) delta -= 360
      if (delta < -180) delta += 360
      currentHeading.current += delta * 0.3

      // Normalize to 0-360
      currentHeading.current = ((currentHeading.current % 360) + 360) % 360

      const offset = CX - currentHeading.current * DEG_PX

      if (tapeRef.current) {
        tapeRef.current.setAttribute("transform", `translate(${offset}, 0)`)
      }

      if (headingTextRef.current) {
        const hdg = Math.round(currentHeading.current) % 360
        headingTextRef.current.textContent = `${String(hdg).padStart(3, "0")}\u00B0`
      }

      rafId = requestAnimationFrame(animate)
    }

    rafId = requestAnimationFrame(animate)
    return () => cancelAnimationFrame(rafId)
  }, [])

  // Build tape ticks for two full rotations (for wrap-around)
  const ticks: { x: number; deg: number; major: boolean; label?: string }[] = []
  for (let pass = 0; pass < 2; pass++) {
    for (let deg = 0; deg < 360; deg += 10) {
      const x = (pass * 360 + deg) * DEG_PX
      const label = CARDINALS[deg] ?? (deg % 30 === 0 ? `${deg}` : undefined)
      ticks.push({ x, deg, major: deg % 30 === 0, label })
    }
  }

  return (
    <svg viewBox={`0 0 ${WIDTH} ${HEIGHT}`} className="w-full h-full">
      <defs>
        <clipPath id="compass-clip">
          <rect x={0} y={16} width={WIDTH} height={HEIGHT - 16} />
        </clipPath>
      </defs>

      <rect x={0} y={0} width={WIDTH} height={HEIGHT} fill="#111" rx={4} />

      {/* Heading readout box at top */}
      <rect
        x={CX - 24}
        y={1}
        width={48}
        height={14}
        rx={2}
        fill="#111"
        stroke="#fbbf24"
        strokeWidth={1}
      />
      <text
        ref={headingTextRef}
        x={CX}
        y={12}
        fill="#fbbf24"
        fontSize={10}
        fontFamily={FONT}
        fontWeight="bold"
        textAnchor="middle"
        style={{ fontVariantNumeric: "tabular-nums" }}
      >
        000°
      </text>

      {/* Tape area — clipped below heading readout */}
      <g clipPath="url(#compass-clip)">
        <g ref={tapeRef}>
          {ticks.map(({ x, major, label }, i) => (
            <g key={i}>
              <line
                x1={x}
                y1={HEIGHT}
                x2={x}
                y2={HEIGHT - (major ? 12 : 7)}
                stroke="white"
                strokeWidth={major ? 1.2 : 0.6}
              />
              {label && (
                <text
                  x={x}
                  y={HEIGHT - 15}
                  fill={CARDINALS[ticks[i].deg] ? "#fbbf24" : "white"}
                  fontSize={CARDINALS[ticks[i].deg] ? 11 : 9}
                  fontFamily={FONT}
                  fontWeight={CARDINALS[ticks[i].deg] ? "bold" : "normal"}
                  textAnchor="middle"
                >
                  {label}
                </text>
              )}
            </g>
          ))}
        </g>
      </g>

      {/* Center indicator */}
      <polygon
        points={`${CX},${HEIGHT} ${CX - 4},${HEIGHT - 6} ${CX + 4},${HEIGHT - 6}`}
        fill="#fbbf24"
      />
    </svg>
  )
}
