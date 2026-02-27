// Types matching Go LTM structs (JSON tags)

export interface GPSData {
  lat: number
  lon: number
  ground_speed: number
  altitude: number
  fix: number
  sats: number
}

export interface AttitudeData {
  pitch: number
  roll: number
  heading: number
}

export interface StatusData {
  vbat: number
  mah_drawn: number
  rssi: number
  airspeed: number
  armed: boolean
  failsafe: boolean
  flight_mode: number
}

export interface OriginData {
  lat: number
  lon: number
  alt: number
  osd_on: boolean
  fix: number
}

export interface NavData {
  gps_mode: number
  nav_mode: number
  nav_action: number
  waypoint_num: number
  nav_error: number
  flags: number
}

export interface ExtraData {
  hdop: number
  hw_status: number
  x_counter: number
  disarm_reason: number
}

export interface StatsPayload {
  uptime_sec: number
  total: number
  fps: number
  crc_errors: number
  decode_errors: number
}

// Full WebSocket message envelope
export interface TelemetryMessage {
  ts: number

  gps?: GPSData
  gps_ts?: number
  attitude?: AttitudeData
  attitude_ts?: number
  status?: StatusData
  status_ts?: number
  origin?: OriginData
  origin_ts?: number
  nav?: NavData
  nav_ts?: number
  extra?: ExtraData
  extra_ts?: number

  stats?: StatsPayload
}
