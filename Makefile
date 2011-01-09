include $(GOROOT)/src/Make.inc

TARG=ncurses

CGOFILES=ncurses.go

CGO_LDFLAGS=-lncurses 

CLEANFILES+=sample

include $(GOROOT)/src/Make.pkg

# Simple test programs

%: install %.go
	$(GC) $*.go
	$(LD) -o $@ $*.$O
