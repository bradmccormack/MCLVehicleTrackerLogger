EXT=sqlite3_mod_distance.so
RM=rm
LDFLAG=-fPIC

all : $(EXT)

$(EXT) : sqlite3_mod_distance.c
	gcc $(LDFLAG) -shared -o $@ $< -lm -lsqlite3

clean :
	@-$(RM) $(EXE) $(EXT)
