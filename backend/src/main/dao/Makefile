EXT=sqlite3_mod_regexp.so
RM=rm
LDFLAG=-fPIC

all : $(EXT)

$(EXT) : sqlite3_mod_distance.c
	gcc $(LDFLAG) -shared -o $@ $< -lsqlite3 -lpcre

clean :
	@-$(RM) $(EXE) $(EXT)
