export const Permissions = {
  dashboard: {
    index: {
      read: "index:read",
    },
  },

  guests: {
    guests: {
      read: "guests:read",
      create: "guests:create",
      update: "guests:update",
      delete: "guests:delete",
      export: "guests:export",
    },
    guestsSettle: {
      read: "guests-settle:read",
    },
  },

  rooms: {
    rooms: {
      read: "rooms:read",
      create: "rooms:create",
      update: "rooms:update",
      delete: "rooms:delete",
    },
    roomsRack: {
      read: "rooms-rack:read",
    },
  },

  reservations: {
    reservations: {
      read: "reservations:read",
      create: "reservations:create",
      update: "reservations:update",
      delete: "reservations:delete",
    },
  },

  users: {
    users: {
      read: "users:read",
      create: "users:create",
      update: "users:update",
      delete: "users:delete",
    },
  },

  accounting: {
    accounting: {
      read: "accounting:read",
    },
    accountingAccounts: {
      read: "accounting-accounts:read",
      create: "accounting-accounts:create",
      update: "accounting-accounts:update",
      delete: "accounting-accounts:delete",
    },
    accountingExpenses: {
      read: "accounting-expenses:read",
      create: "accounting-expenses:create",
      update: "accounting-expenses:update",
      delete: "accounting-expenses:delete",
      export: "accounting-expenses:export",
    },
    accountingIncome: {
      read: "accounting-income:read",
      create: "accounting-income:create",
      update: "accounting-income:update",
      delete: "accounting-income:delete",
      export: "accounting-income:export",
    },
  },

  parking: {
    parking: {
      read: "parking:read",
    },
    parkingLots: {
      read: "parking-lots:read",
      create: "parking-lots:create",
      update: "parking-lots:update",
      delete: "parking-lots:delete",
    },
    parkingSpots: {
      read: "parking-spots:read",
      create: "parking-spots:create",
      update: "parking-spots:update",
      delete: "parking-spots:delete",
    },
    parkingVehicles: {
      read: "parking-vehicles:read",
      create: "parking-vehicles:create",
      update: "parking-vehicles:update",
      delete: "parking-vehicles:delete",
    },
    parkingTransactions: {
      read: "parking-transactions:read",
      create: "parking-transactions:create",
      update: "parking-transactions:update",
      export: "parking-transactions:export",
    },
    parkingTransactionsCheckIn: {
      read: "parking-transactions-check-in:read",
    },
  },
};
