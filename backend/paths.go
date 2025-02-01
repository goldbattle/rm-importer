package backend

import (
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

/*
Returns a path for creating a file.
Normalizes folderName and item.TabletPath.
*/
func getFilePath(location string, folderName string, item DocInfo, ext string) (string, error) {
	itemPath := slices.Clone(item.TabletPath)
	if len(itemPath) == 0 {
		return "", fmt.Errorf("Item path is empty! id=%v", item.Id)
	}

	last := itemPath[len(itemPath)-1]
	if path.Ext(last) != "."+ext {
		itemPath[len(itemPath)-1] = last + "." + ext
	}

	for i := range itemPath {
		itemPath[i] = normalize(itemPath[i])
	}
	folderName = normalize(folderName)

	toJoin := []string{filepath.ToSlash(location), folderName}
	toJoin = append(toJoin, itemPath...)
	return path.Join(toJoin...), nil
}

/* Normalizes folder or a file name. */
/* Replaces banned characters (0x00 up to 0x1F, /"*:<>?\|),
   changes the name if reserved on Windows (CON, AUX, NUL, etc.)
*/
func normalize(name string) string {
	banned := []rune{'"', '*', '/', ':', '<', '>', '?', '\\', '|'}
	banned = append(banned, 0x7F)
	for i := range 0x20 {
		banned = append(banned, rune(i))
	}
	slices.Sort(banned)

	sb := strings.Builder{}

	for _, c := range name {
		if _, ok := slices.BinarySearch(banned, c); ok {
			sb.WriteRune('-')
		} else {
			sb.WriteRune(c)
		}
	}

	name = sb.String()
	ext := path.Ext(name)

	reservedWin := []string{"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3",
		"COM4", "COM5", "COM6",
		"COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3",
		"LPT4", "LPT5", "LPT6",
		"LPT7", "LPT8", "LPT9"}

	nameWithoutExt := strings.TrimSuffix(name, ext)
	if slices.Contains(reservedWin, nameWithoutExt) {
		nameWithoutExt += "-1"
	}

	return nameWithoutExt + ext
}
