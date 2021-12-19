# sliceconverter
Fast conversion from a slice of any type to a slice of any other type.

Conversions can be made as long as the src slice len and cap are multiples of
the dst slice item size. If those requirements are not met, code will panic.

## USage:

src := []MyFancyStruct{ ... }

var dst []byte  // Just so type inference can quick in.
dst = sliceconverter.Convert(src)
