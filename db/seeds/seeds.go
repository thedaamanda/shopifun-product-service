package seeds

import (
	"codebase-app/internal/adapter"
	"context"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var userIDs = []string{
	"c396f23e-a097-476d-aae5-cfc9973634f3",
	"a4b7a3f1-751a-4a10-b506-99202581427b",
	"c2776740-3885-444d-a73b-d6426a407792",
	"3b07156c-b997-4803-8428-a4917a12792d",
	"9d523336-6452-4212-89b4-09564f0d5842",
}

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// Run seeds.
func (s *Seed) run(table string, total int) {
	switch table {
	case "categories":
		s.categoriesSeed()
	case "brands":
		s.brandsSeed()
	case "shops":
		s.shopsSeed()
	case "products":
		s.productsSeed()
	case "reviews":
		s.reviewsSeed(total)
	case "all":
		s.categoriesSeed()
		s.brandsSeed()
		s.shopsSeed()
		s.productsSeed()
		s.reviewsSeed(total)
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}

	if table != "" {
		log.Info().Msg("Seed ran successfully")
		log.Info().Msg("Exiting ...")
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
		os.Exit(0)
	}
}

func (s *Seed) deleteAll() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Msg("Error committing transaction")
			}
		}
	}()

	_, err = tx.Exec(`DELETE FROM users`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting users")
		return
	}
	log.Info().Msg("users table deleted successfully")

	_, err = tx.Exec(`DELETE FROM roles`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting roles")
		return
	}
	log.Info().Msg("roles table deleted successfully")

	log.Info().Msg("=== All tables deleted successfully ===")
}

func (s *Seed) categoriesSeed() {
	categoriesMaps := []map[string]any{
		{"name": "Electronics"},
		{"name": "Clothing"},
		{"name": "Books"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO categories (name)
		VALUES (:name)
	`, categoriesMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating categories")
		return
	}

	log.Info().Msg("categories table seeded successfully")
}

func (s *Seed) brandsSeed() {
	categoriesMaps := []map[string]any{
		{"name": "Nike"},
		{"name": "Samsung"},
		{"name": "Lenovo"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO brands (name)
		VALUES (:name)
	`, categoriesMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating brands")
		return
	}

	log.Info().Msg("brands table seeded successfully")
}

func (s *Seed) shopsSeed() {
	shopMaps := make([]map[string]interface{}, 5)

	for i := 0; i < 5; i++ {
		shopMaps[i] = map[string]interface{}{
			"user_id":     userIDs[i],
			"name":        gofakeit.Company(),
			"description": gofakeit.Sentence(10),
			"terms":       gofakeit.Paragraph(2, 2, 10, "\n"),
		}
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO shops (user_id, name, description, terms)
		VALUES (:user_id, :name, :description, :terms)
	`, shopMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating shops")
		return
	}

	log.Info().Msg("shops table seeded successfully")
}

func (s *Seed) productsSeed() {
	var (
		categories []struct{ ID string }
		brands     []struct{ ID string }
		shops      []struct {
			ID     string
			UserId string
		}
	)

	if err := s.db.Select(&categories, "SELECT id FROM categories"); err != nil {
		log.Error().Err(err).Msg("Error fetching categories")
		return
	}
	if err := s.db.Select(&brands, "SELECT id FROM brands"); err != nil {
		log.Error().Err(err).Msg("Error fetching brands")
		return
	}
	if err := s.db.Select(&shops, "SELECT id, user_id as UserId FROM shops"); err != nil {
		log.Error().Err(err).Msg("Error fetching shops")
		return
	}

	productMaps := make([]map[string]interface{}, 0, len(shops)*5)

	for _, shop := range shops {
		for i := 0; i < 5; i++ {
			product := map[string]interface{}{
				"shop_id":     shop.ID,
				"category_id": categories[gofakeit.Number(0, len(categories)-1)].ID,
				"brand_id":    brands[gofakeit.Number(0, len(brands)-1)].ID,
				"name":        gofakeit.ProductName(),
				"description": gofakeit.ProductDescription(),
				"price":       gofakeit.Price(10, 1000),
				"stock":       gofakeit.Number(0, 100),
				"user_id":     shop.UserId,
			}
			productMaps = append(productMaps, product)
		}
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO products (shop_id, category_id, brand_id, name, description, price, stock, user_id)
		VALUES (:shop_id, :category_id, :brand_id, :name, :description, :price, :stock, :user_id)
	`, productMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating products")
		return
	}

	log.Info().Msg("products table seeded successfully")
}

func (s *Seed) reviewsSeed(total int) {
	var (
		products []struct{ ID string }
	)

	if err := s.db.Select(&products, "SELECT id FROM products"); err != nil {
		log.Error().Err(err).Msg("Error fetching products")
		return
	}

	reviewMaps := make([]map[string]interface{}, 0, len(products)*20)

	for _, product := range products {
		for i := 0; i < total; i++ {
			review := map[string]interface{}{
				"product_id": product.ID,
				"user_id":    gofakeit.UUID(),
				"rating":     gofakeit.Number(1, 5),
				"review":     gofakeit.Sentence(10),
			}
			reviewMaps = append(reviewMaps, review)
		}
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO reviews (product_id, user_id, rating, review)
		VALUES (:product_id, :user_id, :rating, :review)
	`, reviewMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating reviews")
		return
	}

	log.Info().Msg("reviews table seeded successfully")
}
