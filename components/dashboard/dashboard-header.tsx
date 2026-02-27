"use client"

import { Activity } from "lucide-react"
import { Badge } from "@/components/ui/badge"

export function DashboardHeader() {
  return (
    <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex items-center gap-3">
        <div className="flex size-10 items-center justify-center rounded-lg bg-primary/10">
          <Activity className="size-5 text-primary" />
        </div>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground sm:text-2xl text-balance">
            EcayCar Market Analytics
          </h1>
          <p className="text-sm text-muted-foreground">
            Live car market data and insights from ecaytrade.com
          </p>
        </div>
      </div>
      <div className="flex items-center gap-3">
        <Badge
          variant="outline"
          className="border-success/30 bg-success/10 text-success gap-1.5"
        >
          <span className="relative flex size-2">
            <span className="absolute inline-flex size-full animate-ping rounded-full bg-success opacity-75" />
            <span className="relative inline-flex size-2 rounded-full bg-success" />
          </span>
          Live Data
        </Badge>

      </div>
    </header>
  )
}
