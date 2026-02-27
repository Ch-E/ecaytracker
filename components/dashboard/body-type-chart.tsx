"use client"

import { Bar, BarChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer, Cell } from "recharts"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"
import { bodyTypeData } from "@/lib/mock-data"
import type { BodyTypeStat } from "@/lib/api"

const COLORS = [
  "oklch(0.65 0.2 160)",
  "oklch(0.6 0.15 250)",
  "oklch(0.7 0.18 80)",
  "oklch(0.6 0.2 25)",
  "oklch(0.55 0.15 300)",
  "oklch(0.65 0.18 130)",
  "oklch(0.58 0.15 200)",
]

interface BodyTypeChartProps {
  bodyTypes?: BodyTypeStat[]
}

export function BodyTypeChart({ bodyTypes }: BodyTypeChartProps) {
  const chartData = bodyTypes && bodyTypes.length > 0
    ? bodyTypes.map((b) => ({ type: b.type, count: b.count, avgPrice: Math.round(b.avg_price) }))
    : bodyTypeData.map((b) => ({ type: b.type, count: b.count, avgPrice: b.avgPrice }))

  return (
    <Card className="border-border/50">
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-semibold">Body Types</CardTitle>
        <CardDescription>Distribution and average price by body type</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer
          config={{
            count: {
              label: "Listings",
              color: "oklch(0.65 0.2 160)",
            },
          }}
          className="h-[200px] w-full"
        >
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData} margin={{ top: 0, right: 10, left: 0, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="oklch(0.26 0.005 260)" />
              <XAxis
                dataKey="type"
                tick={{ fill: "oklch(0.62 0 0)", fontSize: 12 }}
                axisLine={{ stroke: "oklch(0.26 0.005 260)" }}
                tickLine={false}
              />
              <YAxis
                tick={{ fill: "oklch(0.62 0 0)", fontSize: 12 }}
                axisLine={{ stroke: "oklch(0.26 0.005 260)" }}
                tickLine={false}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    formatter={(value) => `${value} listings`}
                  />
                }
              />
              <Bar dataKey="count" radius={[4, 4, 0, 0]} barSize={36}>
                {chartData.map((_, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
        <div className="mt-4 grid grid-cols-2 gap-3 lg:grid-cols-4 xl:grid-cols-7">
          {chartData.map((item) => (
            <div key={item.type} className="flex flex-col gap-0.5">
              <span className="text-xs text-muted-foreground">{item.type} Avg</span>
              <span className="text-sm font-semibold text-foreground">
                ${item.avgPrice.toLocaleString()}
              </span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
