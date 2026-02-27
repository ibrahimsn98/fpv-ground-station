interface StatProps {
  label: string
  value: string | number | undefined
  unit?: string
}

export function Stat({ label, value, unit }: StatProps) {
  return (
    <div className="flex justify-between items-baseline gap-2">
      <span className="text-[10px] text-muted-foreground uppercase tracking-wider shrink-0">
        {label}
      </span>
      <span className="text-sm font-medium tabular-nums truncate">
        {value ?? "—"}
        {unit && value !== undefined && (
          <span className="text-[10px] text-muted-foreground ml-0.5">{unit}</span>
        )}
      </span>
    </div>
  )
}
