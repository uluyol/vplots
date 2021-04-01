package builtin

// Specify OS-specific C++ compiler, linker flags

// #cgo CXXFLAGS: -I${SRCDIR}/xpdf_darwin_amd64/include -std=c++11
// #cgo LDFLAGS: -L${SRCDIR}/xpdf_darwin_amd64/lib -lxpdf
import "C"
