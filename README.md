# mapstruct

mapstruct is a Go library for decoding generic map values to structures. Unlike similar open source projects, such as [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) or [tiaotiao/mapstruct](https://github.com/tiaotiao/mapstruct), the purpose of this library is try to be **permissive** and **robust**.  

## features

* does **not** define proprietary struct tags, use `json:` and `yaml:` tag instead.
* support the majority of Go data types, from simple types like `int`, `string` to compound types like `map`, `slice` or `struct`.  Following types are **not** supported:
	* Complex64
	* Complex128
	* Chan
	* Func
	* Interface
	* UnsafePointer
* pointer to `array`, `slice` or `map` are **not** supported.
* support recursive data structures.
* tollerate data mismatch whenever possible.  For example, if the input data is `[3]int`, and struct field is `[2]int`, only the first two elements are used. On the contrary, if input data is `[2]int` and field type is `[3]int`, the last element of the field has its "zero value" automatically.
* input data can be either `map[string]interface{}` or `map[interface{}]interface{}` (however, the real data type of keys must be string).

See the test case for illustration of usage and supported data types.
