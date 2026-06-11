export const useCountriesQuery = () =>
  useQuery({
    key: ["countries"],
    query: async () => (await getApiCommonCountries({ query: { limit: -1 } })).data?.data ?? [],
  });
