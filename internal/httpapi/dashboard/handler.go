package dashboard

import (
	"fmt"
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type DashboardModule struct{}

func (m DashboardModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/stats", m.statsHandler(api))
	fuego.Get(s, "/recent-reservations", m.recentReservationsHandler(api))
}

type DashboardStats struct {
	TotalRooms     int     `json:"totalRooms"`
	AvailableRooms int     `json:"availableRooms"`
	OccupiedRooms  int     `json:"occupiedRooms"`
	OccupancyRate  float64 `json:"occupancyRate"`
	TodaysRevenue  float64 `json:"todaysRevenue"`
	CheckInsToday  int     `json:"checkInsToday"`
	CheckOutsToday int     `json:"checkOutsToday"`
}

type RecentReservation struct {
	ID              uint      `json:"id"`
	ReservationCode string    `json:"reservationCode"`
	GuestName       string    `json:"guestName"`
	RoomNumber      string    `json:"roomNumber"`
	Status          string    `json:"status"`
	EntryDate       time.Time `json:"entryDate"`
}

type RecentReservationsResponse struct {
	Data []RecentReservation `json:"data"`
}

func (m DashboardModule) statsHandler(api *h.API) func(c fuego.ContextNoBody) (DashboardStats, error) {
	return func(c fuego.ContextNoBody) (DashboardStats, error) {
		var stats DashboardStats
		ctx := c.Context()
		db := api.Db.WithContext(ctx)
		hotelID := h.GetHotelIDFromContext(ctx)

		// Total rooms
		var totalRooms int64
		if err := db.Model(&models.Room{}).Where("hotel_id = ?", hotelID).Count(&totalRooms).Error; err != nil {
			return stats, fmt.Errorf("count rooms: %w", err)
		}
		stats.TotalRooms = int(totalRooms)

		// Occupied rooms: rooms with status = occupied
		var occupiedRooms int64
		if err := db.Model(&models.Room{}).
			Joins("JOIN room_statuses ON rooms.status_id = room_statuses.id").
			Where("rooms.hotel_id = ? AND room_statuses.slug = ?", hotelID, string(models.RoomStatusOccupied)).
			Count(&occupiedRooms).Error; err != nil {
			return stats, fmt.Errorf("count occupied: %w", err)
		}
		stats.OccupiedRooms = int(occupiedRooms)

		// Available rooms: total minus occupied (excluding cleaning/under_repair)
		var unavailableRooms int64
		if err := db.Model(&models.Room{}).
			Joins("JOIN room_statuses ON rooms.status_id = room_statuses.id").
			Where("rooms.hotel_id = ? AND room_statuses.slug IN (?, ?)", hotelID, string(models.RoomStatusCleaning), string(models.RoomStatusUnderRepair)).
			Count(&unavailableRooms).Error; err != nil {
			return stats, fmt.Errorf("count unavailable: %w", err)
		}
		stats.AvailableRooms = int(totalRooms - occupiedRooms - unavailableRooms)
		if stats.AvailableRooms < 0 {
			stats.AvailableRooms = 0
		}

		// Occupancy rate
		if stats.TotalRooms > 0 {
			stats.OccupancyRate = float64(stats.OccupiedRooms) / float64(stats.TotalRooms) * 100
		}

		now := time.Now()
		startOfDay := now.Truncate(24 * time.Hour)
		endOfDay := startOfDay.Add(24 * time.Hour)

		// Today's revenue: sum of income amounts for today
		var todaysRevenue float64
		if err := db.Model(&models.Income{}).
			Where("hotel_id = ? AND income_date >= ? AND income_date < ?", hotelID, startOfDay, endOfDay).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&todaysRevenue).Error; err != nil {
			return stats, fmt.Errorf("sum revenue: %w", err)
		}
		stats.TodaysRevenue = todaysRevenue

		// Check-ins today: stays with entry_date = today
		var checkInsToday int64
		if err := db.Model(&models.Stay{}).
			Where("hotel_id = ? AND entry_date >= ? AND entry_date < ?", hotelID, startOfDay, endOfDay).
			Count(&checkInsToday).Error; err != nil {
			return stats, fmt.Errorf("count check-ins: %w", err)
		}
		stats.CheckInsToday = int(checkInsToday)

		// Check-outs today: stays with departure_date = today
		var checkOutsToday int64
		if err := db.Model(&models.Stay{}).
			Where("hotel_id = ? AND departure_date >= ? AND departure_date < ?", hotelID, startOfDay, endOfDay).
			Count(&checkOutsToday).Error; err != nil {
			return stats, fmt.Errorf("count check-outs: %w", err)
		}
		stats.CheckOutsToday = int(checkOutsToday)

		return stats, nil
	}
}

func (m DashboardModule) recentReservationsHandler(api *h.API) func(c fuego.ContextNoBody) (RecentReservationsResponse, error) {
	return func(c fuego.ContextNoBody) (RecentReservationsResponse, error) {
		ctx := c.Context()
		db := api.Db.WithContext(ctx)
		hotelID := h.GetHotelIDFromContext(ctx)

		var reservations []models.Reservation
		if err := db.Where("hotel_id = ?", hotelID).
			Order("id DESC").
			Limit(5).
			Preload("Guest").
			Preload("Status").
			Preload("Rooms").
			Find(&reservations).Error; err != nil {
			return RecentReservationsResponse{}, fmt.Errorf("fetch reservations: %w", err)
		}

		var result []RecentReservation
		for _, r := range reservations {
			roomNumber := ""
			if len(r.Rooms) > 0 {
				roomNumber = r.Rooms[0].RoomNumber
			}
			result = append(result, RecentReservation{
				ID:              r.ID,
				ReservationCode: r.ReservationCode,
				GuestName:       r.Guest.FirstName + " " + r.Guest.LastName,
				RoomNumber:      roomNumber,
				Status:          r.Status.Label,
				EntryDate:       r.EntryDate,
			})
		}

		return RecentReservationsResponse{Data: result}, nil
	}
}
