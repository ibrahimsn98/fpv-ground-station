// LTM flight mode enum (0-21)
export const FLIGHT_MODE_NAMES: Record<number, string> = {
  0: "Manual",
  1: "Rate",
  2: "Angle",
  3: "Horizon",
  4: "Acro",
  5: "Stabilized1",
  6: "Stabilized2",
  7: "Stabilized3",
  8: "Alt Hold",
  9: "GPS Hold",
  10: "Waypoints",
  11: "Head Free",
  12: "Circle",
  13: "RTH",
  14: "Follow Me",
  15: "Land",
  16: "FBW A",
  17: "FBW B",
  18: "Cruise",
  19: "Unknown",
  20: "Launch",
  21: "Autotune",
}

export const GPS_MODE_NAMES: Record<number, string> = {
  0: "None",
  1: "PosHold",
  2: "RTH",
  3: "Mission",
}

export const NAV_MODE_NAMES: Record<number, string> = {
  0: "None",
  1: "RTH Start",
  2: "RTH Enroute",
  3: "PosHold Infinite",
  4: "PosHold Timed",
  5: "WP Enroute",
  6: "Process Next",
  7: "Jump",
  8: "Start Land",
  9: "Land In Progress",
  10: "Landed",
  11: "Settling Before Land",
  12: "Start Descent",
  13: "Hover Above Home",
  14: "Emergency Landing",
  15: "Critical GPS",
}

export const NAV_ERROR_NAMES: Record<number, string> = {
  0: "OK",
  1: "WP Too Far",
  2: "WP Sanity",
  3: "WP CRC",
  4: "Finish",
  5: "Timer Complete",
  6: "Invalid Jump",
  7: "Invalid Data",
  8: "Wait For RTH Alt",
  9: "GPS Fix Lost",
  10: "Disarmed",
  11: "Landing Check",
}

export const NAV_ACTION_NAMES: Record<number, string> = {
  0: "Unassigned",
  1: "Waypoint",
  2: "PosHold Unlim",
  3: "PosHold Timed",
  4: "RTH",
  5: "Set POI",
  6: "Jump",
  7: "Set Head",
  8: "Land",
}

export function getFlightModeName(mode: number): string {
  return FLIGHT_MODE_NAMES[mode] ?? `Mode ${mode}`
}
