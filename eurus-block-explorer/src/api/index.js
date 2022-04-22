const BASE_URL = process.env.VUE_APP_BASE_URL
export function searchQueryAPI(id='') {
  return `${process.env.VUE_APP_SEARCH_URL}/api/search?query=${id}`
}

export function transactionAPI(href) {
  return `${BASE_URL}${href}`
}
export function transactionEventAPI(href) {
  return `${BASE_URL}${href}/events`
}

export function blockAPI(id) {
  return `${BASE_URL}/blocks/${id}`
}
export function blockTransAPI(id, query='') {
  return `${BASE_URL}/blocks/${id}/transactions${query}`
}