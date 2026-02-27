import { useCallback, useEffect, useRef, useState } from "react"
import type { TelemetryMessage } from "@/types/telemetry"

export type ConnectionStatus = "connecting" | "connected" | "disconnected"
export type DataStatus = "receiving" | "stale" | "none"
export type Listener = (msg: TelemetryMessage) => void

export interface TelemetryHandle {
  status: ConnectionStatus
  dataStatus: DataStatus
  messageRef: React.RefObject<TelemetryMessage | null>
  subscribe: (fn: Listener) => () => void
}

export function useTelemetry(): TelemetryHandle {
  const [status, setStatus] = useState<ConnectionStatus>("connecting")
  const [dataStatus, setDataStatus] = useState<DataStatus>("none")
  const messageRef = useRef<TelemetryMessage | null>(null)
  const listenersRef = useRef<Set<Listener>>(new Set())
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimer = useRef<ReturnType<typeof setTimeout>>(undefined)

  const subscribe = useCallback((fn: Listener) => {
    listenersRef.current.add(fn)
    return () => {
      listenersRef.current.delete(fn)
    }
  }, [])

  useEffect(() => {
    let disposed = false

    function connect() {
      if (disposed) return

      const protocol = location.protocol === "https:" ? "wss:" : "ws:"
      const ws = new WebSocket(`${protocol}//${location.host}/ws`)
      wsRef.current = ws

      setStatus("connecting")

      ws.onopen = () => {
        if (!disposed) setStatus("connected")
      }

      ws.onmessage = (ev) => {
        try {
          const msg: TelemetryMessage = JSON.parse(ev.data)
          messageRef.current = msg

          if (msg.attitude_ts != null) {
            setDataStatus(msg.ts - msg.attitude_ts < 2000 ? "receiving" : "stale")
          } else {
            setDataStatus("none")
          }

          for (const fn of listenersRef.current) {
            fn(msg)
          }
        } catch {
          // ignore parse errors
        }
      }

      ws.onclose = () => {
        if (disposed) return
        setStatus("disconnected")
        setDataStatus("none")
        reconnectTimer.current = setTimeout(connect, 2000)
      }

      ws.onerror = () => {
        ws.close()
      }
    }

    connect()

    return () => {
      disposed = true
      clearTimeout(reconnectTimer.current)
      wsRef.current?.close()
    }
  }, [])

  return { status, dataStatus, messageRef, subscribe }
}
