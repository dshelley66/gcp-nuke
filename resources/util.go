package resources

import "strings"

func UnPtrBool(ptr *bool, def bool) bool {
	if ptr == nil {
		return def
	}
	return *ptr
}

func UnPtrString(ptr *string, def string) string {
	if ptr == nil {
		return def
	}
	return *ptr
}

func EqualStringPtr(v1, v2 *string) bool {
	if v1 == nil && v2 == nil {
		return true
	}

	if v1 == nil || v2 == nil {
		return false
	}

	return *v1 == *v2
}

func ZoneInRegionList(zone string, regions []string) bool {
	for _, region := range regions {
		if strings.Contains(zone, region) {
			return true
		}
	}
	return false
}
