"use client"

import { useState, useMemo } from "react"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { carListings } from "@/lib/mock-data"
import type { ApiListing } from "@/lib/api"
import { Search, ArrowUpDown, ExternalLink } from "lucide-react"

type SortField = "price" | "year" | "mileage" | "listedDate"
type SortOrder = "asc" | "desc"

interface DisplayListing {
  id: string
  make: string
  model: string
  year: number | null
  price: number
  mileage: number | null
  fairPrice: number | null
  listedDate: string
  condition: string
  transmission: string
  fuelType: string
  bodyType: string
  dealRating: string | null
  url?: string
}

function getDealBadgeClasses(rating: string) {
  switch (rating) {
    case "Great Deal":
      return "border-transparent bg-success/15 text-success"
    case "Good Deal":
      return "border-transparent bg-info/15 text-info"
    case "Fair Deal":
      return "border-transparent bg-warning/15 text-warning"
    case "Overpriced":
      return "border-transparent bg-destructive/15 text-destructive"
    default:
      return ""
  }
}

function mapApiListing(l: ApiListing): DisplayListing {
  return {
    id: l.id,
    make: l.make,
    model: l.model,
    year: l.year,
    price: l.price,
    mileage: l.mileage,
    fairPrice: null,
    listedDate: l.first_seen ?? l.created_at ?? new Date().toISOString(),
    condition: l.condition,
    transmission: l.transmission,
    fuelType: l.fuel_type,
    bodyType: "",
    dealRating: null,
    url: l.url,
  }
}

function mapMockListings(): DisplayListing[] {
  return carListings.map((c) => ({
    id: c.id,
    make: c.make,
    model: c.model,
    year: c.year,
    price: c.price,
    mileage: c.mileage,
    fairPrice: c.fairPrice,
    listedDate: c.listedDate,
    condition: c.condition,
    transmission: c.transmission,
    fuelType: c.fuelType,
    bodyType: c.bodyType,
    dealRating: c.dealRating,
  }))
}

interface ListingsTableProps {
  initialListings?: ApiListing[]
}

