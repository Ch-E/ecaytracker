import { DashboardHeader } from "@/components/dashboard/dashboard-header"
import { KpiCards } from "@/components/dashboard/kpi-cards"
import { PriceTrendChart } from "@/components/dashboard/price-trend-chart"
import { BrandDistributionChart } from "@/components/dashboard/brand-distribution-chart"
import { MileagePriceChart } from "@/components/dashboard/mileage-price-chart"
import { ListingVolumeChart } from "@/components/dashboard/listing-volume-chart"
import { BodyTypeChart } from "@/components/dashboard/body-type-chart"
import { YearDistributionChart } from "@/components/dashboard/year-distribution-chart"
import { FairPriceTool } from "@/components/dashboard/fair-price-tool"
import { ListingsTable } from "@/components/dashboard/listings-table"
import { fetchStats, fetchListings } from "@/lib/api"

export default async function DashboardPage() {
  const [stats, listings] = await Promise.all([fetchStats(), fetchListings()])

  return (
    <div className="min-h-screen bg-background">
      <main className="mx-auto max-w-[1440px] px-4 py-6 sm:px-6 lg:px-8">
        <div className="flex flex-col gap-6">
          {/* Header */}
          <DashboardHeader />

          {/* KPI Cards */}
          <KpiCards stats={stats} />

          {/* Fair Price Estimator */}
          <FairPriceTool />

          {/* Charts Row 1: Price Trends + Brand Distribution */}
          <div className="grid gap-6 lg:grid-cols-2">
            <PriceTrendChart />
            <BrandDistributionChart brands={stats?.top_brands} />
          </div>

          {/* Charts Row 2: Mileage vs Price + Listing Volume */}
          <div className="grid gap-6 lg:grid-cols-2">
            <MileagePriceChart />
            <ListingVolumeChart />
          </div>

          {/* Charts Row 3: Body Types + Year Distribution */}
          <div className="grid gap-6 lg:grid-cols-2">
            <BodyTypeChart />
            <YearDistributionChart />
          </div>

          {/* Listings Table */}
          <ListingsTable initialListings={listings} />

          {/* Footer */}
          <footer className="flex items-center justify-between border-t border-border/50 pt-6 pb-4">
            <p className="text-xs text-muted-foreground">
              Data sourced from ecaytrade.com. Prices and availability subject to change.
            </p>
            <p className="text-xs text-muted-foreground">
              Last updated: Feb 25, 2026
            </p>
          </footer>
        </div>
      </main>
    </div>
  )
}

