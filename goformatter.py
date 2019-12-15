import ctypes

class GoSlice(ctypes.Structure):
	_fields_ = [
		('data', ctypes.c_void_p),
		('len', ctypes.c_longlong),
		('cap', ctypes.c_longlong),
	]

def decode(lib, name, blob):
	buf = ctypes.create_string_buffer(blob, len(blob))
	s = GoSlice(ctypes.addressof(buf), len(blob), len(blob))
	lib.CGobFormat.restype = ctypes.c_void_p
	lib.CGobFormat.argtypes = [ctypes.c_char_p, GoSlice]
	c = lib.CGobFormat(name, s)
	json_text = ctypes.cast(c, ctypes.c_char_p).value
	lib.CGobFree.argtypes = [ctypes.c_void_p]
	lib.CGobFree(c)
	return json_text

#lib = ctypes.CDLL('./gob.so')
#print(decode(lib, b'', b''))