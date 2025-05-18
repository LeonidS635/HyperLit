package hyperlit

import (
	"fmt"

	"github.com/LeonidS635/HyperLit/internal/info"
)

func (h *HyperLit) printSectionsInfoByStatus(status int) error {
	sections := h.sectionsStates.Get(status)
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

	for _, sectionState := range sections {
		fmt.Println(sectionState.Path)
	}
	fmt.Println()

	return nil
}

func (h *HyperLit) printSectionsInfo() (bool, error) {
	if !h.sectionsStates.Check(info.StatusCreated) &&
		!h.sectionsStates.Check(info.StatusDeleted) &&
		!h.sectionsStates.Check(info.StatusDocsOutdated) &&
		!h.sectionsStates.Check(info.StatusCodeOutdated) &&
		!h.sectionsStates.Check(info.StatusModified) {

		fmt.Println("No changes found")
		return false, nil
	}

	if err := h.printSectionsInfoByStatus(info.StatusCreated); err != nil {
		return false, err
	}
	if err := h.printSectionsInfoByStatus(info.StatusDeleted); err != nil {
		return false, err
	}
	if err := h.printSectionsInfoByStatus(info.StatusDocsOutdated); err != nil {
		return false, err
	}
	if err := h.printSectionsInfoByStatus(info.StatusCodeOutdated); err != nil {
		return false, err
	}
	if err := h.printSectionsInfoByStatus(info.StatusModified); err != nil {
		return false, err
	}
	return true, nil
}
