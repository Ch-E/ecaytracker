"use client"

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
} from "recharts"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer } from "@/components/ui/chart"
import { mileagePriceData } from "@/lib/mock-data"
import type { ApiListing } from "@/lib/api"

interface BucketPoint {
  label: string
  avgPrice: number
  count: number
}

function buildBuckets(listings: ApiListing[]): BucketPoint[] {
  const BUCKET_SIZE = 25_000
  const MAX_MILEAGE = 300_000
  const buckets: Record<number, { sum: number; count: number }> = {}

  for (const l of listings) {
    if (l.mileage == null || l.price <= 0) continue
    const bucket = Math.min(Math.floor(l.mileage / BUCKET_SIZE) * BUCKET_SIZE, MAX_MILEAGE)
    if (!buckets[bucket]) buckets[bucket] = { sum: 0, count: 0 }
    buckets[bucket].sum += l.price
    buckets[bucket].count += 1
  }

  return Object.entries(buckets)
    .sort(([a], [b]) => Number(a) - Number(b))
    .filter(([, v]) => v.count >= 2)
    .map(([key, v]) => ({
      label: `${(Number(key) / 1000).toFixed(0)}k`,
      avgPrice: Math.round(v.sum / v.count),
      count: v.count,
    }))
}

function buildMockBuckets(): BucketPoint[] {
  const BUCKET_SIZE = 25_000
  const buckets: Record<number, { sum: number; count: number }> = {}
  for (const p of mileagePriceData) {
    const bucket = Math.floor(p.mileage / BUCKET_SIZE) * BUCKET_SIZE
    if (!buckets[bucket]) buckets[bucket] = { sum: 0, count: 0 }
    buckets[bucket].sum += p.price
    buckets[bucket].count += 1
  }
  return Object.entries(buckets)
    .sort(([a], [b]) => Number(a) - Number(b))
    .map(([key, v]) => ({
      label: `${(Number(key) / 1000).toFixed(0)}k`,
      avgPrice: Math.round(v.sum / v.count),
      count: v.count,
    }))
}

function CustomTooltip({ active, payload, label }: { active?: boolean; label?: string; payload?: Array<{ value: number; payload: BucketPoint }> }) {
  if (active && payload && payload.length) {
    const { avgPrice, count } = payload[0].payload
    return (
      <div className="rounded-lg border border-border/50 bg-background px-3 py-2 text-xs shadow-xl">
        <p className="font-semibold text-foreground">{label} km</p>
        <p className="text-primary font-medium">Avg ${avgPrice.toLocaleString()}</p>
        <p className="text-muted-foreground">{count} listings</p>
      </div>
    )
  }
  return null
}

interface MileagePriceChartProps {
  listings?: ApiListing[]
}

export function MileagePriceChart({ listings }: MileagePriceChartProps) {
  const chartData = listings && listings.length > 0
    ? buildBuckets(listings)
    : buildMockBuckets()

  return (
    <Card className="border-border/50">
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-semibold">Mileage vs Price</CardTitle>
        <CardDescription>Average price by mileage range (25k km buckets)</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer
          config={{
            avgPrice: {
              label: "Avg Price",
              color: "oklch(0.65 0.2 160)",
            },
          }}
          className="h-[300px] w-full"
        >
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chartData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="oklch(0.26 0.005 260)" />
              <XAxis
                dataKey="label"
                tick={{ fill: "oklch(0.62 0 0)", fontSize: 12 }}
                axisLine={{ stroke: "oklch(0.26 0.005 260)" }}
                tickLine={false}
              />
              <YAxis
                tick={{ fill: "oklch(0.62 0 0)", fontSize: 12 }}
                axisLine={{ stroke: "oklch(0.26 0.005 260)" }}
                tickLine={false}
                tickFormatter={(v: number) => `$${(v / 1000).toFixed(0)}k`}
              />
              <Tooltip content={<CustomTooltip />} />
              <Line
                type="monotone"
                dataKey="avgPrice"
                stroke="oklch(0.65 0.2 160)"
                strokeWidth={2}
                dot={{ r: 4, fill: "oklch(0.65 0.2 160)" }}
                activeDot={{ r: 6 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
