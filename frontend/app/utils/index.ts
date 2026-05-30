type FilterOp = "eq" | "neq" | "gt" | "lt" | "gte" | "lte" | "like";

type Filter<T extends object> = {
  [K in keyof T]?: { op: FilterOp; value: T[K] };
};

export function buildFilters<T extends object>(filters: Filter<T>): string {
  return Object.entries(filters)
    .filter(([, v]) => v != null)
    .map(([col, filter]) => {
      const { op, value } = filter as { op: FilterOp; value: unknown };
      return `${col}:${op}:${value}`;
    })
    .join(";");
}
