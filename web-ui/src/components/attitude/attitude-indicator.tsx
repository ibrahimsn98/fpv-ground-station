import { useContext, useEffect, useRef } from "react"
import { TelemetryContext } from "@/providers/telemetry-provider"

const SIZE = 200
const CX = SIZE / 2
const CY = SIZE / 2
const R = 90 // instrument radius
const SKY = "#1a6fc4"
const GROUND = "#6b3a1f"
const PPD = 3 // pixels per degree of pitch
const FONT = "'JetBrains Mono Variable', monospace"

export function AttitudeIndicator() {
  const { subscribe } = useContext(TelemetryContext)
  const svgRef = useRef<SVGSVGElement>(null)
  const horizonRef = useRef<SVGGElement>(null)
  const rollPointerRef = useRef<SVGPolygonElement>(null)
  const pitchTextRef = useRef<SVGTextElement>(null)
  const rollTextRef = useRef<SVGTextElement>(null)

  // Animation state (not React state — direct DOM)
  const targetRoll = useRef(0)
  const targetPitch = useRef(0)
  const currentRoll = useRef(0)
  const currentPitch = useRef(0)

  useEffect(() => {
    return subscribe((msg) => {
      if (msg.attitude) {
        targetRoll.current = msg.attitude.roll
        targetPitch.current = msg.attitude.pitch
      }
    })
  }, [subscribe])

  useEffect(() => {
    let rafId: number

    function animate() {
      const lerp = 0.3
      currentRoll.current += (targetRoll.current - currentRoll.current) * lerp
      currentPitch.current += (targetPitch.current - currentPitch.current) * lerp

      const roll = currentRoll.current
      const pitch = currentPitch.current
      const pitchPx = pitch * PPD

      if (horizonRef.current) {
        horizonRef.current.setAttribute(
          "transform",
          `rotate(${-roll}, ${CX}, ${CY}) translate(0, ${pitchPx})`,
        )
      }

      if (rollPointerRef.current) {
        rollPointerRef.current.setAttribute(
          "transform",
          `rotate(${-roll}, ${CX}, ${CY})`,
        )
      }

      if (pitchTextRef.current) {
        pitchTextRef.current.textContent = `${pitch.toFixed(1)}\u00B0`
      }

      if (rollTextRef.current) {
        rollTextRef.current.textContent = `${roll.toFixed(1)}\u00B0`
      }

      rafId = requestAnimationFrame(animate)
    }

    rafId = requestAnimationFrame(animate)
    return () => cancelAnimationFrame(rafId)
  }, [])

  // Pitch ladder ticks
  const pitchTicks: { y: number; deg: number; major: boolean }[] = []
  for (let deg = -60; deg <= 60; deg += 5) {
    if (deg === 0) continue
    pitchTicks.push({ y: CY - deg * PPD, deg, major: deg % 10 === 0 })
  }

  // Roll arc ticks
  const rollAngles = [0, 10, 20, 30, 45, 60, 90, -10, -20, -30, -45, -60, -90]

  return (
    <svg
      ref={svgRef}
      viewBox={`0 0 ${SIZE} ${SIZE}`}
      className="w-full h-full"
    >
      <defs>
        <clipPath id="bezel-clip">
          <circle cx={CX} cy={CY} r={R} />
        </clipPath>
        <linearGradient id="bezel-gradient" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor="#666" />
          <stop offset="50%" stopColor="#444" />
          <stop offset="100%" stopColor="#222" />
        </linearGradient>
        <filter id="bezel-shadow">
          <feDropShadow dx={0} dy={1} stdDeviation={2} floodColor="black" floodOpacity={0.5} />
        </filter>
      </defs>

      {/* Background (visible if horizon moves away) */}
      <circle cx={CX} cy={CY} r={R} fill="#111" />

      {/* Horizon group — clipped to bezel */}
      <g clipPath="url(#bezel-clip)">
        <g ref={horizonRef}>
          {/* Sky */}
          <rect
            x={CX - R * 2}
            y={CY - R * 4}
            width={R * 4}
            height={R * 4}
            fill={SKY}
          />
          {/* Ground */}
          <rect
            x={CX - R * 2}
            y={CY}
            width={R * 4}
            height={R * 4}
            fill={GROUND}
          />
          {/* Horizon line */}
          <line
            x1={CX - R * 2}
            y1={CY}
            x2={CX + R * 2}
            y2={CY}
            stroke="white"
            strokeWidth={1.5}
          />

          {/* Pitch ladder — chevron-style */}
          {pitchTicks.map(({ y, deg, major }) => {
            const halfW = major ? 28 : 14
            // Chevron: outer tips droop for nose-up (deg>0), rise for nose-down (deg<0)
            const tipDy = major ? (deg > 0 ? 3 : -3) : 0
            return (
              <g key={deg}>
                {major ? (
                  <polyline
                    points={`${CX - halfW},${y + tipDy} ${CX - halfW + 8},${y} ${CX + halfW - 8},${y} ${CX + halfW},${y + tipDy}`}
                    fill="none"
                    stroke="white"
                    strokeWidth={1.2}
                  />
                ) : (
                  <line
                    x1={CX - halfW}
                    y1={y}
                    x2={CX + halfW}
                    y2={y}
                    stroke="white"
                    strokeWidth={0.8}
                  />
                )}
                {major && (
                  <>
                    <text
                      x={CX + halfW + 4}
                      y={y + tipDy + 3}
                      fill="white"
                      fontSize={8}
                      fontFamily={FONT}
                    >
                      {Math.abs(deg)}
                    </text>
                    <text
                      x={CX - halfW - 4}
                      y={y + tipDy + 3}
                      fill="white"
                      fontSize={8}
                      fontFamily={FONT}
                      textAnchor="end"
                    >
                      {Math.abs(deg)}
                    </text>
                  </>
                )}
              </g>
            )
          })}
        </g>
      </g>

      {/* Roll arc (fixed) */}
      <g>
        {rollAngles.map((deg) => {
          const major = deg % 30 === 0
          const rad = ((deg - 90) * Math.PI) / 180
          const innerR = R - (major ? 10 : 6)
          const x1 = CX + innerR * Math.cos(rad)
          const y1 = CY + innerR * Math.sin(rad)
          const x2 = CX + R * Math.cos(rad)
          const y2 = CY + R * Math.sin(rad)
          return (
            <line
              key={deg}
              x1={x1}
              y1={y1}
              x2={x2}
              y2={y2}
              stroke="white"
              strokeWidth={major ? 1.5 : 0.8}
            />
          )
        })}
      </g>

      {/* Roll pointer (moves with roll) */}
      <polygon
        ref={rollPointerRef}
        points={`${CX},${CY - R + 2} ${CX - 5},${CY - R + 10} ${CX + 5},${CY - R + 10}`}
        fill="white"
      />

      {/* Roll readout overlay */}
      <g>
        <rect
          x={CX - 20}
          y={CY - R + 14}
          width={40}
          height={14}
          rx={2}
          fill="black"
          fillOpacity={0.55}
        />
        <text
          ref={rollTextRef}
          x={CX}
          y={CY - R + 25}
          fill="white"
          fontSize={9}
          fontFamily={FONT}
          textAnchor="middle"
          style={{ fontVariantNumeric: "tabular-nums" }}
        >
          0.0°
        </text>
      </g>

      {/* Fixed aircraft symbol — dihedral wings with tail stabilizer */}
      <g>
        {/* Left wing with dihedral */}
        <polyline
          points={`${CX - 8},${CY} ${CX - 28},${CY} ${CX - 32},${CY - 4}`}
          fill="none"
          stroke="#fbbf24"
          strokeWidth={3}
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        {/* Right wing with dihedral */}
        <polyline
          points={`${CX + 8},${CY} ${CX + 28},${CY} ${CX + 32},${CY - 4}`}
          fill="none"
          stroke="#fbbf24"
          strokeWidth={3}
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        {/* Center body dot */}
        <circle cx={CX} cy={CY} r={3} fill="#fbbf24" />
        {/* Fuselage stub */}
        <line
          x1={CX}
          y1={CY + 5}
          x2={CX}
          y2={CY + 14}
          stroke="#fbbf24"
          strokeWidth={2}
          strokeLinecap="round"
        />
        {/* Horizontal tail stabilizer */}
        <line
          x1={CX - 10}
          y1={CY + 14}
          x2={CX + 10}
          y2={CY + 14}
          stroke="#fbbf24"
          strokeWidth={2}
          strokeLinecap="round"
        />
      </g>

      {/* Pitch readout overlay */}
      <g>
        <rect
          x={CX - 22}
          y={CY + 20}
          width={44}
          height={14}
          rx={2}
          fill="black"
          fillOpacity={0.55}
        />
        <text
          ref={pitchTextRef}
          x={CX}
          y={CY + 31}
          fill="white"
          fontSize={9}
          fontFamily={FONT}
          textAnchor="middle"
          style={{ fontVariantNumeric: "tabular-nums" }}
        >
          0.0°
        </text>
      </g>

      {/* Bezel ring — gradient with shadow */}
      <circle
        cx={CX}
        cy={CY}
        r={R + 1}
        fill="none"
        stroke="url(#bezel-gradient)"
        strokeWidth={5}
        filter="url(#bezel-shadow)"
      />
      {/* Inner bezel edge */}
      <circle
        cx={CX}
        cy={CY}
        r={R - 1}
        fill="none"
        stroke="rgba(0,0,0,0.6)"
        strokeWidth={1}
      />
    </svg>
  )
}
