package hyperlit

import (
	"fmt"

	"github.com/LeonidS635/HyperLit/internal/info"
)

func (h *HyperLit) printSectionsInfoByStatus(status int) error {
	sections := h.sectionsStatuses.Get(status)
	if len(sections) == 0 {
		return nil
	}

	switch status {
	case info.StatusCreated:
		fmt.Println("Created sections:")
	case info.StatusDeleted:
		fmt.Println("Deleted sections:")
	case info.StatusDocsOutdated:
		fmt.Println("Docs outdated sections:")
	case info.StatusCodeOutdated:
		fmt.Println("Code outdated sections:")
	case info.StatusModified:
		fmt.Println("Updated sections:")
	case info.StatusUnmodified, info.StatusProbablyModified:
	default:
		return fmt.Errorf("error printing sections info: unknown status: %d", status)
	}

	for _, section := range sections {
		fmt.Println(section.Path)
	}
	fmt.Println()

	return nil
}

func (h *HyperLit) printSectionsInfo() error {
	if !h.sectionsStatuses.Check(info.StatusCreated) &&
		!h.sectionsStatuses.Check(info.StatusDeleted) &&
		!h.sectionsStatuses.Check(info.StatusDocsOutdated) &&
		!h.sectionsStatuses.Check(info.StatusCodeOutdated) &&
		!h.sectionsStatuses.Check(info.StatusModified) {

		fmt.Println("No changes found")
		return nil
	}

	if err := h.printSectionsInfoByStatus(info.StatusCreated); err != nil {
		return err
	}
	if err := h.printSectionsInfoByStatus(info.StatusDeleted); err != nil {
		return err
	}
	if err := h.printSectionsInfoByStatus(info.StatusDocsOutdated); err != nil {
		return err
	}
	if err := h.printSectionsInfoByStatus(info.StatusCodeOutdated); err != nil {
		return err
	}
	if err := h.printSectionsInfoByStatus(info.StatusModified); err != nil {
		return err
	}
	return nil
}
