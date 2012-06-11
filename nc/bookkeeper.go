package nc

// Bookkeeper allows and boxes and channels to be registered
// so that an informative graphviz graph can be made.

import (
	"fmt"
	//"log"
	"reflect"
	"strings"
	"unicode"
	"unsafe"
)

type ChanPtr unsafe.Pointer

var (
	connections [][2]Chan

	boxes []Box

	boxOfChan = make(map[ChanPtr]Box)

	// TODO: Quant type
	//srcBoxForQuant  = make(map[string]Box)  //IDEA: store vector/scalar kind here too: Quant
	//srcChanForQuant = make(map[string]Chan) //IDEA: store vector/scalar kind here too

	//boxForChanPtr
	//quantForChanPtr
	//->

	//Quant: Box, srcChan

	//QuantForName
	//BoxForChan
)

func WriteGraph() {
	var dot graphvizwriter
	dot.Init("plumbing.dot")
	defer dot.Close()

	for _, box := range boxes {
		dot.AddBox(boxname(box))
	}
}

func RegisterConn(dst, src Chan) {
	connections = append(connections, [2]Chan{dst, src})
}

// Checks if ptr is the address of a supported channel type.
func isChan(ptr interface{}) bool {
	return isInputChan(ptr) || isOutputChan(ptr)
}

func Register(box ...Box){
	for _,b:=range box{RegisterBox(b)}
}

// Registers quantities for all tagged output channels.
func RegisterBox(box Box) {
	boxes = append(boxes, box)

	// find and map all output channels
	val := reflect.ValueOf(box).Elem()
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		tag := typ.Field(i).Tag
		field := val.Field(i)

		// skip non-exported fields
		if !unicode.IsUpper(rune(name[0])) { // who uses unicode identifiers anyway?
			continue
		}

		ptr := field.Addr().Interface()

		if isChan(ptr) {
			boxOfChan[ChanPtr(unsafe.Pointer(field.UnsafeAddr()))] = box // .Addr().Unsa..?
		}

		if isInputChan(ptr) {

		}

		// use the field's struct tag as channel name
		if tag == "" {
			continue
		}

		//	if isOutputChan(chanPtr) {
		//		if prev, ok := srcBoxForQuant[name]; ok {
		//			panic(name + " provided by both " + boxname(prev) + " and " + boxname(box))
		//		}
		//		RegisterTag(box, chanPtr, name)
		//	}
	}
}

// Register a quantity taken form channel, give it a name.
//func RegisterTag(box Box, chanPtr Chan, name string) { // rm box, use boxforchanptr, mv SendQuant
//	Assert(isOutputChan(chanPtr))
//	setBoxFor(name, box)
//	setChanFor(name, chanPtr)
//	log.Println("[plumber]", boxname(box), "provides", name)
//	registerComponentQuants(box, chanPtr, name)
//}

// Automatically register components of vector quantities as "quant.x", ...
//func registerComponentQuants(box Box, chanPtr Chan, name string) { // rm box
//	Assert(isOutputChan(chanPtr))
//	switch c := chanPtr.(type) {
//	default:
//		return // nothing to see here
//	case *[3][]chan<- []float32:
//		RegisterTag(box, &c[X], name+".x")
//		RegisterTag(box, &c[Y], name+".y")
//		RegisterTag(box, &c[Z], name+".z")
//	case *[3][]chan<- float64:
//		RegisterTag(box, &c[X], name+".x")
//		RegisterTag(box, &c[Y], name+".y")
//		RegisterTag(box, &c[Z], name+".z")
//	}
//}

//func srcBoxFor(quant string) Box {
//	if b, ok := srcBoxForQuant[quant]; ok {
//		return b
//	}
//	panic("no such quantity " + quant)
//}
//
//func srcChanFor(quant string) Chan {
//	if b, ok := srcChanForQuant[quant]; ok {
//		return b
//	}
//	panic("no such quantity " + quant)
//}
//
//func setBoxFor(quant string, box Box) {
//	if _, ok := srcBoxForQuant[quant]; ok {
//		panic("already defined " + quant)
//	}
//	srcBoxForQuant[quant] = box
//}
//
//func setChanFor(quant string, c Chan) {
//	if _, ok := srcChanForQuant[quant]; ok {
//		panic("already defined " + quant)
//	}
//	srcChanForQuant[quant] = c
//}

func boxname(value interface{}) string {
	typ := fmt.Sprintf("%T", value)
	clean := typ[strings.Index(typ, ".")+1:] // strip "*mm."
	if strings.HasSuffix(clean, "Box") {
		clean = clean[:len(clean)-len("Box")]
	}
	return clean
}
