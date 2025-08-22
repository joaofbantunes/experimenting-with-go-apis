package shared

const RootPointer = JSONPointer("#")

type JSONPointer string

func (jp JSONPointer) String() string {
	return string(jp)
}

func JsonPointerForSegments(segments []string) JSONPointer {
	if len(segments) == 0 {
		return RootPointer
	}
	pointer := RootPointer.String()
	for _, segment := range segments {
		pointer += "/" + segment
	}
	return JSONPointer(pointer)
}
