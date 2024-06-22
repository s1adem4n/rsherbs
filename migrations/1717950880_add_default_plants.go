package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	m "github.com/pocketbase/pocketbase/migrations"
)

type Plants [][2]string

var data = `
[
  [
    "Achillea millefolium",
    "Schafgarbe"
  ],
  [
    "Alchemilla",
    "Frauenmantel"
  ],
  [
    "Allium ursinum",
    "Bärlauch"
  ],
  [
    "Artemisia Absinthium",
    "Wermut"
  ],
  [
    "Bellis perennis",
    "Gänseblümchen"
  ],
  [
    "Cynara scolymus",
    "Artischocke"
  ],
  [
    "Echinacea purpurea",
    "Sonnenhut "
  ],
  [
    "Gingko",
    "Gingko"
  ],
  [
    "Lavendula",
    "Lavendel"
  ],
  [
    "Melissa officinalis",
    "Zitronenmelisse"
  ],
  [
    "Mentha piperita",
    "Pfefferminze"
  ],
  [
    "Salvia",
    "Salbei"
  ],
  [
    "Salvia rosmarinus",
    "Rosmarin"
  ],
  [
    "Thymus vulgaris",
    "Thymian"
  ],
  [
    "Valeriana",
    "Baldriane"
  ],
  [
    "Viola tricolor",
    "Stiefmütterchen"
  ]
]
`

func init() {
	var plants Plants
	err := json.Unmarshal([]byte(data), &plants)
	if err != nil {
		panic(err)
	}

	m.Register(func(db dbx.Builder) error {
		for _, plant := range plants {
			var count int
			err := db.
				Select("count(*)").
				From("plants").
				Where(dbx.NewExp("latin = {:latin}", dbx.Params{
					"latin": plant[0],
				})).
				Row(&count)
			if err != nil {
				return err
			}
			if count == 0 {
				_, err := db.Insert("plants", dbx.Params{
					"latin": plant[0],
					"name":  plant[1],
				}).Execute()
				if err != nil {
					return err
				}
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
