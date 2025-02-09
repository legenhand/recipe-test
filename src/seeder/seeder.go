package seeder

import (
	"github.com/google/uuid"
	"github.com/legenhand/recipe-test/src/model"
	"gorm.io/gorm"
	"log"
	"time"
)

func SeedUnits(db *gorm.DB) {
	units := []model.Unit{
		{ID: uuid.New(), Code: "kg", Name: "Kilogram", IsBase: true, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Code: "gr", Name: "Gram", IsBase: false, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Code: "l", Name: "Liter", IsBase: true, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Code: "ml", Name: "Milliliter", IsBase: false, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Code: "pcs", Name: "Piece", IsBase: true, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, unit := range units {
		err := db.Where(model.Unit{Code: unit.Code}).FirstOrCreate(&unit).Error
		if err != nil {
			log.Printf("Failed to seed unit %s: %v", unit.Code, err)
		}
	}
	log.Println("✅ Units seeded successfully!")
}

// SeedUnitConversions inserts unit conversion data
func SeedUnitConversions(db *gorm.DB) {
	var kg, g, l, ml model.Unit

	if err := db.Where("code = ?", "kg").First(&kg).Error; err != nil {
		log.Fatalf("❌ Failed to find unit kg: %v", err)
	}
	if err := db.Where("code = ?", "gr").First(&g).Error; err != nil {
		log.Fatalf("❌ Failed to find unit g: %v", err)
	}
	if err := db.Where("code = ?", "l").First(&l).Error; err != nil {
		log.Fatalf("❌ Failed to find unit l: %v", err)
	}
	if err := db.Where("code = ?", "ml").First(&ml).Error; err != nil {
		log.Fatalf("❌ Failed to find unit ml: %v", err)
	}

	conversions := []model.UnitConversion{
		{ID: uuid.New(), BaseID: kg.ID, TargetID: g.ID, Value: 1000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), BaseID: l.ID, TargetID: ml.ID, Value: 1000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, conv := range conversions {
		err := db.Where("base_id = ? AND target_id = ?", conv.BaseID, conv.TargetID).FirstOrCreate(&conv).Error
		if err != nil {
			log.Printf("Failed to seed conversion from %v to %v: %v", conv.BaseID, conv.TargetID, err)
		}
	}
	log.Println("✅ Unit conversions seeded successfully!")
}

func SeedInventory(db *gorm.DB) {
	// Inventory items with UnitID as string for easier mapping
	inventoryItems := []struct {
		Name     string
		Qty      float64
		UnitName string // Use UnitName to find UnitID
		Price    float64
	}{
		{Name: "Aren Sugar", Qty: 1, UnitName: "kg", Price: 60000},
		{Name: "Milk", Qty: 1, UnitName: "l", Price: 30000},
		{Name: "Ice Cube", Qty: 1, UnitName: "kg", Price: 15000},
		{Name: "Plastic Cup", Qty: 10, UnitName: "pcs", Price: 5000},
		{Name: "Coffee Bean", Qty: 1, UnitName: "kg", Price: 100000},
		{Name: "Mineral Water", Qty: 1, UnitName: "l", Price: 5000},
	}

	for _, item := range inventoryItems {
		var unit model.Unit
		// Find the Unit ID based on the Unit name
		if err := db.Where("code = ?", item.UnitName).First(&unit).Error; err != nil {
			log.Printf("Error finding unit '%s': %v", item.UnitName, err)
			continue
		}

		// Map found UnitID to the inventory item and create or update
		inventory := model.Inventory{
			Name:   item.Name,
			Qty:    item.Qty,
			UnitID: unit.ID, // Set the correct UnitID
			Price:  item.Price,
		}

		// Create the inventory item or update if it already exists
		db.Where(model.Inventory{Name: item.Name}).FirstOrCreate(&inventory)
	}

	log.Println("✅ Unit Inventory seeded successfully!")
}

func RunSeeds(db *gorm.DB) {
	SeedUnits(db)
	SeedUnitConversions(db)
	SeedInventory(db)
	log.Println("✅ Database seeding completed!")
}
