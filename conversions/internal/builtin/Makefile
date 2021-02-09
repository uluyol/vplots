CXXFLAGS += -Ixpdf_darwin_amd64/include
LDFLAGS += -Lxpdf_darwin_amd64/lib -lxpdf

debug_convert: debug_convert.c convert.cc
	$(CXX) -c convert.cc $(CXXFLAGS)
	$(CC) -c debug_convert.c
	$(CC) -o debug_convert debug_convert.o convert.o $(LDFLAGS) -lc++
