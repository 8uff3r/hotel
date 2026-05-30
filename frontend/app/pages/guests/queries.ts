export const useCountriesQuery = () =>
  useQuery({
    key: ["countries"],
    query: async () => (await getApiCommonCountries()).data,
  });
