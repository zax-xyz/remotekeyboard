C := go build
FLAGS := -ldflags "-s -w"

TARGET := remotekeyboard
INSTALL_PATH := /usr/local/bin/

.PHONY: all
.PHONY: install
.PHONY: clean

all: ${TARGET}

${TARGET}: ${TARGET}.go
	${C} ${FLAGS} $^
	upx $@

install:
	install ${TARGET} ${INSTALL_PATH}

clean:
	rm ${TARGET}
