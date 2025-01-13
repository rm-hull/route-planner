package cmds

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rm-hull/route-planner/db"
	"github.com/rm-hull/route-planner/models"
	"github.com/rm-hull/route-planner/repository"
	"github.com/schollz/progressbar/v3"
)

const BATCH_SIZE = 1000
const ESTIMATED_TOTAL_RECORDS = 14_472_914

func ImportGmlData(path string) error {
	config := db.ConfigFromEnv()
	ctx := context.Background()

	pool, err := db.NewDBPool(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	repo, err := repository.NewGmlRepository(pool)
	if err != nil {
		return fmt.Errorf("failed to initialize repo: %v", err)
	}

	err = repo.DisableTriggers(ctx)
	if err != nil {
		return fmt.Errorf("failed to disable triggers: %v", err)
	}
	defer repo.EnableTriggers(ctx)

	files, err := walkFiles(path)
	if err != nil {
		return fmt.Errorf("failed to walk path: %v", err)
	}

	roadLinks := make([]models.RoadLink, 0)
	roadNodes := make([]models.RoadNode, 0)
	motorwayJunctions := make([]models.MotorwayJunction, 0)

	bar := progressbar.Default(ESTIMATED_TOTAL_RECORDS)
	updateProgressBarSummary := func(count int, filePath string) {
		bar.Describe(fmt.Sprintf("importing GML: (%02d/%02d) %s", count + 1, len(files), filePath))
	}

	for i, filePath := range files {
		updateProgressBarSummary(i, filePath)

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		decoder := xml.NewDecoder(file)

		for {
			bar.Add(1)
			token, err := decoder.Token()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				return fmt.Errorf("error reading token: %v", err)
			}

			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == "featureMember" {
					var feature models.FeatureMember
					err := decoder.DecodeElement(&feature, &se)
					if err != nil {
						return fmt.Errorf("error decoding element: %v", err)
					}

					if feature.RoadLink != nil {
						roadLinks = append(roadLinks, *feature.RoadLink)
					} else if feature.RoadNode != nil {
						roadNodes = append(roadNodes, *feature.RoadNode)
					} else if feature.MotorwayJunction != nil {
						motorwayJunctions = append(motorwayJunctions, *feature.MotorwayJunction)
					} else {
						panic(fmt.Sprintf("Unhandled: %v\n", se))
					}
				}
			}
			if len(roadLinks) == BATCH_SIZE {
				err := repo.StoreRoadLinks(ctx, roadLinks...)
				if err != nil {
					return fmt.Errorf("failed to save: %v", err)
				}
				roadLinks = make([]models.RoadLink, 0)
			}
			if len(roadNodes) == BATCH_SIZE {
				err := repo.StoreRoadNodes(ctx, roadNodes...)
				if err != nil {
					return fmt.Errorf("failed to save: %v", err)
				}
				roadNodes = make([]models.RoadNode, 0)
			}
			if len(motorwayJunctions) == BATCH_SIZE {
				err := repo.StoreMotorwayJunctions(ctx, motorwayJunctions...)
				if err != nil {
					return fmt.Errorf("failed to save: %v", err)
				}
				motorwayJunctions = make([]models.MotorwayJunction, 0)
			}
		}
	}

	bar.Finish()

	err = repo.StoreRoadLinks(ctx, roadLinks...)
	if err != nil {
		return fmt.Errorf("failed to save: %v", err)
	}
	err = repo.StoreRoadNodes(ctx, roadNodes...)
	if err != nil {
		return fmt.Errorf("failed to save: %v", err)
	}
	err = repo.StoreMotorwayJunctions(ctx, motorwayJunctions...)
	if err != nil {
		return fmt.Errorf("failed to save: %v", err)
	}
	return nil
}

// walkFiles recursively walks through a folder and returns the relative paths for files.
func walkFiles(root string) ([]string, error) {
	var files []string

	// Walk through the root directory and subdirectories.
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only add files, not directories.
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
