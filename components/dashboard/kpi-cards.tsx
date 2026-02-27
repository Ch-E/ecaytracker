"use client"

import { Card, CardContent } from "@/components/ui/card"
import { kpiStats } from "@/lib/mock-data"
import type { DashboardStats } from "@/lib/api"
import {
  Car,
  DollarSign,
  Gauge,
  TrendingUp,
  TrendingDown,
  Sparkles,
  BarChart3,
} from "lucide-react"

interface KpiCardProps {
  title: string
  value: string
  change?: number
  icon: React.ReactNode
}

function KpiCard({ title, value, change, icon }: KpiCardProps) {
  const isPositive = (change ?? 0) >= 0
  return (
    <Card className="border-border/50 relative overflow-hidden">
      <CardContent className="p-5">
        <div className="flex items-center justify-between">
          <div className="flex flex-col gap-1">
            <span className="text-muted-foreground text-xs font-medium uppercase tracking-wider">
              {title}
            </span>
            <span className="text-2xl font-bold tracking-tight text-foreground">
              {value}
            </span>
            {change !== undefined && (
              <div className="flex items-center gap-1">
                {isPositive ? (
                  <TrendingUp className="size-3 text-success" />
                ) : (
                  <TrendingDown className="size-3 text-destructive" />
                )}
                <span
                  className={`text-xs font-medium ${
                    isPositive ? "text-success" : "text-destructive"
                  }`}
                >
                  {isPositive ? "+" : ""}
                  {change}%
                </span>
                <span className="text-muted-foreground text-xs">vs last month</span>
              </div>
            )}
          </div>
          <div className="flex size-11 items-center justify-center rounded-lg bg-primary/10">
            {icon}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

interface KpiCardsProps {
  stats?: DashboardStats | null
}

export function KpiCards({ stats }: KpiCardsProps) {
  const cards = [
    {
      title: "Total Listings",
      value: (stats?.total_listings ?? kpiStats.totalListings).toLocaleString(),
      change: stats ? undefined : kpiStats.newListingsChange,
      icon: <Car className="size-5 text-primary" />,
    },
    {
      title: "Average Price",
      value: `$${Math.round(stats?.avg_price ?? kpiStats.avgPrice).toLocaleString()}`,
      change: stats ? undefined : kpiStats.avgPriceChange,
      icon: <DollarSign className="size-5 text-primary" />,
    },
    {
      title: "Median Price",
      value: `$${Math.round(stats?.median_price ?? kpiStats.medianPrice).toLocaleString()}`,
      change: stats ? undefined : kpiStats.medianPriceChange,
      icon: <BarChart3 className="size-5 text-primary" />,
    },
    {
      title: "Avg Mileage",
      value: stats?.avg_mileage
        ? `${Math.round(stats.avg_mileage).toLocaleString()} km`
        : `${kpiStats.avgMileage.toLocaleString()} mi`,
      change: stats ? undefined : kpiStats.avgMileageChange,
      icon: <Gauge className="size-5 text-primary" />,
    },
    {
      title: "New This Week",
      value: (stats?.new_this_week ?? kpiStats.newListingsThisWeek).toString(),
      change: stats ? undefined : kpiStats.newListingsChange,
      icon: <TrendingUp className="size-5 text-primary" />,
    },
    {
      title: "Great Deals",
      value: kpiStats.greatDeals.toString(),
      change: kpiStats.greatDealsChange,
      icon: <Sparkles className="size-5 text-primary" />,
    },
  ]

  return (
    <div className="grid grid-cols-2 gap-4 lg:grid-cols-3 xl:grid-cols-6">
      {cards.map((card) => (
        <KpiCard key={card.title} {...card} />
      ))}
    </div>
  )
}

