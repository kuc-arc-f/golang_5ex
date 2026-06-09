CXX = g++
CXXFLAG_1= -I./include -I/prog/vcpkg/installed/x64-windows/include
CXXFLAG_LIB_2= -L/prog/vcpkg/installed/x64-windows/lib

CXXFLAGS = -shared -std=c++17 $(CXXFLAG_1) $(CXXFLAG_LIB_2)
#g++ -shared -o example.dll example.cpp 

TARGET = example.dll
all: $(TARGET)

$(TARGET): example.o 
	$(CXX) $(CXXFLAGS) example.o -o $(TARGET)

example.o: example.cpp
	$(CXX) $(CXXFLAGS) -c example.cpp

clean:
	del *.o *.exe
