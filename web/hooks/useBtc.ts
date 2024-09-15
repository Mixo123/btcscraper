import axios from 'axios'
import { useQuery } from '@tanstack/react-query'

export const useCryptocurrencyPrices = () => {
  return useQuery({
    queryKey: ['cryptocurrency_prices'],
    queryFn: async () => {
      try {
        const response = await axios.get('http://localhost:8080/coins')
        return response.data
      } catch (error) {
        console.error('Error fetching cryptocurrency prices:', error)
        throw error
      }
    },
    staleTime: 10000,
  })
}