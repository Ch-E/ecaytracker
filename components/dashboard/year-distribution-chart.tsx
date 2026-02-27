"use client"

import { Bar, BarChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer } from "recharts"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart"
import { yearDistribution } from "@/lib/mock-data"
import type { YearStat } from "@/lib/api"

interface YearDistributionChartProps {
  yearDist?: YearStat[]
}

export function YearDistributionChart({ yearDist }: YearDistributionChartProps) {
  const chartData = yearDist && yearDist.length > 0
    ? yearDist.map((y) => ({ year: y.year.toString(), count: y.count }))
    : yearDistribution

  return (
    <Card className="border-border/50">
      <CardHeader className="pb-2">
        <CardTitle className="text-base font-semibold">Year Distribution</CardTitle>
        <CardDescription>Number of listings by model year</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer
          config={{
            count: {
              label: "Listings",
              color: "oklch(0.6 0.15 250)",
            },
          }}
          className="h-[200px] w-full"
        >
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData} margin={{ top: 0, right: 10, left: 0, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="oklch(0.26 0.005 260)" />
              <XAxis
                dataKey="year"
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
              <Bar
                dataKey="count"
                fill="oklch(0.6 0.15 250)"
                radius={[4, 4, 0, 0]}
                barSize={28}
              />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
