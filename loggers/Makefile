CC=gcc

#Note - I was compiling on the device. If you want to compile on an x86 box then you will need to setup a cross compiler.

#Note - Had stability issues with thumb instructions. Was hoping to save memory with them
#CFLAGS=-g -Wall -lphidget21 -lpthread -march=armv4t -mtune=arm920t -Os -mthumb

CFLAGS=-g -Wall -lphidget21 -lpthread -mthumb
TARGET=logger

all: $(TARGET)

$(TARGET): $(TARGET).c
	$(CC) $(CFLAGS) -o $(TARGET) $(TARGET).c

clean:
	$(RM) $(TARGET)

