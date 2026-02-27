const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"

export interface ApiListing {
  id: string
  external_id: string
  url: string
  title: string
  make: string
  model: string
  year: number | null
  mileage: number | null
  price: number
  currency: string
  condition: string
  transmission: string
  fuel_type: string
  color: string
  description: string
  images: string[] | null
  location: string
  seller_name: string
  is_active: boolean
  first_seen: string | null
  last_seen: string | null
  created_at: string | null
  updated_at: string | null
}

export interface BrandStat {
  name: string
  count: number
  avg_price: number
}

export interface BodyTypeStat {
  type: string
  count: number
  avg_price: number
}

export interface YearStat {
  year: number
  count: number
}

export interface DashboardStats {
  total_listings: number
  avg_price: number
  median_price: number
  new_this_week: number
  avg_mileage: number
  top_brands: BrandStat[]
  body_types: BodyTypeStat[]
  year_distribution: YearStat[]
}

export async function fetchStats(): Promise<DashboardStats | null> {
  try {
    const res = await fetch(`${API_BASE}/api/stats`, {
      next: { revalidate: 60 },
    })
    if (!res.ok) return null
    const json = await res.json()
    return json.data as DashboardStats
  } catch {
    return null
  }
}

export async function fetchListings(): Promise<ApiListing[]> {
  try {
    const res = await fetch(`${API_BASE}/api/listings`, {
      next: { revalidate: 60 },
    })
    if (!res.ok) return []
    const json = await res.json()
    return (json.data ?? []) as ApiListing[]
  } catch {
    return []
  }
}
