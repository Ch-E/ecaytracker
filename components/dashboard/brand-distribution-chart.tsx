"use client"

import { Bar, BarChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer, Cell } from "recharts"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"
import { brandDistribution } from "@/lib/mock-data"
import type { BrandStat } from "@/lib/api"

const COLORS = [
  "oklch(0.65 0.2 160)",
  "oklch(0.6 0.15 250)",
  "oklch(0.7 0.18 80)",
  "oklch(0.6 0.2 25)",
  "oklch(0.55 0.15 300)",
  "oklch(0.65 0.2 160 / 0.7)",
  "oklch(0.6 0.15 250 / 0.7)",
  "oklch(0.7 0.18 80 / 0.7)",
]

interface BrandDistributionChartProps {
  brands?: BrandStat[]
}

export function BrandDistributionChart({ brands }: BrandDistributionChartProps) {
  const data =
    brands && brands.length > 0
      ? brands.map((b) => ({ name: b.name, count: b.count, avgPrice: b.avg_price }))
      : brandDistribution

  return (
    <Card className="border-border/50">
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-semibold">Top Brands</CardTitle>
        <CardDescription>Listings by manufacturer</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer
          config={{
            count: {
              label: "Listings",
              color: "oklch(0.65 0.2 160)",
            },
          }}
          className="h-[300px] w-full"
        >
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={data}
              layout="vertical"
              margin={{ top: 0, right: 10, left: 0, bottom: 0 }}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="oklch(0.26 0.005 260)" horizontal={false} />
              <XAxis
                type="number"
                tick={{ fill: "oklch(0.62 0 0)", fontSize: 12 }}
                axisLine={{ stroke: "oklch(0.26 0.005 260)" }}
                tickLine={false}
              />
              <YAxis
                type="category"
                dataKey="name"
                tick={{ fill: "oklch(0.62 0 0)", fontSize: 12 }}
                axisLine={{ stroke: "oklch(0.26 0.005 260)" }}
                tickLine={false}
                width={70}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    formatter={(value, name) => {
                      if (name === "count") return `${value} listings`
                      return String(value)
                    }}
                  />
                }
              />
              <Bar dataKey="count" radius={[0, 4, 4, 0]} barSize={20}>
                {data.map((_, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}