export function ListingsTable({ initialListings }: ListingsTableProps) {
  const [search, setSearch] = useState("")
  const [makeFilter, setMakeFilter] = useState("all")
  const [sortField, setSortField] = useState<SortField>("listedDate")
  const [sortOrder, setSortOrder] = useState<SortOrder>("desc")

  const allListings = useMemo<DisplayListing[]>(() => {
    if (initialListings && initialListings.length > 0) {
      return initialListings.map(mapApiListing)
    }
    return mapMockListings()
  }, [initialListings])

  const makes = useMemo(
    () => Array.from(new Set(allListings.map((c) => c.make).filter(Boolean))).sort(),
    [allListings]
  )

  const filtered = useMemo(() => {
    let items = [...allListings]

    if (search) {
      const q = search.toLowerCase()
      items = items.filter(
        (c) =>
          c.make.toLowerCase().includes(q) ||
          c.model.toLowerCase().includes(q) ||
          (c.year?.toString() ?? "").includes(q)
      )
    }

    if (makeFilter !== "all") {
      items = items.filter((c) => c.make === makeFilter)
    }

    items.sort((a, b) => {
      let comparison = 0
      switch (sortField) {
        case "price":
          comparison = a.price - b.price
          break
        case "year":
          comparison = (a.year ?? 0) - (b.year ?? 0)
          break
        case "mileage":
          comparison = (a.mileage ?? 0) - (b.mileage ?? 0)
          break
        case "listedDate":
          comparison =
            new Date(a.listedDate).getTime() - new Date(b.listedDate).getTime()
          break
      }
      return sortOrder === "asc" ? comparison : -comparison
    })

    return items
  }, [allListings, search, makeFilter, sortField, sortOrder])

  function toggleSort(field: SortField) {
    if (sortField === field) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc")
    } else {
      setSortField(field)
      setSortOrder("desc")
    }
  }

  return (
    <Card className="border-border/50">
      <CardHeader className="pb-2">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle className="text-base font-semibold">Recent Listings</CardTitle>
            <CardDescription>
              {filtered.length} vehicles found on ecaytrade.com
            </CardDescription>
          </div>
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search make, model..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="h-9 w-48 pl-9 bg-secondary border-border/50"
              />
            </div>
            <Select value={makeFilter} onValueChange={setMakeFilter}>
              <SelectTrigger className="h-9 w-36 bg-secondary border-border/50">
                <SelectValue placeholder="All Makes" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Makes</SelectItem>
                {makes.map((make) => (
                  <SelectItem key={make} value={make}>
                    {make}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardHeader>
      <CardContent className="px-0 pb-0">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent border-border/50">
              <TableHead className="pl-6">Vehicle</TableHead>
              <TableHead
                className="cursor-pointer select-none"
                onClick={() => toggleSort("year")}
              >
                <div className="flex items-center gap-1">
                  Year
                  <ArrowUpDown className="size-3 text-muted-foreground" />
                </div>
              </TableHead>
              <TableHead
                className="cursor-pointer select-none"
                onClick={() => toggleSort("price")}
              >
                <div className="flex items-center gap-1">
                  Price
                  <ArrowUpDown className="size-3 text-muted-foreground" />
                </div>
              </TableHead>
              <TableHead>Fair Price</TableHead>
              <TableHead
                className="cursor-pointer select-none"
                onClick={() => toggleSort("mileage")}
              >
                <div className="flex items-center gap-1">
                  Mileage
                  <ArrowUpDown className="size-3 text-muted-foreground" />
                </div>
              </TableHead>
              <TableHead>Condition</TableHead>
              <TableHead>Deal Rating</TableHead>
              <TableHead
                className="cursor-pointer select-none pr-6"
                onClick={() => toggleSort("listedDate")}
              >
                <div className="flex items-center gap-1">
                  Listed
                  <ArrowUpDown className="size-3 text-muted-foreground" />
                </div>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filtered.map((car) => {
              const priceDiff = car.fairPrice !== null ? car.price - car.fairPrice : null
              const priceDiffPct =
                priceDiff !== null && car.fairPrice
                  ? ((priceDiff / car.fairPrice) * 100).toFixed(1)
                  : null
              const subtitle = [car.bodyType, car.transmission, car.fuelType]
                .filter(Boolean)
                .join(" · ")
              return (
                <TableRow key={car.id} className="border-border/50">
                  <TableCell className="pl-6">
                    <div className="flex flex-col gap-0.5">
                      <div className="flex items-center gap-1.5">
                        <span className="font-medium text-foreground">
                          {car.make} {car.model}
                        </span>
                        {car.url && (
                          <a
                            href={car.url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-muted-foreground hover:text-primary"
                          >
                            <ExternalLink className="size-3" />
                          </a>
                        )}
                      </div>
                      {subtitle && (
                        <span className="text-xs text-muted-foreground">{subtitle}</span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className="font-mono text-foreground">
                    {car.year ?? "—"}
                  </TableCell>
                  <TableCell className="font-mono font-semibold text-foreground">
                    ${car.price.toLocaleString()}
                  </TableCell>
                  <TableCell>
                    {priceDiff !== null && priceDiffPct !== null ? (
                      <div className="flex flex-col">
                        <span className="font-mono text-muted-foreground">
                          ${car.fairPrice!.toLocaleString()}
                        </span>
                        <span
                          className={`text-xs font-mono ${
                            priceDiff > 0 ? "text-destructive" : "text-success"
                          }`}
                        >
                          {priceDiff > 0 ? "+" : ""}
                          {priceDiffPct}%
                        </span>
                      </div>
                    ) : (
                      <span className="text-muted-foreground text-sm">—</span>
                    )}
                  </TableCell>
                  <TableCell className="font-mono text-muted-foreground">
                    {car.mileage != null ? `${car.mileage.toLocaleString()} mi` : "—"}
                  </TableCell>
                  <TableCell>
                    {car.condition ? (
                      <Badge variant="outline" className="text-xs text-foreground border-border/50">
                        {car.condition}
                      </Badge>
                    ) : (
                      <span className="text-muted-foreground text-sm">—</span>
                    )}
                  </TableCell>
                  <TableCell>
                    {car.dealRating ? (
                      <Badge className={getDealBadgeClasses(car.dealRating)}>
                        {car.dealRating}
                      </Badge>
                    ) : (
                      <span className="text-muted-foreground text-sm">—</span>
                    )}
                  </TableCell>
                  <TableCell className="pr-6 text-muted-foreground text-xs font-mono">
                    {new Date(car.listedDate).toLocaleDateString("en-US", {
                      month: "short",
                      day: "numeric",
                    })}
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}

