"use client"

import { useCryptocurrencyPrices } from "../../hooks/useBtc";

export default function Home() {
  const { data, isLoading, isError } = useCryptocurrencyPrices()

  if (isLoading) return <div>Loading...</div>
  if (isError) return <div>Error loading prices</div>

  return (
    <div className="w-full h-screen flex flex-col justify-center items-center gap-3">
      {Object.entries(data.prices).map(([key, value]) => (
        <p key={key} className="text-3xl">{`${key.charAt(0).toUpperCase() + key.slice(1)} ${value}`}</p>
      ))}
    </div>
  )
}
