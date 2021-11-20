package nbt

type Compound struct {
	name        string
	value       []interface{}
	previousTag *Compound //To go back
}

// CreateCompoundTag - Create compound tag
func CreateCompoundTag(name string) Compound {
	C := *new(Compound)
	C.name = name
	C.value = make([]interface{}, 0, 16)
	return C
}

// SetPreviousTag - This is an unsafe function and I would reccomend avoiding it unless you have to set it because of multiple compounds in a list
func SetPreviousTag(C *Compound, previousTag *Compound) {
	C.previousTag = previousTag
}

// AddCompoundTag - Add compound tag and sets it in NBTWriter so objects are written to it
func (NBTE *NBTEncoder) AddCompoundTag(name string) {
	C := new(Compound)
	C.name = name
	C.value = make([]interface{}, 0, 16)
	C.previousTag = NBTE.currentCompound
	NBTE.currentCompound = C
}

func (NBTE *NBTEncoder) EndCompoundTag() {
	if NBTE.currentCompound.previousTag != nil { // if nil then it's root tag
		NBTE.currentCompound.EndTag()
		NBTE.currentCompound.previousTag.value = append(NBTE.currentCompound.previousTag.value, *NBTE.currentCompound)
		NBTE.currentCompound = NBTE.currentCompound.previousTag // go back
	} else { // if in root tag
		NBTE.currentCompound.EndTag()
	}
}

func (NBTE *NBTEncoder) EncodeCompound(C *Compound) {
	NBTE.EncodeTag(TagCompound, C.name)
	for _, v := range C.value {
		switch val := v.(type) {
		case End:
			NBTE.data = append(NBTE.data, TagEnd)
			return
		case Byte:
			NBTE.EncodeByte(val.Name, val.Value)
		case Short:
			NBTE.EncodeShort(val.Name, val.Value)
		case Int:
			NBTE.EncodeInt(val.Name, val.Value)
		case Long:
			NBTE.EncodeLong(val.Name, val.Value)
		case Float:
			NBTE.EncodeFloat(val.Name, val.Value)
		case Double:
			NBTE.EncodeDouble(val.Name, val.Value)
		case ByteArray:
			NBTE.EncodeByteArray(val.Name, val.Value)
		case String:
			NBTE.EncodeString(val.Name, val.Value)
		case List:
			NBTE.EncodeList(val)
		case Compound:
			NBTE.EncodeCompound(&val)
		case IntArray:
			NBTE.EncodeIntArray(val.Name, val.Value)
		case LongArray:
			NBTE.EncodeLongArray(val.Name, val.Value)
		}
	}
}

func (TC *Compound) AddTag(val interface{}) {
	TC.value = append(TC.value, val)
}

func (TC *Compound) EndTag() {
	TC.value = append(TC.value, End{})
}

func (TC *Compound) AddMultipleTags(val []interface{}) {
	TC.value = append(TC.value, val...)
}
