import { Currency, DEV, Token } from 'moonbeamswapdada'

export function currencyId(currency: Currency): string {
  if (currency === DEV) return 'EUN'
  if (currency instanceof Token) return currency.address
  throw new Error('invalid currency')
}
